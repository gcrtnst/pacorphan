package main

type TestFunc func(*T)

type TestEntry struct {
	Name string
	Func TestFunc
}

var TestList = []TestEntry{}
