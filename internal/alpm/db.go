package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

type DB struct {
	o *Handle
	c *C.alpm_db_t
}

func (d *DB) Alive() bool {
	if d == nil {
		return false
	}

	if d.o == nil || d.c == nil || !d.o.Alive() {
		d.o = nil
		d.c = nil
		return false
	}
	return true
}
