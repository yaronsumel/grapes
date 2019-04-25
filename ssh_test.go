package main

import (
	"testing"

	"golang.org/x/crypto/ssh/knownhosts"
)

func TestSetKey(t *testing.T) {
	gSSH := grapeSSH{}
	if err := gSSH.setKey(keyPath("testFiles/id_rsa")); err != nil {
		t.FailNow()
	}
	if err := gSSH.setKey(keyPath("testFiles/id_rsa_123")); err == nil {
		t.FailNow()
	}
	if err := gSSH.setKey(keyPath("testFiles/id_rsa.pub")); err == nil {
		t.FailNow()
	}
}

func TestNewClient(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")
	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	hostKeyCallback, err := knownhosts.New("testFiles/known_hosts")
	if err != nil {
		t.FailNow()
	}

	if _, err := gSSH.newClient(s, hostKeyCallback); err != nil {
		t.FailNow()
	}
	s.Host = "localhost"
	if _, err := gSSH.newClient(s, hostKeyCallback); err == nil {
		t.FailNow()
	}
}

func TestExecCommand(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")
	var client *grapeSSHClient
	var err error
	var output *sshOutput
	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}
	hostKeyCallback, err := knownhosts.New("testFiles/known_hosts")
	if err != nil {
		t.Fatal(err)
	}

	client, err = gSSH.newClient(s, hostKeyCallback)
	if err != nil {
		t.Fatal(err)
	}

	output = client.execCommand("ls -al")
	if output.Std.Err == "could not establish ssh session" {
		t.FailNow()
	}
	// make that panic
	client.Close()
	output = client.execCommand("ls -al")
	if output.Std.Err != "could not establish ssh session" {
		t.FailNow()
	}
}

func TestExecCommands(t *testing.T) {

	demoCommands := commands{
		command("ls -al"),
		command("ls -al"),
	}
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")
	var client *grapeSSHClient
	var err error
	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	hostKeyCallback, err := knownhosts.New("testFiles/known_hosts")
	if err != nil {
		t.FailNow()
	}

	client, err = gSSH.newClient(s, hostKeyCallback)
	if err != nil {
		t.FailNow()
	}
	defer client.Close()
	output := client.execCommands(demoCommands)
	for _, v := range output {
		if v.Std.Err == "could not establish ssh session" {
			t.FailNow()
		}
	}
}
