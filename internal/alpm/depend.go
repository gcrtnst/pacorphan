package alpm

// #cgo pkg-config: libalpm
// #include <stdlib.h>
// #include <alpm.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type Depend struct {
	o *Pkg
	c *C.alpm_depend_t
}

func (d *Depend) Alive() bool {
	if d == nil {
		return false
	}

	if d.c == nil || d.o == nil || !d.o.Alive() {
		d.c = nil
		d.o = nil
		return false
	}
	return true
}

func (d *Depend) String() string {
	if !d.Alive() {
		return ""
	}

	c_string := C.alpm_dep_compute_string(d.c)
	runtime.KeepAlive(d)
	defer C.free(unsafe.Pointer(c_string))
	return C.GoString(c_string)
}
