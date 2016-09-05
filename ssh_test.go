package main

import (
	"testing"
)

func TestSetKey(t *testing.T) {

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		gSSH := grapeSSH{}
		gSSH.setKey(keyPath("testFiles/id_rsa"))
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		gSSH := grapeSSH{}
		gSSH.setKey(keyPath("testFiles/id_rsa_123"))
		t.FailNow()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		gSSH := grapeSSH{}
		gSSH.setKey(keyPath("testFiles/id_rsa.pub"))
		t.FailNow()
	}()

}

func TestNewClient(t *testing.T) {

	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")

	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		gSSH.newClient(s)
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		s.Host = "localhost"
		gSSH.newClient(s)
		t.FailNow()
	}()
}

func TestNewSession(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")

	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	client := gSSH.newClient(s)

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		client.newSession()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		client.Close()
		client.newSession()
		t.FailNow()
	}()

}

func TestExec(t *testing.T) {
	gSSH := grapeSSH{}
	gSSH.setKey("testFiles/id_rsa")

	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	client := gSSH.newClient(s)

	std := client.exec(command("echo"))

	if std.Command != command("echo") {
		t.FailNow()
	}

	if std.Std.Err == "" && std.Std.Out == "" {
		t.FailNow()
	}

}
