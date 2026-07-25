package main

import (
	"os"

	"github.com/gcrtnst/pacorphan/internal/testcmd"
)

var testMain = testcmd.NewTestMain()

func main() {
	os.Exit(testMain.Run())
}
