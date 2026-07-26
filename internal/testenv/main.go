package testenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
)

type testEntry struct {
	Name string
	Func TestFunc
}

type TestMain struct {
	Output io.Writer

	testList []testEntry
}

func NewTestMain() *TestMain {
	return &TestMain{Output: os.Stdout}
}

func (m *TestMain) Register(name string, fn TestFunc) {
	m.testList = append(m.testList, testEntry{
		Name: name,
		Func: fn,
	})
	slices.SortStableFunc(m.testList, func(a, b testEntry) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
}

func (m *TestMain) Run() int {
	failed := false
	for _, c := range m.testList {
		fmt.Fprintf(m.Output, "=== RUN   %s\n", c.Name)
		result := m.runTest(c.Func)
		switch result {
		case testSkip:
			fmt.Fprintf(m.Output, "--- SKIP: %s\n", c.Name)
		case testPass:
			fmt.Fprintf(m.Output, "--- PASS: %s\n", c.Name)
		default:
			failed = true
			fmt.Fprintf(m.Output, "--- FAIL: %s\n", c.Name)
		}
	}

	if failed {
		fmt.Fprintln(m.Output, "FAIL")
		return 1
	}
	fmt.Fprintln(m.Output, "PASS")
	return 0
}

func (m *TestMain) runTest(fn TestFunc) (result testResult) {
	result = testFail

	t := newT()
	t.output = newPrefixWriter(m.Output, "    ")

	defer func() {
		if t.Failed() {
			result = testFail
		} else if t.Skipped() {
			result = testSkip
		} else {
			result = testPass
		}
	}()

	defer func() {
		r := recover()
		if r != nil && r != errTestExit {
			panic(r)
		}
	}()

	defer func() {
		for _, f := range t.cleanups {
			defer f()
		}
	}()

	fn(t)
	return
}

type testResult int

const (
	testPass testResult = iota
	testFail
	testSkip
)

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
