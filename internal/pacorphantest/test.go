package main

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gcrtnst/pacorphan/internal/testcmd"
)

func init() { testMain.Register("TestFlagParseError", TestFlagParseError) }
func TestFlagParseError(t *testcmd.T) {
	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--invalid-option"},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	const wantCode = 2
	if e, ok := errors.AsType[*exec.ExitError](err); ok {
		if e.ExitCode() != wantCode {
			t.Errorf("exec [%s]: exit code = %d, want %d", strings.Join(cmd.Args, " "), e.ExitCode(), wantCode)
		}
	} else if err != nil {
		t.Error(testcmd.WrapExitError(cmd, err))
	} else {
		t.Errorf("exec [%s]: exit code = 0, want %d", strings.Join(cmd.Args, " "), wantCode)
	}

	if stdout.Len() != 0 {
		t.Errorf("exec [%s]: stdout.Len() = %d, want 0", strings.Join(cmd.Args, " "), stdout.Len())
	}

	stderrRe := regexp.MustCompile(`^error: .+\n$`)
	if !stderrRe.Match(stderr.Bytes()) {
		t.Errorf("exec [%s]: stderr.Bytes() = %q, want regexp %q", strings.Join(cmd.Args, " "), stderr.Bytes(), stderrRe.String())
	}
}

func init() { testMain.Register("TestHelp", TestHelp) }
func TestHelp(t *testcmd.T) {
	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--help"},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		t.Error(testcmd.WrapExitError(cmd, err))
	}

	stdoutRe := regexp.MustCompile(`(?s)^Usage of pacorphan:\n.+$`)
	if !stdoutRe.Match(stdout.Bytes()) {
		t.Errorf("exec [%s]: stdout.Bytes() = %q, want regexp %q", strings.Join(cmd.Args, " "), stdout.Bytes(), stdoutRe.String())
	}

	if stderr.Len() != 0 {
		t.Errorf("exec [%s]: stderr.Len() = %d, want 0", strings.Join(cmd.Args, " "), stderr.Len())
	}
}
