package alpm

import "testing"

func TestNilDependName(t *testing.T) {
	d := (*Depend)(nil)

	got := d.Name()
	const want = ""
	if got != want {
		t.Errorf("d.Name() = %q, want %q", got, want)
	}
}

func TestZeroDependName(t *testing.T) {
	d := &Depend{}

	got := d.Name()
	const want = ""
	if got != want {
		t.Errorf("d.Name() = %q, want %q", got, want)
	}
}

func TestNilDependString(t *testing.T) {
	d := (*Depend)(nil)

	got := d.String()
	const want = ""
	if got != want {
		t.Errorf("d.String() = %q, want %q", got, want)
	}
}

func TestZeroDependString(t *testing.T) {
	d := &Depend{}

	got := d.String()
	const want = ""
	if got != want {
		t.Errorf("d.String() = %q, want %q", got, want)
	}
}
