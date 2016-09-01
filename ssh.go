package main

import (
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"log"
	"bytes"
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
	privateBytes, err := ioutil.ReadFile(string(*keyPath))
	if err != nil {
		log.Fatal("setKey", err)
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("setKey", err)
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
		panic(err)
	}
	return grapeSshClient{client}
}

func (client *grapeSshClient)newSession() *grapeSshSession {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &grapeSshSession{session}
}

func (session *grapeSshSession)exec(command command) *grapeCommandStd {

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	session.Run(string(*command))
	session.Close()

	return &grapeCommandStd{
		Command:command,
		Std:Std{
			Err:stderr.String(),
			Out:stdout.String(),
		},
	}

}