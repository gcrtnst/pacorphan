package alpm

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestHandle(t *testing.T) {
	inRoot := t.TempDir()
	inDBPath := filepath.Join(inRoot, "var/lib/pacman/")

	errMkdir := os.MkdirAll(inDBPath, 0755|fs.ModeDir)
	if errMkdir != nil {
		t.Fatal(errMkdir)
	}

	h, errNew := NewHandle(inRoot, inDBPath)
	if errNew != nil {
		t.Fatal(errNew)
	}

	defer func() {
		errClose := h.Close()
		if errClose != nil {
			t.Error(errClose)
		}
	}()

	assertSameFile := func(name string, want, got string) {
		wantInfo, errWant := os.Stat(want)
		if errWant != nil {
			t.Error(errWant)
		}

		gotInfo, errGot := os.Stat(got)
		if errGot != nil {
			t.Error(errGot)
		}

		if !os.SameFile(wantInfo, gotInfo) {
			t.Errorf("%s: expected %#v, got %#v", name, want, got)
		}
	}

	gotRoot := h.Root()
	assertSameFile("root", inRoot, gotRoot)

	gotDBPath := h.DBPath()
	assertSameFile("dbpath", inDBPath, gotDBPath)
}

func TestNewHandleError(t *testing.T) {
	root := t.TempDir()
	dbpath := filepath.Join(root, "var/lib/pacman/")

	h, err := NewHandle(root, dbpath)
	if h != nil {
		t.Errorf("h: expected nil, got %#v", h)
	}
	if err != ErrNotADir {
		t.Errorf("err: expected %#v, got %#v", ErrNotADir, err)
	}
}
