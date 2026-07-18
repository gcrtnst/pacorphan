package alpm

import "testing"

func TestNilHandleClose(t *testing.T) {
	h := (*Handle)(nil)

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

func TestNilHandleDBPath(t *testing.T) {
	h := (*Handle)(nil)

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
