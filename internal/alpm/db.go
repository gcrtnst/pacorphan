package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

type DB struct {
	o *Handle
	c *C.alpm_db_t
}

