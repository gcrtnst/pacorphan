package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gcrtnst/pacorphan/internal/testenv"
)

func init() { testMain.Register("TestHelp", TestHelp) }
func TestHelp(t *testenv.T) {
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
		t.Error(testenv.WrapExitError(cmd, err))
	}

	stdoutRe := regexp.MustCompile(`(?s)^Usage of pacorphan:\n.+$`)
	if !stdoutRe.Match(stdout.Bytes()) {
		t.Errorf("exec [%s]: stdout = %q, want regexp %q", strings.Join(cmd.Args, " "), stdout.Bytes(), stdoutRe.String())
	}

	if stderr.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd.Args, " "), stderr.Bytes())
	}
}

func init() { testMain.Register("TestErrorFlagParse", TestErrorFlagParse) }
func TestErrorFlagParse(t *testenv.T) {
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
		t.Error(testenv.WrapExitError(cmd, err))
	} else {
		t.Errorf("exec [%s]: exit code = 0, want %d", strings.Join(cmd.Args, " "), wantCode)
	}

	if stdout.Len() != 0 {
		t.Errorf(`exec [%s]: stdout = %q, want ""`, strings.Join(cmd.Args, " "), stdout.Bytes())
	}

	stderrRe := regexp.MustCompile(`^error: .+\n$`)
	if !stderrRe.Match(stderr.Bytes()) {
		t.Errorf("exec [%s]: stderr = %q, want regexp %q", strings.Join(cmd.Args, " "), stderr.Bytes(), stderrRe.String())
	}
}

func init() { testMain.Register("TestErrorALPMInit", TestErrorALPMInit) }
func TestErrorALPMInit(t *testenv.T) {
	env := testenv.HelpEnv(t)

	errRemove := os.RemoveAll(env.DBPath)
	if errRemove != nil {
		t.Fatal(errRemove)
	}

	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	const wantCode = 1
	if e, ok := errors.AsType[*exec.ExitError](err); ok {
		if e.ExitCode() != wantCode {
			t.Errorf("exec [%s]: exit code = %d, want %d", strings.Join(cmd.Args, " "), e.ExitCode(), wantCode)
		}
	} else if err != nil {
		t.Error(testenv.WrapExitError(cmd, err))
	} else {
		t.Errorf("exec [%s]: exit code = 0, want %d", strings.Join(cmd.Args, " "), wantCode)
	}

	if stdout.Len() != 0 {
		t.Errorf(`exec [%s]: stdout = %q, want ""`, strings.Join(cmd.Args, " "), stdout.Bytes())
	}

	stderrRe := regexp.MustCompile(`^error: failed to initialize alpm library: .+\n$`)
	if !stderrRe.Match(stderr.Bytes()) {
		t.Errorf("exec [%s]: stderr = %q, want regexp %q", strings.Join(cmd.Args, " "), stderr.Bytes(), stderrRe.String())
	}
}
