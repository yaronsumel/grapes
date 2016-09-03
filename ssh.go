package main

import (
	"bytes"
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
)

func (gSSH *grapeSSH) setKey(keyPath keyPath) {
	privateBytes, err := ioutil.ReadFile(string(keyPath))
	if err != nil {
		fatal(fmt.Sprintf("Could not open idendity file."))
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		fatal(fmt.Sprintf("Could not parse idendity file."))
	}
	gSSH.keySigner = privateKey
}

func (gSSH *grapeSSH) newClient(server server) grapeSSHClient {
	client, err := ssh.Dial("tcp", server.Host, &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(gSSH.keySigner),
		},
	})
	if err != nil {
		fatal(fmt.Sprintf("Could not establish ssh connection to server [%s].", server.Host))
	}
	return grapeSSHClient{client}
}

func (client *grapeSSHClient) newSession() *grapeSSHSession {
	session, err := client.NewSession()
	if err != nil {
		fatal(fmt.Sprintf("Could not establish session [%s].", client.Client.RemoteAddr()))
	}
	return &grapeSSHSession{session}
}

func (client *grapeSSHClient) exec(command command) *grapeCommandStd {

	session := client.newSession()

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	session.Run(string(command))
	session.Close()

	return &grapeCommandStd{
		Command: command,
		Std: std{
			Err: stderr.String(),
			Out: stdout.String(),
		},
	}

}
