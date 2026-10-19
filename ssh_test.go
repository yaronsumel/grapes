package main

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestSetKey(t *testing.T) {
	gSSH := grapeSSH{}
	if err := gSSH.setKey("testFiles/id_rsa"); err != nil {
		t.Fatal(err)
	}
	if err := gSSH.setKey("testFiles/id_rsa_123"); err == nil {
		t.Fatal("accepted missing private key")
	}
	if err := gSSH.setKey("testFiles/id_rsa.pub"); err == nil {
		t.Fatal("accepted public key as private key")
	}
}

func TestNewClient(t *testing.T) {
	f := newSSHFixture(t, false)
	client, err := f.ssh.newClient(f.server, f.hostKey)
	if err != nil {
		t.Fatal(err)
	}
	client.Close()
	if _, err := f.ssh.newClient(server{Host: "missing-port"}, f.hostKey); err == nil {
		t.Fatal("accepted invalid server address")
	}
}

func TestExecCommand(t *testing.T) {
	for _, tc := range []struct {
		command        command
		stdout, stderr string
		status         int
		wantError      bool
	}{
		{"success", "ok\n", "", 0, false},
		{"warning", "ok\n", "warning\n", 0, false},
		{"failure", "partial output\n", "failed\n", 42, true},
		{"silent-failure", "", "", 42, true},
		{"reject", "", "", -1, true},
		{"disconnect", "", "", -1, true},
		{"missing-status", "", "", -1, true},
	} {
		t.Run(string(tc.command), func(t *testing.T) {
			f := newSSHFixture(t, false)
			client, err := f.ssh.newClient(f.server, f.hostKey)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			output := client.execCommand(tc.command)
			if output.Command != tc.command || output.Std.Out != tc.stdout || output.Std.Err != tc.stderr {
				t.Fatalf("unexpected command output: %+v", output)
			}
			if (output.Error != "") != tc.wantError {
				t.Fatalf("error = %q, want error: %v", output.Error, tc.wantError)
			}
			if tc.status < 0 {
				if output.ExitStatus != nil {
					t.Fatalf("invented exit status: %d", *output.ExitStatus)
				}
			} else if output.ExitStatus == nil || *output.ExitStatus != tc.status {
				t.Fatalf("exit status = %v, want %d", output.ExitStatus, tc.status)
			}
			data, err := yaml.Marshal(output)
			if err != nil {
				t.Fatal(err)
			}
			var decoded sshOutput
			if err := yaml.Unmarshal(data, &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded.Error != output.Error || (decoded.ExitStatus == nil) != (output.ExitStatus == nil) {
				t.Fatalf("error or status missing from YAML: %s", data)
			}
			if tc.status == 0 && !strings.Contains(string(data), "exit_status: 0") {
				t.Fatalf("successful exit status omitted from YAML: %s", data)
			}
		})
	}
}

func TestExecCommandSessionFailure(t *testing.T) {
	f := newSSHFixture(t, true)
	client, err := f.ssh.newClient(f.server, f.hostKey)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	output := client.execCommand("success")
	if !strings.Contains(output.Error, "could not establish ssh session") || output.ExitStatus != nil {
		t.Fatalf("session failure not preserved: %+v", output)
	}
}

func TestExecCommands(t *testing.T) {
	f := newSSHFixture(t, false)
	client, err := f.ssh.newClient(f.server, f.hostKey)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	output := client.execCommands(commands{"silent-failure", "success"})
	if len(output) != 2 || output[0].Error == "" || output[1].Error != "" || output[1].Std.Out != "ok\n" {
		t.Fatalf("commands were not all executed and reported: %+v", output)
	}
}
