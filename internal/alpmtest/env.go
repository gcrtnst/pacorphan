package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

type Env struct {
	Root    string
	MakePkg *MakePkg

	Pacman  string
	Unshare string
	Stdout  *os.File
	Stderr  *os.File
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
		MakePkg: makepkg,
		Pacman:  pacman,
		Unshare: unshare,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	return e, nil
}

func (r *Env) Dispose() error {
	errMakePkg := r.MakePkg.Dispose()
	errRoot := os.RemoveAll(r.Root)
	return errors.Join(errMakePkg, errRoot)
}

func (r *Env) Install(pkg string, explicit bool) error {
	asopt := "--asdeps"
	if explicit {
		asopt = "--asexplicit"
	}

	cmd := &exec.Cmd{
		Path: r.Unshare,
		Args: []string{
			r.Unshare, "--map-root-user", "--",
			r.Pacman, "--upgrade", "--noconfirm", asopt, "--sysroot", r.Root, "--noprogress", "--", pkg,
		},
		Stdout: r.Stdout,
		Stderr: r.Stderr,
	}
	err := cmd.Run()
	return WrapExitError(cmd, err)
}

func (r *Env) MakeAndInstall(src *PkgBuild, explicit bool) error {
	var err error

	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp) // ignore error

	pkgList, err := r.MakePkg.Run(tmp, src)
	if err != nil {
		return err
	}

	for _, pkg := range pkgList {
		err = r.Install(pkg, explicit)
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
