package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/gcrtnst/pacorphan/internal/alpm"
	"github.com/spf13/pflag"
)

func main() {
	os.Exit(run())
}

type Pkg struct {
	Name    string
	Version string
}

func run() int {
	var fQuiet bool
	fs := pflag.NewFlagSet("pacorphan", pflag.ContinueOnError)
	fs.BoolVarP(&fQuiet, "quiet", "q", false, "show less information")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		if !errors.Is(errParse, pflag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %s\n", errParse)
		}
		return 1
	}

	orphans, err := FindOrphans()
	if err != nil {
		if errALPM, ok := errors.AsType[*alpm.Error](err); ok {
			switch errALPM.CFunc {
			case "alpm_initialize":
				fmt.Fprintf(os.Stderr, "error: failed to initialize alpm library: %s\n", errALPM.Errno.Message())
				return 1
			case "alpm_db_get_pkgcache":
				fmt.Fprintf(os.Stderr, "error: failed to load pkgcache: %s\n", errALPM.Errno.Message())
				return 1
			default:
				fmt.Fprintf(os.Stderr, "error: %s(): %s\n", errALPM.CFunc, errALPM.Errno.Message())
				return 1
			}
		} else if errors.Is(err, alpm.ErrHandleCloseFailed) {
			fmt.Fprintf(os.Stderr, "warn: failed to release alpm library\n")
			// continue
		} else {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return 1
		}
	}

	for _, pkg := range orphans {
		if fQuiet {
			fmt.Printf("%s\n", pkg.Name)
		} else {
			fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
		}
	}
	return 0
}

func FindOrphans() (orphans []Pkg, err error) {
	h, errHandleNew := alpm.NewHandle("/", "/var/lib/pacman")
	if errHandleNew != nil {
		return nil, err
	}
	defer func() {
		errHandleClose := h.Close()
		if errHandleClose != nil && err == nil {
			err = errHandleClose
		}
	}()

	db := h.LocalDB()
	pkgList, errPkgCache := db.PkgCache()
	if errPkgCache != nil {
		return nil, errPkgCache
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

	orphans = make([]Pkg, 0, len(mark))
	for _, pkg := range mark {
		orphans = append(orphans, Pkg{
			Name:    pkg.Name(),
			Version: pkg.Version(),
		})
	}
	slices.SortStableFunc(orphans, func(a, b Pkg) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return alpm.CompareVersion(a.Version, b.Version)
	})
	return orphans, nil
}

func pop[S ~[]E, E any](s S) (E, S) {
	var zero E

	e := s[len(s)-1]
	s[len(s)-1] = zero
	s = s[:len(s)-1]
	return e, s
}
