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

func (d *Depend) alive() bool {
	if d == nil || d.c == nil {
		return false
	}
	if !d.p.alive() {
		d.c = nil
		return false
	}
	return true
}

func (d *Depend) String() string {
	defer runtime.KeepAlive(d)
	if !d.alive() {
		return ""
	}

	c_string := C.alpm_dep_compute_string(d.c)
	defer C.free(unsafe.Pointer(c_string))
	return C.GoString(c_string)
}

func newDepList(p *Pkg, c_list *C.alpm_list_t) *List[*Depend] {
	return &List[*Depend]{
		o:  p,
		c:  &c_list,
		fi: nil,
		fo: func(c_data unsafe.Pointer) *Depend {
			return newDep(p, (*C.alpm_depend_t)(c_data))
		},
	}
}
