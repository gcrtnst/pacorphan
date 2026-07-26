package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gcrtnst/pacorphan/internal/exiterr"
	"github.com/gcrtnst/pacorphan/internal/testenv"
)

func init() { testMain.Register("TestEmpty", TestEmpty) }
func TestEmpty(t *testenv.T) {
	env := testenv.HelpEnv(t)

	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		t.Error(exiterr.Wrap(cmd, err))
	}

	if stdout.Len() != 0 {
		t.Errorf(`exec [%s]: stdout = %q, want ""`, strings.Join(cmd.Args, " "), stdout.Bytes())
	}

	if stderr.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd.Args, " "), stderr.Bytes())
	}
}

func init() { testMain.Register("TestCustomPath", TestCustomPath) }
func TestCustomPath(t *testenv.T) {
	// When --sysroot is used, if the DBPath defined in the sysroot's pacman.conf
	// does not exist on the host system (outside the sysroot), the following error
	// occurs even though the path exists inside the sysroot:
	// error: 'failed to resolve path '/var/lib/pacman-alt/' passed to 'DBPath': No such file or directory
	//
	// Since this is considered a bug in pacman, this test is skipped until it is fixed.
	t.Skip("skipping due to pacman bug with --sysroot")

	opt := testenv.NewEnvOption()
	opt.DBPath = "/var/lib/pacman-alt"
	opt.CacheDir = "/var/cache/pacman-alt/pkg"
	env := testenv.HelpEnvWithOption(t, opt)

	src := testenv.NewPkgBuild("a", "0.0.1")
	testenv.HelpMakeAndInstall(t, env, src, false)

	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--root", env.Root, "--dbpath", env.DBPath, "--quiet"},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		t.Error(exiterr.Wrap(cmd, err))
	}

	stdoutWant := src.Name + "\n"
	if !bytes.Equal(stdout.Bytes(), []byte(stdoutWant)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd.Args, " "), stdout.Bytes(), stdoutWant)
	}

	if stderr.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd.Args, " "), stderr.Bytes())
	}
}

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
		t.Error(exiterr.Wrap(cmd, err))
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
		t.Error(exiterr.Wrap(cmd, err))
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
		t.Error(exiterr.Wrap(cmd, err))
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
