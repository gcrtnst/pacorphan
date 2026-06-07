package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type DB struct {
	h *Handle
	c *C.alpm_db_t
}

func newDB(h *Handle, c_db *C.alpm_db_t) *DB {
	defer runtime.KeepAlive(h)

	c_handle := C.alpm_db_get_handle(c_db)
	if c_handle != h.c {
		panic("alpm: database handle mismatch")
	}

	return &DB{h: h, c: c_db}
}

func (d *DB) alive() bool {
	if d == nil || d.c == nil {
		return false
	}
	if !d.h.alive() {
		d.c = nil
		return false
	}
	return true
}

func (d *DB) PkgCache() (*List[*Pkg], error) {
	defer runtime.KeepAlive(d)
	if !d.alive() {
		return nil, ErrHandleClosed
	}

	c_list := C.alpm_db_get_pkgcache(d.c)
	if c_list == nil {
		return nil, d.h.error("alpm_db_get_pkgcache")
	}

	l := newPkgList(d.h, d, c_list)
	return l, nil
}

func newDBList(h *Handle) *List[*DB] {
	c_list := (*C.alpm_list_t)(nil)
	l := &List[*DB]{
		c: &c_list,
		fi: func(db *DB) unsafe.Pointer {
			if !db.alive() || db.h != h {
				return nil
			}
			return unsafe.Pointer(db.c)
		},
		fo: func(c_data unsafe.Pointer) *DB {
			return newDB(h, (*C.alpm_db_t)(c_data))
		},
	}
	runtime.AddCleanup(l, func(c_list **C.alpm_list_t) {
		C.alpm_list_free(*c_list)
	}, &c_list)
	return l
}
