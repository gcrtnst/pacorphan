package main

import (
	"os"

	"github.com/gcrtnst/pacorphan/internal/alpm"
)

func init() { Register("TestDBClosed", TestDBClosed) }
func TestDBClosed(t *T) {
	env := HelpEnv(t)

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

	db := h.LocalDB()

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	pkgs, errPkgs := db.PkgCache()
	if pkgs != nil {
		t.Errorf("db.PkgCache() = %#v, want %#v", pkgs, (*alpm.List[*alpm.Pkg])(nil))
	}
	if errPkgs != alpm.ErrHandleClosed {
		t.Errorf(`db.PkgCache() error = "%v", want "%v"`, errPkgs, alpm.ErrHandleClosed)
	}
}

func init() { Register("TestDBInvalid", TestDBInvalid) }
func TestDBInvalid(t *T) {
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

	db := h.LocalDB()

	errRemove := os.RemoveAll(env.DBPath)
	if errRemove != nil {
		t.Fatal(errRemove)
	}

	pkgs, errPkgs := db.PkgCache()
	if pkgs != nil {
		t.Errorf("db.PkgCache() = %#v, want %#v", pkgs, (*alpm.List[*alpm.Pkg])(nil))
	}
	if e, ok := errPkgs.(*alpm.Error); ok {
		const wantCFunc = "alpm_db_get_pkgcache"
		if e.CFunc != wantCFunc {
			t.Errorf("db.PkgCache() error.CFunc = %q, want %q", e.CFunc, wantCFunc)
		}

		if e.Errno != alpm.ErrDBOpen {
			t.Errorf(`db.PkgCache() error.Errno = "%v", want "%v"`, e.Errno, alpm.ErrDBOpen)
		}
	} else {
		t.Errorf("db.PkgCache() error type = %T, want %T", errPkgs, (*alpm.Error)(nil))
	}
}
