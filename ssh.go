package main

import (
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"bytes"
	"fmt"
)

type (
	grapeSsh struct {
		keySigner ssh.Signer
	}
	grapeSshClient struct {
		*ssh.Client
	}
	grapeSshSession struct {
		*ssh.Session
	}
	Std struct {
		Out string
		Err string
	}
	grapeCommandStd struct {
		Command command
		Std     Std
	}
)

func (this *grapeSsh)setKey(keyPath keyPath) {
	privateBytes, err := ioutil.ReadFile(string(keyPath))
	if err != nil {
		fatal(fmt.Sprintf("Could not open idendity file."))
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		fatal(fmt.Sprintf("Could not parse idendity file."))
	}
	this.keySigner = privateKey
}

func (this *grapeSsh)newClient(server server) grapeSshClient {
	client, err := ssh.Dial("tcp", server.Host, &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(this.keySigner),
		},
	})
	if err != nil {
		fatal(fmt.Sprintf("Could not establish ssh connection to server [%s].",server.Host))
	}
	return grapeSshClient{client}
}

func (client *grapeSshClient)newSession() *grapeSshSession {
	session, err := client.NewSession()
	if err != nil {
		fatal(fmt.Sprintf("Could not establish session [%s].",client.Client.RemoteAddr()))
	}
	return &grapeSshSession{session}
}

func (session *grapeSshSession)exec(command command) *grapeCommandStd {

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	session.Run(string(command))
	session.Close()

	return &grapeCommandStd{
		Command:command,
		Std:Std{
			Err:stderr.String(),
			Out:stdout.String(),
		},
	}

}