package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type sshFixture struct {
	server     server
	ssh        grapeSSH
	knownHosts string
	hostKey    ssh.HostKeyCallback
	closed     chan struct{}
}

// Each fixture serves one authenticated connection, with multiple command sessions.
func newSSHFixture(t *testing.T, rejectSessions bool) *sshFixture {
	t.Helper()
	f := &sshFixture{closed: make(chan struct{})}
	if err := f.ssh.setKey("testFiles/id_rsa"); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f.server = server{Name: "local SSH", Host: listener.Addr().String(), User: "tester"}
	f.knownHosts = filepath.Join(t.TempDir(), "known_hosts")
	line := knownhosts.Line([]string{f.server.Host}, f.ssh.keySigner.PublicKey()) + "\n"
	if err := os.WriteFile(f.knownHosts, []byte(line), 0600); err != nil {
		listener.Close()
		t.Fatal(err)
	}
	f.hostKey, err = knownhosts.New(f.knownHosts)
	if err != nil {
		listener.Close()
		t.Fatal(err)
	}
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(meta ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if meta.User() != f.server.User || !bytes.Equal(key.Marshal(), f.ssh.keySigner.PublicKey().Marshal()) {
				return nil, fmt.Errorf("unauthorized test client")
			}
			return nil, nil
		},
	}
	config.AddHostKey(f.ssh.keySigner)
	var mu sync.Mutex
	var conn net.Conn
	stopped := false
	t.Cleanup(func() {
		listener.Close()
		mu.Lock()
		stopped = true
		if conn != nil {
			conn.Close()
		}
		mu.Unlock()
		<-f.closed
	})
	go func() {
		defer close(f.closed)
		raw, err := listener.Accept()
		if err != nil {
			return
		}
		defer raw.Close()
		mu.Lock()
		conn = raw
		if stopped {
			mu.Unlock()
			return
		}
		mu.Unlock()
		raw.SetDeadline(time.Now().Add(10 * time.Second))
		transport, channels, requests, err := ssh.NewServerConn(raw, config)
		if err != nil {
			return
		}
		defer transport.Close()
		go ssh.DiscardRequests(requests)
		for incoming := range channels {
			if rejectSessions || incoming.ChannelType() != "session" {
				incoming.Reject(ssh.Prohibited, "sessions disabled")
				continue
			}
			channel, requests, err := incoming.Accept()
			if err != nil {
				return
			}
			serveTestSession(transport, channel, requests)
		}
	}()
	return f
}

func serveTestSession(transport *ssh.ServerConn, channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()
	for request := range requests {
		var exec struct{ Command string }
		if request.Type != "exec" || ssh.Unmarshal(request.Payload, &exec) != nil {
			request.Reply(false, nil)
			continue
		}
		if exec.Command == "reject" {
			request.Reply(false, nil)
			return
		}
		request.Reply(true, nil)
		var status uint32
		switch exec.Command {
		case "success":
			io.WriteString(channel, "ok\n")
		case "warning":
			io.WriteString(channel, "ok\n")
			io.WriteString(channel.Stderr(), "warning\n")
		case "failure":
			io.WriteString(channel, "partial output\n")
			io.WriteString(channel.Stderr(), "failed\n")
			status = 42
		case "stall":
			io.WriteString(channel, "started\n")
			transport.Wait()
			return
		case "exit-one":
			status = 1
		case "silent-failure":
			status = 42
		case "disconnect":
			transport.Close()
			return
		case "missing-status":
			return
		default:
			status = 127
		}
		channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{status}))
		return
	}
}
