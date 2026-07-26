package testenv

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gcrtnst/pacorphan/internal/exiterr"
)

const (
	defaultPacmanConf = "/etc/pacman.conf"
	defaultDBPath     = "/var/lib/pacman"
	defaultCacheDir   = "/var/cache/pacman/pkg"
)

type EnvOption struct {
	PacmanConf string
	DBPath     string
	CacheDir   string
}

func NewEnvOption() *EnvOption {
	return &EnvOption{
		PacmanConf: defaultPacmanConf,
		DBPath:     defaultDBPath,
		CacheDir:   defaultCacheDir,
	}
}

type Env struct {
	Root    string
	DBPath  string
	MakePkg *MakePkg

	Pacman  string
	Unshare string
	Stdout  io.Writer
	Stderr  io.Writer
}

func NewEnv(opt *EnvOption) (_ *Env, err error) {
	if opt == nil {
		opt = NewEnvOption()
	}

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

	root, err := CreateTempRoot(opt)
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
		DBPath:  filepath.Join(root, opt.DBPath),
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
	return exiterr.Wrap(cmd, err)
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

func CreateTempRoot(opt *EnvOption) (_ string, err error) {
	if opt == nil {
		opt = NewEnvOption()
	}
	// Assume that all file paths in opt are absolute paths and
	// do not contain special characters such as newline characters.

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

	err = os.MkdirAll(filepath.Join(abs, opt.PacmanConf, ".."), 0o755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(abs, opt.DBPath), 0o755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(abs, opt.CacheDir), 0o755)
	if err != nil {
		return "", err
	}

	conf := new(bytes.Buffer)
	_, _ = conf.WriteString("[options]\nSigLevel = Never\nDisableSandbox\n")
	if opt.DBPath != defaultDBPath {
		_, _ = conf.WriteString("DBPath = ")
		_, _ = conf.WriteString(opt.DBPath)
		_, _ = conf.WriteRune('\n')
	}
	if opt.CacheDir != defaultCacheDir {
		_, _ = conf.WriteString("CacheDir = ")
		_, _ = conf.WriteString(opt.CacheDir)
		_, _ = conf.WriteRune('\n')
	}
	err = os.WriteFile(filepath.Join(abs, opt.PacmanConf), conf.Bytes(), 0o644)
	if err != nil {
		return "", err
	}

	return abs, nil
}
