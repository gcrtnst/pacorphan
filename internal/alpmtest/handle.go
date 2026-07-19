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

func init() { Register("TestHandleClosed", TestHandleClosed) }
func TestHandleClosed(t *T) {
	env := HelpEnv(t)

	h, errNew := alpm.NewHandle(env.Root, env.DBPath)
	if errNew != nil {
		t.Fatal(errNew)
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	errCloseDouble := h.Close()
	if errCloseDouble != nil {
		t.Error(errCloseDouble)
	}

	gotRoot := h.Root()
	const wantRoot = ""
	if gotRoot != wantRoot {
		t.Errorf("h.Root() = %q, want %q", gotRoot, wantRoot)
	}

	gotDBPath := h.DBPath()
	const wantDBPath = ""
	if gotDBPath != wantDBPath {
		t.Errorf("h.DBPath() = %q, want %q", gotDBPath, wantDBPath)
	}

	gotDB := h.LocalDB()
	wantDB := (*alpm.DB)(nil)
	if gotDB != wantDB {
		t.Errorf("h.LocalDB() = %#v, got %#v", gotDB, wantDB)
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

func init() { Register("TestFindDBsSatisfier", TestFindDBsSatisfier) }
func TestFindDBsSatisfier(t *T) {
	env := HelpEnv(t)
	HelpMakeAndInstall(t, env, NewPkgBuild("a", "0.0.1"), true)
	HelpMakeAndInstall(t, env, NewPkgBuild("b", "0.0.2"), true)

	h, errHandle := alpm.NewHandle(env.Root, env.DBPath)
	if errHandle != nil {
		t.Fatal(errHandle)
	}
	defer func() {
		err := h.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	l := alpm.NewDBList(h)
	l.PushBack(h.LocalDB())

	pkgA := h.FindDBsSatisfier(l, "a=0.0.1")
	pkgAName := pkgA.Name()
	if pkgAName != "a" {
		t.Errorf(`h.FindDBsSatisfier(l, "a=0.0.1").Name() = %q, want "a"`, pkgAName)
	}

	pkgB := h.FindDBsSatisfier(l, "b=0.0.2")
	pkgBName := pkgB.Name()
	if pkgBName != "b" {
		t.Errorf(`h.FindDBsSatisfier(l, "b=0.0.2").Name() = %q, want "b"`, pkgBName)
	}
}

func init() { Register("TestFindDBsSatisfierHandleClosed", TestFindDBsSatisfierHandleClosed) }
func TestFindDBsSatisfierHandleClosed(t *T) {
	env := HelpEnv(t)
	HelpMakeAndInstall(t, env, NewPkgBuild("a", "0.0.1"), true)

	h, errHandle := alpm.NewHandle(env.Root, env.DBPath)
	if errHandle != nil {
		t.Fatal(errHandle)
	}
	defer func() {
		err := h.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	l := alpm.NewDBList(h)
	l.PushBack(h.LocalDB())

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	pkg := h.FindDBsSatisfier(l, "a")
	if pkg != nil {
		t.Errorf(`h.FindDBsSatisfier(l, "a") = %#v, want nil`, pkg)
	}
}

func init() { Register("TestFindDBsSatisfierHandleNil", TestFindDBsSatisfierHandleNil) }
func TestFindDBsSatisfierHandleNil(t *T) {
	env := HelpEnv(t)
	HelpMakeAndInstall(t, env, NewPkgBuild("a", "0.0.1"), true)

	h, errHandle := alpm.NewHandle(env.Root, env.DBPath)
	if errHandle != nil {
		t.Fatal(errHandle)
	}
	defer func(h *alpm.Handle) {
		err := h.Close()
		if err != nil {
			t.Error(err)
		}
	}(h)

	l := alpm.NewDBList(h)
	l.PushBack(h.LocalDB())

	h = nil

	pkg := h.FindDBsSatisfier(l, "a")
	if pkg != nil {
		t.Errorf(`h.FindDBsSatisfier(l, "a") = %#v, want nil`, pkg)
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
