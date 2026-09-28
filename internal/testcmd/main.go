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
	fs := flag.NewFlagSet("testcmd", flag.ContinueOnError)
	fs.StringVar(&pacorphan, "pacorphan", "", "path to pacorphan binary")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		if errors.Is(errParse, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", errParse)
		return 2
	}

	if pacorphan == "" {
		var err error
		pacorphan, err = exec.LookPath("pacorphan")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return 1
		}
	}

	return testMain.Run()
}
