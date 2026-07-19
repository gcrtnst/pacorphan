package main

import (
	"errors"
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

func init() { Register("TestDBList", TestDBList) }
func TestDBList(t *T) {
	env := HelpEnv(t)

	HelpMakeAndInstall(t, env, NewPkgBuild("a", "1.2.3"), true)
	const pkgName = "a"
	const pkgVer = "1.2.3-1"
	const pkgCnt = 1

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

	dbs := alpm.NewDBList(h)
	defer dbs.Free()

	db1 := h.LocalDB()
	pkgs1, errPkgs1 := db1.PkgCache()
	if errPkgs1 != nil {
		t.Fatal(errPkgs1)
	}

	pkgs1Len := pkgs1.Len()
	if pkgs1Len != pkgCnt {
		t.Fatalf("pkgs1.Len() = %d, want %d", pkgs1Len, pkgCnt)
	}

	pkg1 := pkgs1.Front().Data()
	pkg1Name := pkg1.Name()
	pkg1Ver := pkg1.Version()
	if pkg1Name != pkgName {
		t.Errorf("pkg1.Name() = %q, want %q", pkg1Name, pkgName)
	}
	if pkg1Ver != pkgVer {
		t.Errorf("pkg1.Version() = %q, want %q", pkg1Ver, pkgVer)
	}

	if t.Failed() {
		t.FailNow()
	}

	db2 := dbs.PushBack(db1).Data()
	pkgs2, errPkgs2 := db2.PkgCache()
	if errPkgs2 != nil {
		t.Fatal(errPkgs2)
	}

	pkgs2Len := pkgs2.Len()
	if pkgs2Len != pkgCnt {
		t.Fatalf("pkgs2.Len() = %d, want %d", pkgs2Len, pkgCnt)
	}

	pkg2 := pkgs2.Front().Data()
	pkg2Name := pkg2.Name()
	pkg2Ver := pkg2.Version()
	if pkg2Name != pkgName {
		t.Errorf("pkg2.Name() = %q, want %q", pkg2Name, pkgName)
	}
	if pkg2Ver != pkgVer {
		t.Errorf("pkg2.Version() = %q, want %q", pkg2Ver, pkgVer)
	}

	if t.Failed() {
		t.FailNow()
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	_, errPkgs1Closed := db1.PkgCache()
	if !errors.Is(errPkgs1Closed, alpm.ErrHandleClosed) {
		t.Errorf(`db1.PkgCache() error = "%v", want %v`, errPkgs1Closed, alpm.ErrHandleClosed)
	}

	_, errPkgs2Closed := db2.PkgCache()
	if !errors.Is(errPkgs2Closed, alpm.ErrHandleClosed) {
		t.Errorf(`db2.PkgCache() error = "%v", want %v`, errPkgs2Closed, alpm.ErrHandleClosed)
	}
}
