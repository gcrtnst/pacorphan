package exiterr

import (
	"errors"
	"os/exec"
	"regexp"
	"testing"
)

func TestWrapExitError(t *testing.T) {
	cmd := exec.Command("echo", "ok")
	raw := &exec.ExitError{Stderr: []byte("testing")}
	err := Wrap(cmd, raw)

	if e, ok := errors.AsType[*ExitError](err); ok {
		if e.Cmd != cmd {
			t.Errorf("Wrap(cmd, raw).Cmd = %#v, want %#v", e.Cmd, cmd)
		}
		if e.Err != raw {
			t.Errorf("Wrap(cmd, raw).Err = %#v, want %#v", e.Err, raw)
		}
	} else {
		t.Errorf("Wrap(cmd, raw) type = %T, want %T", err, (*ExitError)(nil))
	}
}

func TestWrapOtherError(t *testing.T) {
	cmd := exec.Command("echo", "ok")
	raw := errors.New("testing")
	err := Wrap(cmd, raw)

	if err != raw {
		t.Errorf("Wrap(cmd, raw) = %#v, want %#v", err, raw)
	}
}

func TestExitErrorError(t *testing.T) {
	tt := []struct {
		name   string
		in     *ExitError
		wantRe *regexp.Regexp
	}{
		{
			name: "Normal",
			in: &ExitError{
				Cmd: exec.Command("echo", "ok"),
				Err: &exec.ExitError{},
			},
			wantRe: regexp.MustCompile(`^exec \[echo ok\]: .+$`),
		},
		{
			name:   "Nil",
			in:     nil,
			wantRe: regexp.MustCompile(`^<nil>$`),
		},
		{
			name: "NilCmd",
			in: &ExitError{
				Cmd: nil,
				Err: &exec.ExitError{},
			},
			wantRe: regexp.MustCompile(`^exec: .+$`),
		},
		{
			name: "NilErr",
			in: &ExitError{
				Cmd: exec.Command("echo", "ok"),
				Err: nil,
			},
			wantRe: regexp.MustCompile(`^exec \[echo ok\]: <nil>$`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Error()
			if !tc.wantRe.MatchString(got) {
				t.Errorf("err.Error() = %q, want regex `%s`", got, tc.wantRe)
			}
		})
	}
}

func TestExitErrorUnwrap(t *testing.T) {
	sentinel := &exec.ExitError{}

	tt := []struct {
		name string
		in   *ExitError
		want error
	}{
		{
			name: "Normal",
			in:   &ExitError{Err: sentinel},
			want: sentinel,
		},
		{
			name: "Nil",
			in:   nil,
			want: nil,
		},
		{
			name: "NilErr",
			in:   &ExitError{},
			want: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Unwrap()
			if got != tc.want {
				t.Errorf("err.Unwrap() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
