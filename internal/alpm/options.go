package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

func (h *Handle) Root() string {
	c_root := C.alpm_option_get_root(h.c)
	return C.GoString(c_root)
}

func (h *Handle) DBPath() string {
	c_dbpath := C.alpm_option_get_dbpath(h.c)
	return C.GoString(c_dbpath)
}
