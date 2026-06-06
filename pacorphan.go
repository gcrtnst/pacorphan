package main

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/gcrtnst/pacorphan/internal/alpm"
)

func main() {
	os.Exit(run())
}

func run() (code int) {
	h, errHandleNew := alpm.NewHandle("/", "/var/lib/pacman")
	if errHandleNew != nil {
		fmt.Fprintf(os.Stderr, "error: cannot initialize alpm: %s\n", errHandleNew.Error())
		return 1
	}
	defer func() {
		errHandleClose := h.Close()
		if errHandleClose != nil {
			fmt.Fprintf(os.Stderr, "error: cannot release alpm: %s\n", errHandleClose.Error())
			code = 1
		}
	}()

	db := h.LocalDB()
	pkgList, errPkgCache := db.PkgCache()
	if errPkgCache != nil {
		fmt.Fprintf(os.Stderr, "error: cannot load package cache: %s\n", errPkgCache.Error())
		return 1
	}

	mark := make(map[string]*alpm.Pkg)
	for pkg := range pkgList.All() {
		mark[pkg.Name()] = pkg
	}

	var stack []*alpm.Pkg
	for pkg := range pkgList.All() {
		if pkg.Reason() == alpm.PkgReasonExplicit {
			delete(mark, pkg.Name())
			stack = append(stack, pkg)
		}
	}

	for len(stack) > 0 {
		var pkg *alpm.Pkg
		pkg, stack = pop(stack)

		for dep := range pkg.Depends().All() {
			depPkg := alpm.FindSatisfier(pkgList, dep.String())
			if depPkg != nil {
				depName := depPkg.Name()
				if _, ok := mark[depName]; ok {
					delete(mark, depName)
					stack = append(stack, depPkg)
				}
			}
		}
	}

	orphans := slices.Collect(maps.Values(mark))
	slices.SortStableFunc(orphans, func(a, b *alpm.Pkg) int {
		aName := a.Name()
		bName := b.Name()
		if aName < bName {
			return -1
		}
		if aName > bName {
			return 1
		}

		aVersion := a.Version()
		bVersion := b.Version()
		if aVersion < bVersion {
			return -1
		}
		if aVersion > bVersion {
			return 1
		}

		return 0
	})

	for _, pkg := range orphans {
		fmt.Printf("%s %s\n", pkg.Name(), pkg.Version())
	}
	return 0
}

func pop[S ~[]E, E any](s S) (E, S) {
	var zero E

	e := s[len(s)-1]
	s[len(s)-1] = zero
	s = s[:len(s)-1]
	return e, s
}
