package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
)

var testList = []testEntry{}

type testEntry struct {
	Name string
	Func TestFunc
}

func Register(name string, fn TestFunc) {
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

type TestFunc func(*T)

var errTestExit = errors.New("alpmtest: test exit")

func runTest(t *T, fn TestFunc) {
	defer func() {
		r := recover()
		if r != nil && r != errTestExit {
			panic(r)
		}
	}()

	fn(t)
}

type T struct {
	failed bool
	output io.Writer
}

func newT() *T {
	return &T{
		failed: false,
		output: newPrefixWriter(os.Stdout, "    "),
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

type prefixWriter struct {
	w      io.Writer
	prefix string
	mid    bool
}

func newPrefixWriter(w io.Writer, prefix string) *prefixWriter {
	return &prefixWriter{
		w:      w,
		prefix: prefix,
		mid:    false,
	}
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		if !w.mid {
			_, err := io.WriteString(w.w, w.prefix)
			if err != nil {
				return n, err
			}

			w.mid = true
		}

		b := p[n:]
		i := bytes.IndexByte(b, '\n')

		if i < 0 {
			m, err := w.w.Write(b)
			n += m
			return n, err
		}

		m, err := w.w.Write(b[:i+1])
		n += m
		if err != nil {
			return n, err
		}
		w.mid = false
	}
	return n, nil
}
