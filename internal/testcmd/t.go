package testcmd

import (
	"errors"
	"fmt"
	"io"
)

var errTestExit = errors.New("testenv: test exit")

type TestFunc func(*T)

type T struct {
	failed   bool
	output   io.Writer
	cleanups []func()
}

func newT() *T {
	return &T{
		failed: false,
		output: io.Discard,
	}
}

func (t *T) Failed() bool {
	return t.failed
}

func (t *T) Output() io.Writer {
	return t.output
}

func (t *T) Fail() {
	t.failed = true
}

func (t *T) FailNow() {
	t.Fail()
	panic(errTestExit)
}

func (t *T) Error(args ...any) {
	t.Log(args...)
	t.Fail()
}

func (t *T) Errorf(format string, args ...any) {
	t.Logf(format, args...)
	t.Fail()
}

func (t *T) Fatal(args ...any) {
	t.Log(args...)
	t.FailNow()
}

func (t *T) Fatalf(format string, args ...any) {
	t.Logf(format, args...)
	t.FailNow()
}

func (t *T) Log(args ...any) {
	_, _ = fmt.Fprintln(t.Output(), args...)
}

func (t *T) Logf(format string, args ...any) {
	b := fmt.Appendf(nil, format, args...)
	if len(b) <= 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	_, _ = t.Output().Write(b)
}

func (t *T) Cleanup(f func()) {
	t.cleanups = append(t.cleanups, f)
}
