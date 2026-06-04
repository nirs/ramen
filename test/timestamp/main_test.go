// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var linePrefix = regexp.MustCompile(`^\[\s*\d+\.\d{3}\] `)

type result struct {
	exit   int
	stdout string
	stderr string
}

func TestStdout(t *testing.T) {
	r := run(t, "./timestamp", "echo", "hello")
	checkExit(t, r.exit, 0)
	checkStream(t, r.stdout, "hello")
	checkStream(t, r.stderr)
}

func TestStderr(t *testing.T) {
	r := run(t, "./timestamp", "sh", "-c", "echo err >&2")
	checkExit(t, r.exit, 0)
	checkStream(t, r.stdout)
	checkStream(t, r.stderr, "err")
}

func TestStdoutAndStderr(t *testing.T) {
	r := run(t, "./timestamp", "sh", "-c", "echo out; echo err >&2")
	checkExit(t, r.exit, 0)
	checkStream(t, r.stdout, "out")
	checkStream(t, r.stderr, "err")
}

func TestExitFailure(t *testing.T) {
	r := run(t, "./timestamp", "false")
	checkExit(t, r.exit, 1)
	checkStream(t, r.stdout)
	checkStream(t, r.stderr)
}

func TestNoCommand(t *testing.T) {
	r := run(t, "./timestamp")
	checkExit(t, r.exit, 1)
	checkStream(t, r.stdout)
	checkStderrText(t, r.stderr, "No command provided")
}

func TestMissingCommand(t *testing.T) {
	r := run(t, "./timestamp", "timestamp-test-nonexistent")
	checkExit(t, r.exit, 127)
	checkStream(t, r.stdout)
	checkStderrText(t, r.stderr, "executable file not found")
}

func run(t *testing.T, program string, args ...string) result {
	t.Helper()

	cmd := exec.Command(program, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	r := result{stdout: stdout.String(), stderr: stderr.String()}
	if err == nil {
		return r
	}

	exit, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run %s %v: %v", program, args, err)
	}

	r.exit = exit.ExitCode()
	return r
}

func checkExit(t *testing.T, exit, want int) {
	t.Helper()

	if exit != want {
		t.Fatalf("exit = %d, want %d", exit, want)
	}
}

func checkStderrText(t *testing.T, stderr, want string) {
	t.Helper()

	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
}

func checkStream(t *testing.T, out string, want ...string) {
	t.Helper()

	if len(want) == 0 {
		if out != "" {
			t.Fatalf("got %q, want empty", out)
		}
		return
	}

	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != len(want) {
		t.Fatalf("got %q, want %v", out, want)
	}

	for i, line := range lines {
		if !linePrefix.MatchString(line) {
			t.Fatalf("line %d missing timestamp: %q", i, line)
		}
		idx := strings.Index(line, "] ")
		if got := line[idx+2:]; got != want[i] {
			t.Fatalf("line %d: got %q, want %q", i, got, want[i])
		}
	}
}
