package alpm

import (
	"errors"
	"regexp"
	"testing"
)

func TestErrorError(t *testing.T) {
	err := &Error{CFunc: "alpm_strerror", Errno: ErrMemory}
	got := err.Error()
	pat := `^alpm: alpm_strerror\(\): (.|\n)+$`
	if !regexp.MustCompile(pat).MatchString(got) {
		t.Errorf("(%#v).Error() == %q, want format `%s`", err, got, pat)
	}
}

func TestErrorIs(t *testing.T) {
	err := &Error{CFunc: "alpm_strerror", Errno: ErrMemory}
	if !errors.Is(err, ErrMemory) {
		t.Errorf("errors.Is(%#v, %#v) = false, want true", err, ErrMemory)
	}
}

func TestErrorAs(t *testing.T) {
	err := &Error{CFunc: "alpm_strerror", Errno: ErrMemory}
	got, ok := errors.AsType[Errno](err)
	if !(ok && got == ErrMemory) {
		t.Errorf("errors.AsType[Errno](%#v) = (%#v, %t), want = (%#v, true)", err, got, ok, ErrMemory)
	}
}

func TestErrnoError(t *testing.T) {
	got := ErrMemory.Error()
	pat := `^alpm: (.|\n)+$`
	if !regexp.MustCompile(pat).MatchString(got) {
		t.Errorf("ErrMemory.Message() == %q, want format `%s`", got, pat)
	}
}

func TestErrnoMessage(t *testing.T) {
	got := ErrMemory.Message()
	if len(got) <= 0 {
		t.Errorf("ErrMemory.Message() == %q, want non-empty message", got)
	}
}
