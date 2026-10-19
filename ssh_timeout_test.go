package main

import (
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestConnectionTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		// Accept TCP, but never send the SSH banner.
		io.Copy(io.Discard, conn)
	}()
	client := grapeSSH{connectTimeout: 100 * time.Millisecond}
	started := time.Now()
	_, err = client.newClient(server{Host: listener.Addr().String()}, ssh.InsecureIgnoreHostKey())
	var timeout net.Error
	if !errors.As(err, &timeout) || !timeout.Timeout() {
		t.Fatalf("expected handshake timeout, got %v", err)
	}
	if time.Since(started) > 2*time.Second {
		t.Fatal("connection timeout was not enforced")
	}
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("timed-out connection was not closed")
	}
}

func TestCommandTimeout(t *testing.T) {
	f := newSSHFixture(t, false)
	f.ssh.commandTimeout = 100 * time.Millisecond
	client, err := f.ssh.newClient(f.server, f.hostKey)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	started := time.Now()
	output := client.execCommands(commands{"stall", "success"})
	if time.Since(started) > 2*time.Second {
		t.Fatal("command timeout was not enforced")
	}
	if len(output) != 2 || !strings.Contains(output[0].Error, "timed out after 100ms") || output[0].ExitStatus != nil {
		t.Fatalf("timeout not reported: %+v", output)
	}
	if output[0].Std.Out != "started\n" {
		t.Fatalf("lost output before timeout: %+v", output[0])
	}
	if !strings.Contains(output[1].Error, "skipped") || output[1].ExitStatus != nil {
		t.Fatalf("subsequent command was not skipped: %+v", output[1])
	}
	select {
	case <-f.closed:
	case <-time.After(time.Second):
		t.Fatal("timeout did not close the SSH connection")
	}
}
