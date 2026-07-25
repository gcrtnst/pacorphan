package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/gcrtnst/pacorphan/internal/testenv"
	"github.com/spf13/pflag"
)

var pacorphan = ""
var testMain = testenv.NewTestMain()

func main() {
	os.Exit(run())
}

func run() int {
	var cmd string
	fs := pflag.NewFlagSet("pacorphantest", pflag.ContinueOnError)
	fs.StringVarP(&cmd, "cmd", "c", "", "path to pacorphan binary")

	errParse := fs.Parse(os.Args[1:])
	if errParse != nil {
		if errors.Is(errParse, pflag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", errParse)
		return 2
	}

	if cmd == "" {
		var err error
		cmd, err = exec.LookPath("pacorphan")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return 1
		}
	}
	pacorphan = cmd

	return testMain.Run()
}
