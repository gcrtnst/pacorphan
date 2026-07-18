package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	failed := false
	for _, c := range testList {
		t := newT()

		fmt.Printf("=== RUN   %s\n", c.Name)
		runTest(t, c.Func)

		if t.Failed() {
			failed = true
			fmt.Printf("--- FAIL: %s\n", c.Name)
		} else {
			fmt.Printf("--- PASS: %s\n", c.Name)
		}
	}

	if failed {
		fmt.Println("FAIL")
		return 1
	}
	fmt.Println("PASS")
	return 0
}
