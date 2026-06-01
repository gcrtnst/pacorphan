package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import "runtime"

type Pkg struct {
	o *DB
	c *C.alpm_pkg_t
}

func (p *Pkg) Alive() bool {
	if p == nil {
		return false
	}

	if p.c == nil || p.o == nil || !p.o.Alive() {
		p.c = nil
		p.o = nil
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
