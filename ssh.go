package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultCommandTimeout = 5 * time.Minute
)

type (
	grapeSSH struct {
		keySigner      ssh.Signer
		connectTimeout time.Duration
		commandTimeout time.Duration
	}
	grapeSSHClient struct {
		*ssh.Client
		conn           net.Conn
		commandTimeout time.Duration
		timedOut       bool
	}
	std struct {
		Out string
		Err string
	}
	sshOutput struct {
		Command    command
		Std        std
		Error      string `yaml:"error,omitempty"`
		ExitStatus *int   `yaml:"exit_status,omitempty"`
	}
	sshOutputArray []*sshOutput
	sshError       error
)

func (gSSH *grapeSSH) newError(errMsg string) sshError {
	return errors.New(errMsg)
}

func (gSSH *grapeSSH) setKey(keyPath keyPath) sshError {
	privateBytes, err := os.ReadFile(string(keyPath))
	if err != nil {
		return gSSH.newError("Could not open identity file.")
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return gSSH.newError(fmt.Sprint("Could not parse identity file."))
	}
	gSSH.keySigner = privateKey
	return nil
}

func (gSSH *grapeSSH) newClient(server server, hostKeyCallback ssh.HostKeyCallback) (*grapeSSHClient, sshError) {
	timeout := gSSH.connectTimeout
	if timeout == 0 {
		timeout = defaultConnectTimeout
	}
	deadline := time.Now().Add(timeout)
	config := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(gSSH.keySigner),
		},
		HostKeyCallback: hostKeyCallback,
	}
	raw, err := net.DialTimeout("tcp", server.Host, timeout)
	if err != nil {
		return nil, fmt.Errorf("could not establish ssh connection: %w", err)
	}
	if err := raw.SetDeadline(deadline); err != nil {
		raw.Close()
		return nil, err
	}
	conn, channels, requests, err := ssh.NewClientConn(raw, server.Host, config)
	if err != nil {
		raw.Close()
		return nil, fmt.Errorf("could not establish ssh connection: %w", err)
	}
	raw.SetDeadline(time.Time{})
	return &grapeSSHClient{Client: ssh.NewClient(conn, channels, requests), conn: raw, commandTimeout: gSSH.commandTimeout}, nil
}

func (client *grapeSSHClient) execCommand(cmd command) *sshOutput {
	output := &sshOutput{
		Command: cmd,
	}
	if client.timedOut {
		output.Error = "command skipped after previous command timed out"
		return output
	}
	timeout := client.commandTimeout
	if timeout == 0 {
		timeout = defaultCommandTimeout
	}
	deadline := time.Now().Add(timeout)
	if err := client.conn.SetDeadline(deadline); err != nil {
		output.Error = err.Error()
		return output
	}
	defer func() {
		client.conn.SetDeadline(time.Time{})
		if !time.Now().Before(deadline) {
			client.Close()
			client.timedOut = true
			output.Error = fmt.Sprintf("command timed out after %s", timeout)
			output.ExitStatus = nil
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		output.Error = fmt.Sprintf("could not establish ssh session: %v", err)
		return output
	}
	defer session.Close()
	var stderr, stdout bytes.Buffer
	session.Stdout, session.Stderr = &stdout, &stderr
	err = session.Run(string(cmd))
	output.Std = std{Out: stdout.String(), Err: stderr.String()}
	if err == nil {
		status := 0
		output.ExitStatus = &status
	} else {
		output.Error = err.Error()
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			status := exitErr.ExitStatus()
			output.ExitStatus = &status
		}
	}
	return output
}

func (client *grapeSSHClient) execCommands(commands commands) sshOutputArray {
	output := sshOutputArray{}
	for _, command := range commands {
		output = append(output, client.execCommand(command))
	}
	return output
}
