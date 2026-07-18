package main

import (
	"os"
	"path/filepath"

	"github.com/gcrtnst/pacorphan/internal/alpm"
)

func init() { Register("TestHandle", TestHandle) }
func TestHandle(t *T) {
	env := HelpEnv(t)

	h, err := alpm.NewHandle(env.Root, env.DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := h.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	root := h.Root()
	if !samefile(root, env.Root) {
		t.Errorf("h.Root() = %q, want %q", root, env.Root)
	}

	dbpath := h.DBPath()
	if !samefile(dbpath, env.DBPath) {
		t.Errorf("h.DBPath() = %q, want %q", dbpath, env.DBPath)
	}
}

func init() { Register("TestHandleInitError", TestHandleInitError) }
func TestHandleInitError(t *T) {
	root, errMkdir := os.MkdirTemp("", "")
	if errMkdir != nil {
		t.Error(errMkdir)
		return
	}
	defer os.RemoveAll(root) // ignore error

	dbpath := filepath.Join(root, "var/lib/pacman")
	// dbpath is not created

	h, err := alpm.NewHandle(root, dbpath)
	if err == nil {
		_ = h.Close()
		t.Fatalf(`alpm.NewHandle(%q, %q) error = nil, want non-nil`, root, dbpath)
	}

	e, ok := err.(*alpm.Error)
	if !ok {
		t.Fatalf(`alpm.NewHandle(%q, %q) error type = %T, want %T`, root, dbpath, err, (*alpm.Error)(nil))
	}

	const wantCFunc = "alpm_initialize"
	if e.CFunc != wantCFunc {
		t.Errorf(`alpm.NewHandle(%q, %q) error.CFunc = %q, want %q`, root, dbpath, e.CFunc, wantCFunc)
	}

	const wantErrno = alpm.ErrNotADir
	if e.Errno != wantErrno {
		t.Errorf(`alpm.NewHandle(%q, %q) error.Errno = %q, want %q`, root, dbpath, e.Errno.Error(), wantErrno.Error())
	}
}

func samefile(fp1, fp2 string) bool {
	fi1, err1 := os.Stat(fp1)
	if err1 != nil {
		return false
	}

	fi2, err2 := os.Stat(fp2)
	if err2 != nil {
		return false
	}

	return os.SameFile(fi1, fi2)
}
