package main

import (
	"testing"
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
	if _, err := gSSH.newClient(s); err != nil {
		t.FailNow()
	}
	s.Host = "localhost"
	if _, err := gSSH.newClient(s); err == nil {
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
	client, err = gSSH.newClient(s)
	if err != nil {
		t.FailNow()
	}
	output = client.execCommand("echo")
	if output.Std.Err == "could not establish ssh session" {
		t.FailNow()
	}
	// make that panic
	client.Close()
	output = client.execCommand("echo")
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
	client, err = gSSH.newClient(s)
	defer client.Close()
	if err != nil {
		t.FailNow()
	}
	output := client.execCommands(demoCommands)
	for _, v := range output {
		if v.Std.Err == "could not establish ssh session" {
			t.FailNow()
		}
	}
}
