package main

import (
	"fmt"
	"slices"

	"github.com/gcrtnst/pacorphan/internal/alpm"
)

func init() { register("TestPkg", TestPkg) }
func TestPkg(t *T) {
	env := HelpEnv(t)

	srcA := NewPkgBuild("a", "0.0.1")
	HelpMakeAndInstall(t, env, srcA, false)

	srcB := NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a")
	HelpMakeAndInstall(t, env, srcB, true)

	srcC := NewPkgBuild("c", "0.0.3")
	srcC.OptDepends = append(srcC.OptDepends, "a")
	HelpMakeAndInstall(t, env, srcC, true)

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

	pkgsLen := pkgs.Len()
	if pkgsLen != 3 {
		t.Fatalf("pkgs.Len() = %d, want 3", pkgsLen)
	}

	elem := pkgs.Front()
	pkgA := elem.Data()

	pkgAName := pkgA.Name()
	if pkgAName != srcA.Name {
		t.Errorf("pkgA.Name() = %q, want %q", pkgAName, srcA.Name)
	}

	pkgAVer := pkgA.Version()
	srcAVer := fmt.Sprintf("%s-%d", srcA.Ver, srcA.Rel)
	if pkgAVer != srcAVer {
		t.Errorf("pkgA.Version() = %q, want %q", pkgAVer, srcAVer)
	}

	pkgADeps := pkgA.Depends()
	pkgADepLen := pkgADeps.Len()
	if pkgADepLen == len(srcA.Depends) {
		for i, pkgADep := range slices.Collect(pkgADeps.All()) {
			pkgADepStr := pkgADep.String()
			if pkgADepStr != srcA.Depends[i] {
				t.Errorf("pkgA.Depends()[%d].String() = %q, want %q", i, pkgADepStr, srcA.Depends[i])
			}
		}
	} else {
		t.Errorf("pkgA.Depends().Len() = %d, want %d", pkgADepLen, len(srcA.Depends))
	}

	pkgAOpts := pkgA.OptDepends()
	pkgAOptLen := pkgAOpts.Len()
	if pkgAOptLen == len(srcA.OptDepends) {
		for i, pkgAOpt := range slices.Collect(pkgAOpts.All()) {
			pkgAOptStr := pkgAOpt.String()
			if pkgAOptStr != srcA.OptDepends[i] {
				t.Errorf("pkgA.OptDepends()[%d].String() = %q, want %q", i, pkgAOptStr, srcA.OptDepends[i])
			}
		}
	} else {
		t.Errorf("pkgA.OptDepends().Len() = %d, want %d", pkgAOptLen, len(srcA.OptDepends))
	}

	pkgAReason := pkgA.Reason()
	if pkgAReason != alpm.PkgReasonDepend {
		t.Errorf("pkgA.Reason() = %d, want %d", pkgAReason, alpm.PkgReasonDepend)
	}

	elem = elem.Next()
	pkgB := elem.Data()

	pkgBName := pkgB.Name()
	if pkgBName != srcB.Name {
		t.Errorf("pkgB.Name() = %q, want %q", pkgBName, srcB.Name)
	}

	pkgBVer := pkgB.Version()
	srcBVer := fmt.Sprintf("%s-%d", srcB.Ver, srcB.Rel)
	if pkgBVer != srcBVer {
		t.Errorf("pkgB.Version() = %q, want %q", pkgBVer, srcBVer)
	}

	pkgBDeps := pkgB.Depends()
	pkgBDepLen := pkgBDeps.Len()
	if pkgBDepLen == len(srcB.Depends) {
		for i, pkgBDep := range slices.Collect(pkgBDeps.All()) {
			pkgBDepStr := pkgBDep.String()
			if pkgBDepStr != srcB.Depends[i] {
				t.Errorf("pkgB.Depends()[%d].String() = %q, want %q", i, pkgBDepStr, srcB.Depends[i])
			}
		}
	} else {
		t.Errorf("pkgB.Depends().Len() = %d, want %d", pkgBDepLen, len(srcB.Depends))
	}

	pkgBOpts := pkgB.OptDepends()
	pkgBOptLen := pkgBOpts.Len()
	if pkgBOptLen == len(srcB.OptDepends) {
		for i, pkgBOpt := range slices.Collect(pkgBOpts.All()) {
			pkgBOptStr := pkgBOpt.String()
			if pkgBOptStr != srcB.OptDepends[i] {
				t.Errorf("pkgB.OptDepends()[%d].String() = %q, want %q", i, pkgBOptStr, srcB.OptDepends[i])
			}
		}
	} else {
		t.Errorf("pkgB.OptDepends().Len() = %d, want %d", pkgBOptLen, len(srcB.OptDepends))
	}

	pkgBReason := pkgB.Reason()
	if pkgBReason != alpm.PkgReasonExplicit {
		t.Errorf("pkgB.Reason() = %d, want %d", pkgBReason, alpm.PkgReasonExplicit)
	}

	elem = elem.Next()
	pkgC := elem.Data()

	pkgCName := pkgC.Name()
	if pkgCName != srcC.Name {
		t.Errorf("pkgC.Name() = %q, want %q", pkgCName, srcC.Name)
	}

	pkgCVer := pkgC.Version()
	srcCVer := fmt.Sprintf("%s-%d", srcC.Ver, srcC.Rel)
	if pkgCVer != srcCVer {
		t.Errorf("pkgC.Version() = %q, want %q", pkgCVer, srcCVer)
	}

	pkgCDeps := pkgC.Depends()
	pkgCDepLen := pkgCDeps.Len()
	if pkgCDepLen == len(srcC.Depends) {
		for i, pkgCDep := range slices.Collect(pkgCDeps.All()) {
			pkgCDepStr := pkgCDep.String()
			if pkgCDepStr != srcC.Depends[i] {
				t.Errorf("pkgC.Depends()[%d].String() = %q, want %q", i, pkgCDepStr, srcC.Depends[i])
			}
		}
	} else {
		t.Errorf("pkgC.Depends().Len() = %d, want %d", pkgCDepLen, len(srcC.Depends))
	}

	pkgCOpts := pkgC.OptDepends()
	pkgCOptLen := pkgCOpts.Len()
	if pkgCOptLen == len(srcC.OptDepends) {
		for i, pkgCOpt := range slices.Collect(pkgCOpts.All()) {
			pkgCOptStr := pkgCOpt.String()
			if pkgCOptStr != srcC.OptDepends[i] {
				t.Errorf("pkgC.OptDepends()[%d].String() = %q, want %q", i, pkgCOptStr, srcC.OptDepends[i])
			}
		}
	} else {
		t.Errorf("pkgC.OptDepends().Len() = %d, want %d", pkgCOptLen, len(srcC.OptDepends))
	}

	pkgCReason := pkgC.Reason()
	if pkgCReason != alpm.PkgReasonExplicit {
		t.Errorf("pkgC.Reason() = %d, want %d", pkgCReason, alpm.PkgReasonExplicit)
	}
}

func init() { register("TestPkgClosed", TestPkgClosed) }
func TestPkgClosed(t *T) {
	env := HelpEnv(t)

	srcA := NewPkgBuild("a", "0.0.1")
	HelpMakeAndInstall(t, env, srcA, false)

	srcB := NewPkgBuild("b", "0.0.2")
	srcB.Depends = append(srcB.Depends, "a")
	HelpMakeAndInstall(t, env, srcB, true)

	srcC := NewPkgBuild("c", "0.0.3")
	srcC.OptDepends = append(srcC.OptDepends, "a")
	HelpMakeAndInstall(t, env, srcC, true)

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

	pkgsLen := pkgs.Len()
	if pkgsLen != 3 {
		t.Fatalf("pkgs.Len() = %d, want 3", pkgsLen)
	}

	elem := pkgs.Front()
	pkgA := elem.Data()
	pkgB := elem.Next().Data()
	pkgC := elem.Next().Next().Data()

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	pkgAName := pkgA.Name()
	if pkgAName != "" {
		t.Errorf(`pkgA.Name() = %q, want ""`, pkgAName)
	}

	pkgAVer := pkgA.Version()
	if pkgAVer != "" {
		t.Errorf(`pkgA.Version() = %q, want ""`, pkgAVer)
	}

	pkgADepLen := pkgA.Depends().Len()
	if pkgADepLen != 0 {
		t.Errorf("pkgA.Depends().Len() = %d, want 0", pkgADepLen)
	}

	pkgAOptLen := pkgA.OptDepends().Len()
	if pkgAOptLen != 0 {
		t.Errorf("pkgA.OptDepends().Len() = %d, want 0", pkgAOptLen)
	}

	pkgAReason := pkgA.Reason()
	if pkgAReason != alpm.PkgReasonUnknown {
		t.Errorf("pkgA.Reason() = %d, want %d", pkgAReason, alpm.PkgReasonUnknown)
	}

	pkgBName := pkgB.Name()
	if pkgBName != "" {
		t.Errorf(`pkgB.Name() = %q, want ""`, pkgBName)
	}

	pkgBVer := pkgB.Version()
	if pkgBVer != "" {
		t.Errorf(`pkgB.Version() = %q, want ""`, pkgBVer)
	}

	pkgBDepLen := pkgB.Depends().Len()
	if pkgBDepLen != 0 {
		t.Errorf("pkgB.Depends().Len() = %d, want 0", pkgBDepLen)
	}

	pkgBOptLen := pkgB.OptDepends().Len()
	if pkgBOptLen != 0 {
		t.Errorf("pkgB.OptDepends().Len() = %d, want 0", pkgBOptLen)
	}

	pkgBReason := pkgB.Reason()
	if pkgBReason != alpm.PkgReasonUnknown {
		t.Errorf("pkgB.Reason() = %d, want %d", pkgBReason, alpm.PkgReasonUnknown)
	}

	pkgCName := pkgC.Name()
	if pkgCName != "" {
		t.Errorf(`pkgC.Name() = %q, want ""`, pkgCName)
	}

	pkgCVer := pkgC.Version()
	if pkgCVer != "" {
		t.Errorf(`pkgC.Version() = %q, want ""`, pkgCVer)
	}

	pkgCDepLen := pkgC.Depends().Len()
	if pkgCDepLen != 0 {
		t.Errorf("pkgC.Depends().Len() = %d, want 0", pkgCDepLen)
	}

	pkgCOptLen := pkgC.OptDepends().Len()
	if pkgCOptLen != 0 {
		t.Errorf("pkgC.OptDepends().Len() = %d, want 0", pkgCOptLen)
	}

	pkgCReason := pkgC.Reason()
	if pkgCReason != alpm.PkgReasonUnknown {
		t.Errorf("pkgC.Reason() = %d, want %d", pkgCReason, alpm.PkgReasonUnknown)
	}
}

func init() { register("TestPkgListClosed", TestPkgListClosed) }
func TestPkgListClosed(t *T) {
	env := HelpEnv(t)
	HelpMakeAndInstall(t, env, NewPkgBuild("a", "0.0.1"), false)

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

	pkgsLen := pkgs.Len()
	if pkgsLen != 1 {
		t.Fatalf("before close: pkgs.Len() = %d, want 1", pkgsLen)
	}

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	pkgsLen = pkgs.Len()
	if pkgsLen != 0 {
		t.Fatalf("before close: pkgs.Len() = %d, want 0", pkgsLen)
	}
}
