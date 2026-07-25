package main

import (
	"os"

	"github.com/gcrtnst/pacorphan/internal/testenv"
)

var testMain = testenv.NewTestMain()

func main() {
	os.Exit(testMain.Run())
}
