package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/gcrtnst/pacorphan/internal/testenv"
)

var pacorphan = ""

var testMain = testenv.NewTestMain()

func main() {
	os.Exit(run())
}

func run() int {
	const name = "testcmd"
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.StringVar(&pacorphan, "pacorphan", "pacorphan", "path to pacorphan binary")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		if errors.Is(errParse, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, errParse)
		return 2
	}

	var errWhichPacorphan error
	pacorphan, errWhichPacorphan = exec.LookPath(pacorphan)
	if errWhichPacorphan != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, errWhichPacorphan)
		return 1
	}

	return testMain.Run()
}
