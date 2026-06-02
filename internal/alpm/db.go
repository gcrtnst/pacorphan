package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

type DB struct {
	o *Handle
	c *C.alpm_db_t
}

func (d *DB) Alive() bool {
	if d == nil || d.c == nil {
		return false
	}
	if !d.o.Alive() {
		d.c = nil
		return false
	}
	return true
}
