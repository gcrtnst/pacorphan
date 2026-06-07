package alpm

// #cgo pkg-config: libalpm
// #include <stdlib.h>
// #include <alpm.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type Pkg struct {
	h *Handle
	d *DB
	c *C.alpm_pkg_t
}

func newPkg(h *Handle, d *DB, c_pkg *C.alpm_pkg_t) *Pkg {
	defer runtime.KeepAlive(h)
	defer runtime.KeepAlive(d)

	if h != d.h {
		panic("alpm: database handle mismatch")
	}

	c_handle := C.alpm_pkg_get_handle(c_pkg)
	if c_handle != h.c {
		panic("alpm: package handle mismatch")
	}

	c_db := C.alpm_pkg_get_db(c_pkg)
	if c_db != d.c {
		panic("alpm: package database mismatch")
	}

	return &Pkg{h: h, d: d, c: c_pkg}
}

func (p *Pkg) alive() bool {
	if p == nil || p.c == nil {
		return false
	}
	if !p.d.alive() || !p.h.alive() {
		p.c = nil
		return false
	}
	return true
}

func (p *Pkg) Name() string {
	defer runtime.KeepAlive(p)
	if !p.alive() {
		return ""
	}

	c_string := C.alpm_pkg_get_name(p.c)
	return C.GoString(c_string)
}

func (p *Pkg) Version() string {
	defer runtime.KeepAlive(p)
	if !p.alive() {
		return ""
	}

	c_string := C.alpm_pkg_get_version(p.c)
	return C.GoString(c_string)
}

func (p *Pkg) Depends() *List[*Depend] {
	defer runtime.KeepAlive(p)
	if !p.alive() {
		return nil
	}

	c_list := C.alpm_pkg_get_depends(p.c)
	return newDepList(p, c_list)
}

func (p *Pkg) OptDepends() *List[*Depend] {
	defer runtime.KeepAlive(p)
	if !p.alive() {
		return nil
	}

	c_list := C.alpm_pkg_get_optdepends(p.c)
	return newDepList(p, c_list)
}

func (p *Pkg) Reason() PkgReason {
	defer runtime.KeepAlive(p)
	if !p.alive() {
		return PkgReasonUnknown
	}

	c_reason := C.alpm_pkg_get_reason(p.c)
	return PkgReason(c_reason)
}

type PkgReason int

const (
	PkgReasonExplicit PkgReason = C.ALPM_PKG_REASON_EXPLICIT
	PkgReasonDepend   PkgReason = C.ALPM_PKG_REASON_DEPEND
	PkgReasonUnknown  PkgReason = C.ALPM_PKG_REASON_UNKNOWN
)

func CompareVersion(a, b string) int {
	c_a := C.CString(a)
	defer C.free(unsafe.Pointer(c_a))
	c_b := C.CString(b)
	defer C.free(unsafe.Pointer(c_b))

	c_ret := C.alpm_pkg_vercmp(c_a, c_b)
	return c2goInt(c_ret)
}

func newPkgList(h *Handle, d *DB, c_list *C.alpm_list_t) *List[*Pkg] {
	return &List[*Pkg]{
		o:  d,
		c:  &c_list,
		fi: nil,
		fo: func(c_data unsafe.Pointer) *Pkg {
			return newPkg(h, d, (*C.alpm_pkg_t)(c_data))
		},
	}
}
