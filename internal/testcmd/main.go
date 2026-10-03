package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var pacorphan = ""

var testMain = NewTestMain()

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

	if pacorphan == "" {
		var err error
		pacorphan, err = exec.LookPath("pacorphan")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", name, err)
			return 1
		}
	}

	return testMain.Run()
}
