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
		fmt.Printf("=== RUN   %s\n", c.Name)
		f := runTest(c.Func)
		if f {
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
