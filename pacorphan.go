package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strconv"

	"github.com/gcrtnst/pacorphan/internal/alpm"
	"github.com/gcrtnst/pacorphan/internal/exiterr"
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
	fQuiet := false
	fHelp := false
	fOptDepends := true
	fDBPath := ""
	fRoot := ""
	fConfig := ""
	fSysRoot := ""
	fVersion := false
	fs := pflag.NewFlagSet("pacorphan", pflag.ContinueOnError)
	fs.SortFlags = false
	fs.BoolVarP(&fQuiet, "quiet", "q", fQuiet, "show less information")
	fsRespectVarP(fs, &fOptDepends, "respect-optdepends", "o", "treat optdepends as required (default)")
	fsIgnoreVarP(fs, &fOptDepends, "ignore-optdepends", "O", "ignore optdepends")
	fs.StringVar(&fDBPath, "dbpath", fDBPath, "set an alternate database location")
	fs.StringVar(&fRoot, "root", fRoot, "set an alternate installation root")
	fs.StringVar(&fConfig, "config", fConfig, "set an alternate configuration file")
	fs.StringVar(&fSysRoot, "sysroot", fSysRoot, "set an alternate system root")
	fs.BoolVarP(&fHelp, "help", "h", fHelp, "show this help message and exit")
	fs.BoolVar(&fVersion, "version", fVersion, "show version information and exit")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", errParse)
		return 2
	}

	if fHelp {
		fmt.Printf("Usage of %s:\n", fs.Name())
		fmt.Print(fs.FlagUsages())
		return 0
	}

	if fVersion {
		fmt.Printf("%s %s - libalpm v%s\n", fs.Name(), Version(), alpm.Version())
		return 0
	}

	conf := &PacmanConf{
		Config:  fConfig,
		Root:    fRoot,
		SysRoot: fSysRoot,
		Stderr:  os.Stderr,
	}

	sysroot := fSysRoot
	if sysroot == "" {
		sysroot = "/"
	}

	var root string
	if fRoot != "" {
		root = filepath.Join(sysroot, fRoot)
	} else {
		var err error
		root, err = conf.Get("RootDir")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return 1
		}
	}

	var dbpath string
	if fDBPath != "" {
		dbpath = filepath.Join(sysroot, fDBPath)
	} else {
		var err error
		dbpath, err = conf.Get("DBPath")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return 1
		}
	}

	mask := Depends
	if fOptDepends {
		mask |= OptDepends
	}

	orphans, miss, err := FindOrphans(root, dbpath, mask)
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
			fmt.Fprintf(os.Stderr, "warning: failed to release alpm library\n")
			// continue
		} else {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return 1
		}
	}

	for _, mdep := range miss {
		if mdep.InstalledVersion == "" {
			fmt.Fprintf(os.Stderr, "warning: '%s' requires '%s', which is not installed\n", mdep.DependentPkgName, mdep.DepString)
		} else if mdep.DepType&Depends != 0 {
			fmt.Fprintf(os.Stderr, "warning: '%s' requires '%s', but version %s is installed\n", mdep.DependentPkgName, mdep.DepString, mdep.InstalledVersion)
		} else {
			fmt.Fprintf(os.Stderr, "warning: '%s' recommends '%s', but version %s is installed\n", mdep.DependentPkgName, mdep.DepString, mdep.InstalledVersion)
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

func fsRespectVarP(fs *pflag.FlagSet, p *bool, name string, shorthand string, usage string) {
	fs.BoolFuncP(name, shorthand, usage, func(s string) error {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		*p = v
		return nil
	})
}

func fsIgnoreVarP(fs *pflag.FlagSet, p *bool, name string, shorthand string, usage string) {
	fs.BoolFuncP(name, shorthand, usage, func(s string) error {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		*p = !v
		return nil
	})
}

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}

	main := &info.Main
	for main.Replace != nil {
		main = main.Replace
	}

	return main.Version
}

func FindOrphans(root string, dbpath string, mask DependType) (_ []Pkg, _ []MissingDep, err error) {
	h, errHandleNew := alpm.NewHandle(root, dbpath)
	if errHandleNew != nil {
		return nil, nil, errHandleNew
	}
	defer func() {
		errHandleClose := h.Close()
		if errHandleClose != nil && err == nil {
			err = errHandleClose
		}
	}()

	db := h.LocalDB()
	dbList := alpm.NewDBList(h)
	defer dbList.Free()
	dbList.PushBack(db)

	pkgList, errPkgCache := db.PkgCache()
	if errPkgCache != nil {
		return nil, nil, errPkgCache
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

	var miss []MissingDep
	for len(stack) > 0 {
		var pkg *alpm.Pkg
		pkg, stack = pop(stack)

		for meta := range IterDepend(pkg, mask) {
			dep := meta.Depend
			depPkg := h.FindDBsSatisfier(dbList, dep.String())
			if depPkg != nil {
				depName := depPkg.Name()
				if _, ok := mark[depName]; ok {
					delete(mark, depName)
					stack = append(stack, depPkg)
				}
			} else {
				mdep := MissingDep{
					DependentPkgName: pkg.Name(),
					DepString:        dep.String(),
					DepType:          meta.Type,
				}

				depPkg = h.FindDBsSatisfier(dbList, dep.Name())
				if depPkg != nil {
					mdep.InstalledVersion = depPkg.Version()
				}

				if mdep.DepType&Depends != 0 || mdep.InstalledVersion != "" {
					miss = append(miss, mdep)
				}
			}
		}
	}

	orphans := make([]Pkg, 0, len(mark))
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

	slices.SortStableFunc(miss, func(a, b MissingDep) int {
		if a.DependentPkgName < b.DependentPkgName {
			return -1
		}
		if a.DependentPkgName > b.DependentPkgName {
			return 1
		}
		if a.DepType < b.DepType {
			return -1
		}
		if a.DepType > b.DepType {
			return 1
		}
		if a.DepString < b.DepString {
			return -1
		}
		if a.DepString > b.DepString {
			return 1
		}
		if a.InstalledVersion < b.InstalledVersion {
			return -1
		}
		if a.InstalledVersion > b.InstalledVersion {
			return 1
		}
		return 0
	})

	return orphans, miss, nil
}

type PacmanConf struct {
	Config  string
	Root    string
	SysRoot string

	Stderr io.Writer
}

func (c *PacmanConf) Get(directive string) (string, error) {
	args := make([]string, 0, 4)
	if v := c.Config; v != "" {
		args = append(args, "--config="+v)
	}
	if v := c.Root; v != "" {
		args = append(args, "--rootdir="+v)
	}
	if v := c.SysRoot; v != "" {
		args = append(args, "--sysroot="+v)
	}
	args = append(args, directive)

	buf := new(bytes.Buffer)
	cmd := exec.Command("pacman-conf", args...)
	cmd.Stdout = buf
	cmd.Stderr = c.Stderr

	err := cmd.Run()
	if err != nil {
		return "", exiterr.Wrap(cmd, err)
	}

	out := buf.Bytes()
	if len(out) >= 1 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

type Depend struct {
	Type   DependType
	Depend *alpm.Depend
}

func IterDepend(pkg *alpm.Pkg, mask DependType) iter.Seq[*Depend] {
	return func(yield func(*Depend) bool) {
		if mask&Depends != 0 {
			for dep := range pkg.Depends().All() {
				e := &Depend{Type: Depends, Depend: dep}
				if !yield(e) {
					return
				}
			}
		}
		if mask&OptDepends != 0 {
			for dep := range pkg.OptDepends().All() {
				e := &Depend{Type: OptDepends, Depend: dep}
				if !yield(e) {
					return
				}
			}
		}
	}
}

type MissingDep struct {
	DependentPkgName string
	DepString        string
	DepType          DependType
	InstalledVersion string
}

type DependType int

const (
	Depends DependType = 1 << iota
	OptDepends
)

func pop[S ~[]E, E any](s S) (E, S) {
	var zero E

	e := s[len(s)-1]
	s[len(s)-1] = zero
	s = s[:len(s)-1]
	return e, s
}
