package alpm

// #cgo pkg-config: libalpm
// #include <stdlib.h>
// #include <alpm.h>
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

var (
	ErrHandleCloseFailed = errors.New("alpm: failed to close handle")
	ErrHandleClosed      = errors.New("alpm: handle already closed")
)

type Handle struct {
	c *C.alpm_handle_t
	d runtime.Cleanup
}

func NewHandle(root string, dbpath string) (*Handle, error) {
	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))
	c_dbpath := C.CString(dbpath)
	defer C.free(unsafe.Pointer(c_dbpath))
	c_errno := C.alpm_errno_t(C.ALPM_ERR_OK)

	c_handle := C.alpm_initialize(c_root, c_dbpath, &c_errno)
	if c_handle == nil {
		return nil, Errno(c_errno)
	}

	h := &Handle{c: c_handle}
	h.d = runtime.AddCleanup(h, func(c_handle *C.alpm_handle_t) {
		_ = C.alpm_release(c_handle)
	}, c_handle)
	return h, nil
}

func (h *Handle) Alive() bool {
	return h != nil && h.c != nil
}

func (h *Handle) Close() error {
	if !h.Alive() {
		return nil
	}

	c_ret := C.alpm_release(h.c)
	h.c = nil

	h.d.Stop()
	runtime.KeepAlive(h)

	if c_ret != 0 {
		return ErrHandleCloseFailed
	}
	return nil
}

func (h *Handle) LocalDB() *DB {
	if !h.Alive() {
		return nil
	}

	c_db := C.alpm_get_localdb(h.c)
	runtime.KeepAlive(h)

	return newDB(h, c_db)
}

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

func (h *Handle) errno() error {
	if !h.Alive() {
		return ErrHandleClosed
	}

	errno := Errno(C.alpm_errno(h.c))
	if errno == ErrOK {
		return nil
	}
	return errno
}
