package testcmd

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Env struct {
	Root    string
	DBPath  string
	MakePkg *MakePkg

	Pacman  string
	Unshare string
	Stdout  io.Writer
	Stderr  io.Writer
}

func NewEnv() (_ *Env, err error) {
	unshare, err := exec.LookPath("unshare")
	if err != nil {
		return nil, err
	}

	pacman, err := exec.LookPath("pacman")
	if err != nil {
		return nil, err
	}

	makepkg, err := NewMakePkg()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = makepkg.Dispose()
		}
	}()

	root, err := CreateTempRoot()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(root)
		}
	}()

	e := &Env{
		Root:    root,
		DBPath:  filepath.Join(root, "var/lib/pacman"),
		MakePkg: makepkg,
		Pacman:  pacman,
		Unshare: unshare,
		Stdout:  nil,
		Stderr:  nil,
	}
	return e, nil
}

func (e *Env) Dispose() error {
	errMakePkg := e.MakePkg.Dispose()
	errRoot := os.RemoveAll(e.Root)
	return errors.Join(errMakePkg, errRoot)
}

func (e *Env) SetStdout(w io.Writer) {
	e.Stdout = w
	e.MakePkg.Stdout = w
}

func (e *Env) SetStderr(w io.Writer) {
	e.Stderr = w
	e.MakePkg.Stderr = w
}

func (e *Env) Install(pkg string, explicit bool) error {
	asopt := "--asdeps"
	if explicit {
		asopt = "--asexplicit"
	}

	cmd := &exec.Cmd{
		Path: e.Unshare,
		Args: []string{
			e.Unshare, "--map-root-user", "--",
			e.Pacman, "--upgrade", "--noconfirm", asopt, "--sysroot", e.Root, "--noprogress", "--", pkg,
		},
		Stdout: e.Stdout,
		Stderr: e.Stderr,
	}
	err := cmd.Run()
	return WrapExitError(cmd, err)
}

func (e *Env) MakeAndInstall(src *PkgBuild, explicit bool) error {
	var err error

	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp) // ignore error

	pkgList, err := e.MakePkg.Run(tmp, src)
	if err != nil {
		return err
	}

	for _, pkg := range pkgList {
		err = e.Install(pkg, explicit)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateTempRoot() (_ string, err error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(dir)
		}
	}()

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(abs, "var/lib/pacman"), 0o755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(abs, "var/cache/pacman/pkg"), 0o755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(abs, "etc"), 0o755)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(
		filepath.Join(abs, "etc/pacman.conf"),
		[]byte("[options]\nSigLevel = Never\nDisableSandbox\n"),
		0o644,
	)
	if err != nil {
		return "", err
	}

	return abs, nil
}
