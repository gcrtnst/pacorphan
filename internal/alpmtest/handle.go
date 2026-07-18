package main

import (
	"os"
	"path/filepath"

	"github.com/gcrtnst/pacorphan/internal/alpm"
)

func init() { Register("TestHandleErrorInit", TestHandleErrorInit) }
func TestHandleErrorInit(t *T) {
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
