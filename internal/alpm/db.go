package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import "runtime"

type DB struct {
	h *Handle
	c *C.alpm_db_t
}

func newDB(h *Handle, c_db *C.alpm_db_t) *DB {
	c_handle := C.alpm_db_get_handle(c_db)
	runtime.KeepAlive(h)

	if c_handle != h.c {
		panic("alpm: database handle mismatch")
	}

	return &DB{h: h, c: c_db}
}

func (d *DB) Alive() bool {
	if d == nil || d.c == nil {
		return false
	}
	if !d.h.Alive() {
		d.c = nil
		return false
	}
	return true
}
