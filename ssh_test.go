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

func TestNewSession(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")
	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}
	client, err := gSSH.newClient(s)
	if err != nil {
		t.FailNow()
	}
	if _, err := client.newSession(); err != nil {
		t.FailNow()
	}
	client.Close()
	if _, err := client.newSession(); err == nil {
		t.FailNow()
	}
}

func TestExec(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")
	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}
	client, err := gSSH.newClient(s)
	if err != nil {
		t.FailNow()
	}
	std, err := client.exec(command("echo"))
	if err != nil {
		t.FailNow()
	}
	client.Close()
	_, err2 := client.exec(command("echo"))
	if err2 == nil {
		t.FailNow()
	}
	if std.Std.Err == "" && std.Std.Out == "" {
		t.FailNow()
	}
}