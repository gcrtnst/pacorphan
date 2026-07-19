package main

import "github.com/gcrtnst/pacorphan/internal/alpm"

func init() { Register("TestDBClosed", TestDBClosed) }
func TestDBClosed(t *T) {
	env := HelpEnv(t)

	h, errHandle := alpm.NewHandle(env.Root, env.DBPath)
	if errHandle != nil {
		t.Fatal(errHandle)
	}
	defer func() {
		err := h.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	db := h.LocalDB()

	errClose := h.Close()
	if errClose != nil {
		t.Fatal(errClose)
	}

	pkgs, errPkgs := db.PkgCache()
	if pkgs != nil {
		t.Errorf("db.PkgCache() = %#v, want %#v", pkgs, (*alpm.List[*alpm.Pkg])(nil))
	}
	if errPkgs != alpm.ErrHandleClosed {
		t.Errorf(`db.PkgCache() error = "%v", want "%v"`, errPkgs, alpm.ErrHandleClosed)
	}
}
