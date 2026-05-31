package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

type DB struct {
	h *Handle
	c *C.alpm_db_t
}

