package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewGrape(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "test",
			configPath:  "test",
			keyPath:     "test",
			serverGroup: "test",
		})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "staging",
		})
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "not_working",
		})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "not_working",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "staging",
		})
		t.FailNow()
	}()

	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa.pub",
			serverGroup: "staging",
		})
		t.FailNow()
	}()

}

func TestRun(t *testing.T) {
	for _, async := range []bool{false, true} {
		for _, tc := range []struct {
			name      string
			command   command
			badHost   bool
			wantError bool
		}{
			{"success", "warning", false, false},
			{"command failure", "silent-failure", false, true},
			{"mixed servers", "success", true, true},
		} {
			t.Run(fmt.Sprintf("%s/async=%v", tc.name, async), func(t *testing.T) {
				f := newSSHFixture(t, false)
				app := grape{ssh: f.ssh, input: input{asyncFlag: asyncFlag(async)}, servers: servers{f.server}, commands: commands{tc.command}}
				if tc.badHost {
					app.servers = append(servers{{Host: "missing-port"}}, app.servers...)
				}
				var out bytes.Buffer
				err := app.run(f.knownHosts, &out)
				if (err != nil) != tc.wantError {
					t.Fatalf("run error = %v, want error: %v", err, tc.wantError)
				}
				if strings.Count(out.String(), "host:") != len(app.servers) || !strings.Contains(out.String(), string(tc.command)) {
					t.Fatalf("missing server results: %s", &out)
				}
				if tc.badHost && !strings.Contains(err.Error(), "missing-port") {
					t.Fatalf("connection failure missing from aggregate error: %v", err)
				}
				if tc.command == "silent-failure" && !strings.Contains(err.Error(), "42") {
					t.Fatalf("command failure missing from aggregate error: %v", err)
				}
				select {
				case <-f.closed:
				case <-time.After(2 * time.Second):
					t.Fatal("SSH client was not closed after running commands")
				}
			})
		}
	}
}

type blockedOutput struct {
	entered chan struct{}
	release chan struct{}
	bytes.Buffer
}

func (w *blockedOutput) Write(data []byte) (int, error) {
	close(w.entered)
	<-w.release
	return w.Buffer.Write(data)
}

func TestAsyncRunWaitsForOutput(t *testing.T) {
	writer := &blockedOutput{entered: make(chan struct{}), release: make(chan struct{})}
	app := grape{input: input{asyncFlag: true}, servers: servers{{Host: "missing-port"}}}
	done := make(chan error, 1)
	go func() { done <- app.run("testFiles/known_hosts", writer) }()
	select {
	case <-writer.entered:
	case <-time.After(2 * time.Second):
		close(writer.release)
		t.Fatal("worker did not start writing output")
	}
	returnedEarly := false
	select {
	case <-done:
		returnedEarly = true
	case <-time.After(100 * time.Millisecond):
	}
	close(writer.release)
	if returnedEarly {
		t.Fatal("run returned before its output finished")
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("connection failure was not returned")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run did not finish after output was released")
	}
	if !strings.Contains(writer.String(), "missing-port") {
		t.Fatal("server result was not written")
	}
}

type failedOutput struct{ err error }

func (w failedOutput) Write([]byte) (int, error) { return 0, w.err }

func TestRunOutputFailure(t *testing.T) {
	want := errors.New("output unavailable")
	app := grape{servers: servers{{Host: "missing-port"}}}
	if err := app.run("testFiles/known_hosts", failedOutput{want}); !errors.Is(err, want) {
		t.Fatalf("output error was not returned: %v", err)
	}
}

func TestRunMissingKnownHosts(t *testing.T) {
	app := grape{}
	if err := app.run(filepath.Join(t.TempDir(), "missing"), io.Discard); err == nil {
		t.Fatal("missing known_hosts did not return an error")
	}
}

func TestPrintOutput(t *testing.T) {
	o := sshOutput{
		Command: "a",
		Std: std{
			Err: "err",
			Out: "out",
		},
	}
	s := server{
		User:  "a",
		Name:  "b",
		Host:  "c",
		Fatal: "d",
		Output: sshOutputArray{
			&o,
		},
	}
	var out bytes.Buffer
	if err := s.printOutput(&out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "command: a") {
		t.Fatalf("missing command output: %s", &out)
	}
}

func TestVerifyAction(t *testing.T) {

	app := grape{
		commands: []command{
			command("ls -al #1"),
			command("ls -al #2"),
		},
		servers: []server{
			{
				Host:  "a",
				Name:  "a",
				User:  "a",
				Fatal: "a",
			},
		},
		input: input{
			verifyFlag:  true,
			commandName: "a",
		},
	}

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		app.verifyAction()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		app.input.verifyFlag = false
		app.verifyAction()
		t.FailNow()
	}()

}
