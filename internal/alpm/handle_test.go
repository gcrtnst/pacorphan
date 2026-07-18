package alpm

import "testing"

func TestNilHandleClose(t *testing.T) {
	h := (*Handle)(nil)

	err := h.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestZeroHandleClose(t *testing.T) {
	h := &Handle{}

	err := h.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestNilHandleRoot(t *testing.T) {
	h := (*Handle)(nil)

	got := h.Root()
	const want = ""
	if got != want {
		t.Errorf("h.Root() = %q, want %q", got, want)
	}
}

func TestZeroHandleRoot(t *testing.T) {
	h := &Handle{}

	got := h.Root()
	const want = ""
	if got != want {
		t.Errorf("h.Root() = %q, want %q", got, want)
	}
}

func TestNilHandleDBPath(t *testing.T) {
	h := (*Handle)(nil)

	got := h.DBPath()
	const want = ""
	if got != want {
		t.Errorf("h.DBPath() = %q, want %q", got, want)
	}
}

func TestZeroHandleDBPath(t *testing.T) {
	h := &Handle{}

	got := h.DBPath()
	const want = ""
	if got != want {
		t.Errorf("h.DBPath() = %q, want %q", got, want)
	}
}

func TestNilHandleLocalDB(t *testing.T) {
	h := (*Handle)(nil)

	got := h.LocalDB()
	want := (*DB)(nil)
	if got != want {
		t.Errorf("h.LocalDB() = %#v, want %#v", got, want)
	}
}

func TestZeroHandleLocalDB(t *testing.T) {
	h := &Handle{}

	got := h.LocalDB()
	want := (*DB)(nil)
	if got != want {
		t.Errorf("h.LocalDB() = %#v, want %#v", got, want)
	}
}
