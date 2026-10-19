package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

// The helper exercises the CLI's exit code in a fresh process without changing HOME.
func TestCLIProcess(t *testing.T) {
	if os.Getenv("GRAPES_TEST_CLI") == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				os.Args = append([]string{"grapes"}, os.Args[i+1:]...)
				flag.CommandLine = flag.NewFlagSet("grapes", flag.ExitOnError)
				os.Exit(runCLI(os.Getenv("GRAPES_TEST_KNOWN_HOSTS")))
			}
		}
		t.Fatal("missing CLI arguments")
	}
	for _, async := range []bool{false, true} {
		for _, tc := range []struct {
			name                            string
			command                         command
			badHost, badUser, rejectSession bool
			wantCode                        int
		}{
			{"success with stderr", "warning", false, false, false, 0},
			{"command failure", "silent-failure", false, false, false, 1},
			{"rejected execution", "reject", false, false, false, 1},
			{"rejected session", "success", false, false, true, 1},
			{"connection failure", "success", true, false, false, 1},
			{"authentication failure", "success", false, true, false, 1},
			{"command timeout", "stall", false, false, false, 1},
		} {
			t.Run(fmt.Sprintf("%s/async=%v", tc.name, async), func(t *testing.T) {
				f := newSSHFixture(t, tc.rejectSession)
				target := f.server
				if tc.badHost {
					target.Host = "missing-port"
				}
				if tc.badUser {
					target.User = "unauthorized"
				}
				conf := config{Version: "1", Servers: map[string]servers{"test": {target}}, Commands: map[string]commands{"test": {tc.command}}}
				data, err := yaml.Marshal(conf)
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(t.TempDir(), "config.yml")
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
				args := []string{"-test.run=^TestCLIProcess$", "--", "-y", "-c", path, "-i", "testFiles/id_rsa", "-s", "test", "-cmd", "test"}
				if tc.command == "stall" {
					args = append(args, "-command-timeout", "100ms")
				}
				if async {
					args = append(args, "-async")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, os.Args[0], args...)
				cmd.Env = append(os.Environ(), "GRAPES_TEST_CLI=1", "GRAPES_TEST_KNOWN_HOSTS="+f.knownHosts)
				var stdout, stderr bytes.Buffer
				cmd.Stdout, cmd.Stderr = &stdout, &stderr
				err = cmd.Run()
				if ctx.Err() != nil || cmd.ProcessState == nil {
					t.Fatalf("CLI did not finish: %v; stderr: %s", err, &stderr)
				}
				if got := cmd.ProcessState.ExitCode(); got != tc.wantCode {
					t.Fatalf("exit code = %d, want %d; stdout: %s; stderr: %s", got, tc.wantCode, &stdout, &stderr)
				}
				if !strings.Contains(stdout.String(), "host: "+target.Host) {
					t.Fatalf("CLI exited without printing its server result: %s", &stdout)
				}
				if tc.wantCode != 0 && stderr.Len() == 0 {
					t.Fatal("CLI failure did not include a diagnostic")
				}
				if tc.command == "stall" && !strings.Contains(stdout.String(), "timed out after 100ms") {
					t.Fatalf("configured timeout missing from CLI output: %s", &stdout)
				}
				if tc.name == "command failure" && !strings.Contains(stdout.String(), "exit_status: 42") {
					t.Fatalf("command exit status missing from CLI output: %s", &stdout)
				}
			})
		}
	}
}
