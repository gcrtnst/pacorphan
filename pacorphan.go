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
	"slices"

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
	fStrict := false
	fDBPath := ""
	fRoot := ""
	fConfig := ""
	fSysRoot := ""
	fs := pflag.NewFlagSet("pacorphan", pflag.ContinueOnError)
	fs.BoolVarP(&fQuiet, "quiet", "q", fQuiet, "show less information")
	fs.BoolVarP(&fStrict, "strict", "t", fStrict, "ignore optdepends")
	fs.StringVarP(&fDBPath, "dbpath", "b", fDBPath, "set an alternate database location")
	fs.StringVarP(&fRoot, "root", "R", fRoot, "set an alternate installation root")
	fs.StringVarP(&fConfig, "config", "C", fConfig, "set an alternate configuration file")
	fs.StringVarP(&fSysRoot, "sysroot", "S", fSysRoot, "set an alternate system root")
	fs.BoolVarP(&fHelp, "help", "h", fHelp, "show this help message and exit")

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

	conf := &PacmanConf{
		Config:  fConfig,
		Root:    fRoot,
		SysRoot: fSysRoot,
		Stderr:  os.Stderr,
	}

	sysroot := fSysRoot
	if sysroot != "" {
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

	orphans, err := FindOrphans(root, dbpath, fStrict)
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

	for _, pkg := range orphans {
		if fQuiet {
			fmt.Printf("%s\n", pkg.Name)
		} else {
			fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
		}
	}
	return 0
}

func FindOrphans(root string, dbpath string, strict bool) (orphans []Pkg, err error) {
	h, errHandleNew := alpm.NewHandle(root, dbpath)
	if errHandleNew != nil {
		return nil, errHandleNew
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

		depSeq := pkg.Depends().All()
		if !strict {
			depSeq = concat(depSeq, pkg.OptDepends().All())
		}

		for dep := range depSeq {
			depPkg := h.FindDBsSatisfier(dbList, dep.String())
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

func pop[S ~[]E, E any](s S) (E, S) {
	var zero E

	e := s[len(s)-1]
	s[len(s)-1] = zero
	s = s[:len(s)-1]
	return e, s
}

func concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}
