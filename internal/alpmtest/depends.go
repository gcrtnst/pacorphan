package main

import (
	"github.com/gcrtnst/pacorphan/internal/alpm"
	"github.com/gcrtnst/pacorphan/internal/testenv"
)

func init() { testMain.Register("TestALPMDepend", TestALPMDepend) }
func TestALPMDepend(t *testenv.T) {
	env := testenv.HelpEnv(t)

	srcA := testenv.NewPkgBuild("a", "0.0.1")
	testenv.HelpMakeAndInstall(t, env, srcA, false)

	srcB := testenv.NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a=0.0.1")
	testenv.HelpMakeAndInstall(t, env, srcB, true)

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
	pkgs, errPkgs := db.PkgCache()
	if errPkgs != nil {
		t.Fatal(errPkgs)
	}

	pkgB := (*alpm.Pkg)(nil)
	for pkg := range pkgs.All() {
		if pkg.Name() == "b" {
			pkgB = pkg
			break
		}
	}
	if pkgB == nil {
		t.Fatal("package b not found")
	}

	dep := pkgB.Depends().Front().Data()

	depName := dep.Name()
	if depName != srcA.Name {
		t.Errorf("dep.Name() = %q, want %q", depName, srcA.Name)
	}

	depStr := dep.String()
	if depStr != srcB.Depends[0] {
		t.Errorf("dep.String() = %q, want %q", depStr, srcB.Depends[0])
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	depNameClosed := dep.Name()
	if depNameClosed != "" {
		t.Errorf(`after close: dep.Name() = %q, want ""`, depNameClosed)
	}

	depStrClosed := dep.String()
	if depStrClosed != "" {
		t.Errorf(`after close: dep.String() = %q, want ""`, depStrClosed)
	}
}

func init() { testMain.Register("TestALPMDependListClosed", TestALPMDependListClosed) }
func TestALPMDependListClosed(t *testenv.T) {
	env := testenv.HelpEnv(t)

	srcA := testenv.NewPkgBuild("a", "0.0.1")
	testenv.HelpMakeAndInstall(t, env, srcA, false)

	srcB := testenv.NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a=0.0.1")
	testenv.HelpMakeAndInstall(t, env, srcB, true)

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
	pkgs, errPkgs := db.PkgCache()
	if errPkgs != nil {
		t.Fatal(errPkgs)
	}

	pkgB := (*alpm.Pkg)(nil)
	for pkg := range pkgs.All() {
		if pkg.Name() == "b" {
			pkgB = pkg
			break
		}
	}
	if pkgB == nil {
		t.Fatal("package b not found")
	}

	deps := pkgB.Depends()

	depsLen := deps.Len()
	if depsLen != 1 {
		t.Errorf("before close: deps.Len() = %d, want 1", depsLen)
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	depsLenClosed := deps.Len()
	if depsLenClosed != 0 {
		t.Errorf("after close: deps.Len() = %d, want 0", depsLenClosed)
	}
}
