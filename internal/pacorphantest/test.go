package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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

func init() { testMain.Register("TestNormal", TestNormal) }
func TestNormal(t *testenv.T) {
	env := testenv.HelpEnv(t)

	srcAExp := testenv.NewPkgBuild("a-explicit", "1.1.1")
	srcADep1 := testenv.NewPkgBuild("a-dependency-1", "1.2.1")
	srcADep2 := testenv.NewPkgBuild("a-dependency-2", "1.2.2")
	srcAOpt1 := testenv.NewPkgBuild("a-optional-1", "1.3.1")
	srcAOpt2 := testenv.NewPkgBuild("a-optional-2", "1.3.2")
	srcAExp.Depends = append(srcAExp.Depends, "a-dependency-1")
	srcAExp.OptDepends = append(srcAExp.OptDepends, "a-optional-1")
	srcADep1.Depends = append(srcADep1.Depends, "a-dependency-2")
	srcAOpt1.OptDepends = append(srcAOpt1.OptDepends, "a-optional-2")
	testenv.HelpMakeAndInstall(t, env, srcAExp, true)
	testenv.HelpMakeAndInstall(t, env, srcADep1, false)
	testenv.HelpMakeAndInstall(t, env, srcADep2, false)
	testenv.HelpMakeAndInstall(t, env, srcAOpt1, false)
	testenv.HelpMakeAndInstall(t, env, srcAOpt2, false)

	srcBOrp1 := testenv.NewPkgBuild("b-orphan-1", "2.4.1")
	srcBOrp2 := testenv.NewPkgBuild("b-orphan-2", "2.4.2")
	srcBOrp3 := testenv.NewPkgBuild("b-orphan-3", "2.4.3")
	srcBOrp1.Depends = append(srcBOrp1.Depends, "b-orphan-2")
	srcBOrp1.OptDepends = append(srcBOrp1.OptDepends, "b-orphan-3")
	testenv.HelpMakeAndInstall(t, env, srcBOrp1, false)
	testenv.HelpMakeAndInstall(t, env, srcBOrp2, false)
	testenv.HelpMakeAndInstall(t, env, srcBOrp3, false)

	srcCExp := testenv.NewPkgBuild("c-explicit", "3.1.1")
	srcCDep := testenv.NewPkgBuild("c-dependency", "3.2.1")
	srcCOpt1 := testenv.NewPkgBuild("c-optional-1", "3.3.1")
	srcCOpt2 := testenv.NewPkgBuild("c-optional-2", "3.3.2")
	srcCExp.Depends = append(srcCExp.Depends, "c-dependency")
	srcCDep.OptDepends = append(srcCDep.OptDepends, "c-optional-1")
	srcCOpt1.Depends = append(srcCOpt1.Depends, "c-optional-2")
	testenv.HelpMakeAndInstall(t, env, srcCExp, true)
	testenv.HelpMakeAndInstall(t, env, srcCDep, false)
	testenv.HelpMakeAndInstall(t, env, srcCOpt1, false)
	testenv.HelpMakeAndInstall(t, env, srcCOpt2, false)

	srcDOrp1 := testenv.NewPkgBuild("d-orphan-1", "4.4.1")
	srcDOrp2 := testenv.NewPkgBuild("d-orphan-2", "4.4.2")
	srcDOrp1.Depends = append(srcDOrp1.Depends, "d-orphan-2")
	srcDOrp2.Depends = append(srcDOrp2.Depends, "d-orphan-1")
	testenv.HelpMakeAndInstall(t, env, srcDOrp1, false)
	testenv.HelpMakeAndInstall(t, env, srcDOrp2, false)

	srcEOrp1 := testenv.NewPkgBuild("e-orphan-1", "5.4.1")
	srcEOrp2 := testenv.NewPkgBuild("e-orphan-2", "5.4.2")
	srcEOrp1.OptDepends = append(srcEOrp1.OptDepends, "e-orphan-2")
	srcEOrp2.OptDepends = append(srcEOrp2.OptDepends, "e-orphan-1")
	testenv.HelpMakeAndInstall(t, env, srcEOrp1, false)
	testenv.HelpMakeAndInstall(t, env, srcEOrp2, false)

	srcFExp := testenv.NewPkgBuild("f-explicit", "6.1.1")
	srcFDep1 := testenv.NewPkgBuild("f-dependency-1", "6.2.1")
	srcFDep2 := testenv.NewPkgBuild("f-dependency-2", "6.2.2")
	srcFOpt1 := testenv.NewPkgBuild("f-optional-1", "6.3.1")
	srcFOpt2 := testenv.NewPkgBuild("f-optional-2", "6.3.2")
	srcFExp.Depends = append(srcFExp.Depends, "f-dependency-1")
	srcFExp.OptDepends = append(srcFExp.OptDepends, "f-optional-1")
	srcFDep1.Depends = append(srcFDep1.Depends, "f-dependency-2")
	srcFDep2.Depends = append(srcFDep2.Depends, "f-dependency-1")
	srcFOpt1.OptDepends = append(srcFOpt1.OptDepends, "f-optional-2")
	srcFOpt2.OptDepends = append(srcFOpt2.OptDepends, "f-optional-1")
	testenv.HelpMakeAndInstall(t, env, srcFExp, true)
	testenv.HelpMakeAndInstall(t, env, srcFDep1, false)
	testenv.HelpMakeAndInstall(t, env, srcFDep2, false)
	testenv.HelpMakeAndInstall(t, env, srcFOpt1, false)
	testenv.HelpMakeAndInstall(t, env, srcFOpt2, false)

	cmd1 := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root},
	}

	stdout1 := new(bytes.Buffer)
	stderr1 := new(bytes.Buffer)
	cmd1.Stdout = stdout1
	cmd1.Stderr = stderr1

	err1 := cmd1.Run()
	if err1 != nil {
		t.Error(exiterr.Wrap(cmd1, err1))
	}

	stdout1Want := strings.Join([]string{
		"b-orphan-1 2.4.1-1", "b-orphan-2 2.4.2-1", "b-orphan-3 2.4.3-1",
		"d-orphan-1 4.4.1-1", "d-orphan-2 4.4.2-1",
		"e-orphan-1 5.4.1-1", "e-orphan-2 5.4.2-1",
	}, "\n") + "\n"
	if !bytes.Equal(stdout1.Bytes(), []byte(stdout1Want)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd1.Args, " "), stdout1.Bytes(), stdout1Want)
	}

	if stderr1.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd1.Args, " "), stderr1.Bytes())
	}

	cmd2 := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root, "--strict"},
	}

	stdout2 := new(bytes.Buffer)
	stderr2 := new(bytes.Buffer)
	cmd2.Stdout = stdout2
	cmd2.Stderr = stderr2

	err2 := cmd2.Run()
	if err2 != nil {
		t.Error(exiterr.Wrap(cmd2, err2))
	}

	stdout2Want := strings.Join([]string{
		"a-optional-1 1.3.1-1", "a-optional-2 1.3.2-1",
		"b-orphan-1 2.4.1-1", "b-orphan-2 2.4.2-1", "b-orphan-3 2.4.3-1",
		"c-optional-1 3.3.1-1", "c-optional-2 3.3.2-1",
		"d-orphan-1 4.4.1-1", "d-orphan-2 4.4.2-1",
		"e-orphan-1 5.4.1-1", "e-orphan-2 5.4.2-1",
		"f-optional-1 6.3.1-1", "f-optional-2 6.3.2-1",
	}, "\n") + "\n"
	if !bytes.Equal(stdout2.Bytes(), []byte(stdout2Want)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd2.Args, " "), stdout2.Bytes(), stdout2Want)
	}

	if stderr2.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd2.Args, " "), stderr2.Bytes())
	}

	cmd3 := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root, "--quiet"},
	}

	stdout3 := new(bytes.Buffer)
	stderr3 := new(bytes.Buffer)
	cmd3.Stdout = stdout3
	cmd3.Stderr = stderr3

	err3 := cmd3.Run()
	if err3 != nil {
		t.Error(exiterr.Wrap(cmd3, err3))
	}

	stdout3Want := strings.Join([]string{
		"b-orphan-1", "b-orphan-2", "b-orphan-3",
		"d-orphan-1", "d-orphan-2",
		"e-orphan-1", "e-orphan-2",
	}, "\n") + "\n"
	if !bytes.Equal(stdout3.Bytes(), []byte(stdout3Want)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd3.Args, " "), stdout3.Bytes(), stdout3Want)
	}

	if stderr3.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd3.Args, " "), stderr3.Bytes())
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
	opt.PacmanConf = "/etc/pacman-alt.conf"
	opt.DBPath = "/var/lib/pacman-alt"
	opt.CacheDir = "/var/cache/pacman-alt/pkg"
	env := testenv.HelpEnvWithOption(t, opt) // generates pacman.conf including DBPath

	src := testenv.NewPkgBuild("a", "0.0.1")
	testenv.HelpMakeAndInstall(t, env, src, false) // pacman is executed internally with --sysroot

	cmd1 := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--root", env.Root, "--dbpath", env.DBPath, "--quiet"},
	}

	stdout1 := new(bytes.Buffer)
	stderr1 := new(bytes.Buffer)
	cmd1.Stdout = stdout1
	cmd1.Stderr = stderr1

	err1 := cmd1.Run()
	if err1 != nil {
		t.Error(exiterr.Wrap(cmd1, err1))
	}

	stdout1Want := src.Name + "\n"
	if !bytes.Equal(stdout1.Bytes(), []byte(stdout1Want)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd1.Args, " "), stdout1.Bytes(), stdout1Want)
	}

	if stderr1.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd1.Args, " "), stderr1.Bytes())
	}

	cmd2 := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--sysroot", env.Root, "--config", opt.PacmanConf, "--quiet"},
	}

	stdout2 := new(bytes.Buffer)
	stderr2 := new(bytes.Buffer)
	cmd2.Stdout = stdout2
	cmd2.Stderr = stderr2

	err2 := cmd2.Run()
	if err2 != nil {
		t.Error(exiterr.Wrap(cmd2, err2))
	}

	stdout2Want := src.Name + "\n"
	if !bytes.Equal(stdout2.Bytes(), []byte(stdout2Want)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd2.Args, " "), stdout2.Bytes(), stdout2Want)
	}

	if stderr2.Len() != 0 {
		t.Errorf(`exec [%s]: stderr = %q, want ""`, strings.Join(cmd2.Args, " "), stderr2.Bytes())
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

func init() { testMain.Register("TestVersion", TestVersion) }
func TestVersion(t *testenv.T) {
	cmd := &exec.Cmd{
		Path: pacorphan,
		Args: []string{"pacorphan", "--version"},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		t.Error(exiterr.Wrap(cmd, err))
	}

	const semver = `(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
	stdoutRe := regexp.MustCompile(`^pacorphan ((v(` + semver + `))|(\(devel\))) - libalpm v(` + semver + `)\n$`)
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

func init() { testMain.Register("TestErrorPacmanConfRoot", TestErrorPacmanConfRoot) }
func TestErrorPacmanConfRoot(t *testenv.T) {
	opt := testenv.NewEnvOption()
	env := testenv.HelpEnvWithOption(t, opt)

	errRemove := os.Remove(filepath.Join(env.Root, opt.PacmanConf))
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

	stderrRe := regexp.MustCompile(`^(.|\n)+\nerror: exec \[pacman-conf .+\]: .+\n$`)
	if !stderrRe.Match(stderr.Bytes()) {
		t.Errorf("exec [%s]: stderr = %q, want regexp %q", strings.Join(cmd.Args, " "), stderr.Bytes(), stderrRe.String())
	}
}

func init() { testMain.Register("TestErrorPacmanConfDBPath", TestErrorPacmanConfDBPath) }
func TestErrorPacmanConfDBPath(t *testenv.T) {
	opt := testenv.NewEnvOption()
	env := testenv.HelpEnvWithOption(t, opt)

	f, errOpen := os.OpenFile(
		filepath.Join(env.Root, opt.PacmanConf),
		os.O_WRONLY|os.O_APPEND,
		0o664,
	)
	if errOpen != nil {
		t.Fatal(errOpen)
	}

	_, errWrite := f.WriteString("DBPath = /nonexistent\n")
	if errWrite != nil {
		t.Fatal(errWrite)
	}

	errClose := f.Close()
	if errClose != nil {
		t.Fatal(errClose)
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

	stderrRe := regexp.MustCompile(`^(.|\n)+\nerror: exec \[pacman-conf .+\]: .+\n$`)
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

func init() { testMain.Register("TestWarnMissingDeps", TestWarnMissingDeps) }
func TestWarnMissingDeps(t *testenv.T) {
	env := testenv.HelpEnv(t)

	srcExp := testenv.NewPkgBuild("explicit", "1.1.1")
	srcDep := testenv.NewPkgBuild("dependency", "1.2.1")
	srcOpt := testenv.NewPkgBuild("optional", "1.3.1")
	srcExp.Depends = append(srcExp.Depends, "dependency=2.2.1")
	srcExp.OptDepends = append(srcExp.OptDepends, "optional=2.3.1")
	srcExp.Depends = append(srcExp.Depends, "dependency-alt")
	srcExp.OptDepends = append(srcExp.OptDepends, "optional-alt")
	testenv.HelpMakeAndInstall(t, env, srcExp, true)
	testenv.HelpMakeAndInstall(t, env, srcDep, false)
	testenv.HelpMakeAndInstall(t, env, srcOpt, false)

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

	const stdoutWant = "dependency 1.2.1-1\noptional 1.3.1-1\n"
	if !bytes.Equal(stdout.Bytes(), []byte(stdoutWant)) {
		t.Errorf("exec [%s]: stdout = %q, want %q", strings.Join(cmd.Args, " "), stdout.Bytes(), stdoutWant)
	}

	stderrWant := strings.Join([]string{
		"warning: 'explicit' requires 'dependency-alt', which is not installed",
		"warning: 'explicit' requires 'dependency=2.2.1', but version 1.2.1-1 is installed",
		"warning: 'explicit' recommends 'optional=2.3.1', but version 1.3.1-1 is installed",
	}, "\n") + "\n"
	if !bytes.Equal(stderr.Bytes(), []byte(stderrWant)) {
		t.Errorf("exec [%s]: stderr = %q, want %q", strings.Join(cmd.Args, " "), stderr.Bytes(), stderrWant)
	}
}
