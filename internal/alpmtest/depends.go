package main

import "github.com/gcrtnst/pacorphan/internal/alpm"

func init() { Register("TestDepend", TestDepend) }
func TestDepend(t *T) {
	env := HelpEnv(t)

	srcA := NewPkgBuild("a", "0.0.1")
	HelpMakeAndInstall(t, env, srcA, false)

	srcB := NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a=0.0.1")
	HelpMakeAndInstall(t, env, srcB, true)

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

	depStr := dep.String()
	if depStr != srcB.Depends[0] {
		t.Errorf("dep.String() = %q, want %q", depStr, srcB.Depends[0])
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	depStrClosed := dep.String()
	if depStrClosed != "" {
		t.Errorf(`after close: dep.String() = %q, want ""`, depStrClosed)
	}
}

func init() { Register("TestDependListClosed", TestDependListClosed) }
func TestDependListClosed(t *T) {
	env := HelpEnv(t)

	srcA := NewPkgBuild("a", "0.0.1")
	HelpMakeAndInstall(t, env, srcA, false)

	srcB := NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a=0.0.1")
	HelpMakeAndInstall(t, env, srcB, true)

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
