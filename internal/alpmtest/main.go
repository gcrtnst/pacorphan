package main

import (
	"fmt"
	"os"
	"slices"
)

var testList = []testEntry{}

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

type testEntry struct {
	Name string
	Func TestFunc
}

func register(name string, fn TestFunc) {
	testList = append(testList, testEntry{
		Name: name,
		Func: fn,
	})
	slices.SortStableFunc(testList, func(a, b testEntry) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
}
