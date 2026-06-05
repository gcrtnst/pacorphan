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
	p *Pkg
	c *C.alpm_depend_t
}

func newDep(p *Pkg, c_dep *C.alpm_depend_t) *Depend {
	return &Depend{p: p, c: c_dep}
}

func (d *Depend) Alive() bool {
	if d == nil || d.c == nil {
		return false
	}
	if !d.p.Alive() {
		d.c = nil
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
