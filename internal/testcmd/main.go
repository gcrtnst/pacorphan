package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var testMain = NewTestMain()

func register(name string, fn TestFunc) {
	testMain.Register(name, fn)
}

func main() {
	os.Exit(run())
}

func run() int {
	const name = "testcmd"
	var pacorphan string
	var pacman string
	var makepkg string
	var unshare string
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.StringVar(&pacorphan, "pacorphan", "", "path to pacorphan binary")
	fs.StringVar(&pacman, "pacman", "", "path to pacman binary")
	fs.StringVar(&makepkg, "makepkg", "", "path to makepkg binary")
	fs.StringVar(&unshare, "unshare", "", "path to unshare binary")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		if errors.Is(errParse, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, errParse)
		return 2
	}

	if pacorphan == "" {
		var err error
		pacorphan, err = exec.LookPath("pacorphan")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, err)
			return 1
		}
	}

	if pacman == "" {
		var err error
		pacman, err = exec.LookPath("pacman")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, err)
			return 1
		}
	}

	if makepkg == "" {
		var err error
		makepkg, err = exec.LookPath("makepkg")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, err)
			return 1
		}
	}

	if unshare == "" {
		var err error
		unshare, err = exec.LookPath("unshare")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, err)
			return 1
		}
	}

	testMain.Config.pacorphan = pacorphan
	testMain.Config.pacman = pacman
	testMain.Config.makepkg = makepkg
	testMain.Config.unshare = unshare
	return testMain.Run()
}
