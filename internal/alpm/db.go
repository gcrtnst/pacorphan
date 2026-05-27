package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import "runtime"

type DB struct {
	h *Handle
	c *C.alpm_db_t
}

func (h *Handle) LocalDB() *DB {
	c_db := C.alpm_get_localdb(h.c)
	runtime.KeepAlive(h)

	return &DB{h: h, c: c_db}
}
