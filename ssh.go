package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type (
	grapeSSH struct {
		keySigner ssh.Signer
	}
	grapeSSHClient struct {
		*ssh.Client
	}
	grapeSSHSession struct {
		*ssh.Session
	}
	std struct {
		Out string
		Err string
	}
	grapeCommandStd struct {
		Command command
		Std     std
	}
	sshError error
)

func (gSSH *grapeSSH) newError(errMsg string) sshError {
	return errors.New(errMsg)
}

func (client *grapeSSHClient) newError(errMsg string) sshError {
	return errors.New(errMsg)
}

func (gSSH *grapeSSH) setKey(keyPath keyPath) sshError {
	privateBytes, err := ioutil.ReadFile(string(keyPath))
	if err != nil {
		return gSSH.newError("Could not open idendity file.")
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return gSSH.newError(fmt.Sprintf("Could not parse idendity file."))
	}
	gSSH.keySigner = privateKey
	return nil
}

func (gSSH *grapeSSH) newClient(server server) (*grapeSSHClient, sshError) {
	client, err := ssh.Dial("tcp", server.Host, &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(gSSH.keySigner),
		},
	})
	if err != nil {
		return nil, gSSH.newError(fmt.Sprintf("Could not established ssh connection to server [%s].", server.Host))
	}
	return &grapeSSHClient{client}, nil
}

func (client *grapeSSHClient) newSession() (*grapeSSHSession, sshError) {
	session, err := client.NewSession()
	if err != nil {
		return nil, client.newError(fmt.Sprintf("Could not establish session [%s].", client.Client.RemoteAddr()))
	}
	return &grapeSSHSession{session}, nil
}

func (client *grapeSSHClient) exec(command command) (*grapeCommandStd, sshError) {

	commandStd := &grapeCommandStd{
		Command: command,
		Std: std{
			Err: "",
			Out: "",
		},
	}

	session, err := client.newSession()

	if err != nil {
		return nil, err
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	session.Run(string(command))
	session.Close()

	commandStd.Std = std{
		Err: stderr.String(),
		Out: stdout.String(),
	}

	return commandStd, nil

}
