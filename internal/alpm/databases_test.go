package alpm

import "testing"

func TestNilDBPkgCache(t *testing.T) {
	db := (*DB)(nil)

	l, err := db.PkgCache()
	if l != nil {
		t.Errorf("db.PkgCache() = %#v, want %#v", l, (*List[*Pkg])(nil))
	}
	if err != ErrHandleClosed {
		t.Errorf(`db.PkgCache() error = "%v", want "%v"`, err, ErrHandleClosed)
	}
}

func TestZeroDBPkgCache(t *testing.T) {
	db := &DB{}

	l, err := db.PkgCache()
	if l != nil {
		t.Errorf("db.PkgCache() = %#v, want %#v", l, (*List[*Pkg])(nil))
	}
	if err != ErrHandleClosed {
		t.Errorf(`db.PkgCache() error = "%v", want "%v"`, err, ErrHandleClosed)
	}
}
