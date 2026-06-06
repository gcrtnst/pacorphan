package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import "runtime"

type Pkg struct {
	h *Handle
	d *DB
	c *C.alpm_pkg_t
}

func newPkg(h *Handle, d *DB, c_pkg *C.alpm_pkg_t) *Pkg {
	if h != d.h {
		panic("alpm: database handle mismatch")
	}

	c_handle := C.alpm_pkg_get_handle(c_pkg)
	runtime.KeepAlive(h)
	runtime.KeepAlive(d)

	if c_handle != h.c {
		panic("alpm: package handle mismatch")
	}

	c_db := C.alpm_pkg_get_db(c_pkg)
	runtime.KeepAlive(h)
	runtime.KeepAlive(d)

	if c_db != d.c {
		panic("alpm: package database mismatch")
	}

	return &Pkg{h: h, d: d, c: c_pkg}
}

func (p *Pkg) Alive() bool {
	if p == nil || p.c == nil {
		return false
	}
	if !p.d.Alive() || !p.h.Alive() {
		p.c = nil
		return false
	}
	return true
}

func (p *Pkg) Reason() PkgReason {
	if !p.Alive() {
		return PkgReasonUnknown
	}

	c_reason := C.alpm_pkg_get_reason(p.c)
	runtime.KeepAlive(p)
	return PkgReason(c_reason)
}

type PkgReason int

const (
	PkgReasonExplicit PkgReason = C.ALPM_PKG_REASON_EXPLICIT
	PkgReasonDepend   PkgReason = C.ALPM_PKG_REASON_DEPEND
	PkgReasonUnknown  PkgReason = C.ALPM_PKG_REASON_UNKNOWN
)
