package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import "runtime"

func (h *Handle) Root() string {
	if !h.Alive() {
		return ""
	}

	root := C.GoString(C.alpm_option_get_root(h.c))
	runtime.KeepAlive(h)
	return root
}

func (h *Handle) DBPath() string {
	if !h.Alive() {
		return ""
	}

	dbpath := C.GoString(C.alpm_option_get_dbpath(h.c))
	runtime.KeepAlive(h)
	return dbpath
}
