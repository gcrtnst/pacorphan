package alpm

import "testing"

func TestErrnoError(t *testing.T) {
	got := ErrMemory.Error()
	want := "out of memory!"
	if got != want {
		t.Errorf("expected %#v, got %#v", want, got)
	}
}
