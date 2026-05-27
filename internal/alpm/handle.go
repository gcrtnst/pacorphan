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

var ErrCloseHandle = errors.New("alpm: failed to close handle")

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

func (h *Handle) Close() error {
	if h.c == nil {
		return nil
	}

	c_ret := C.alpm_release(h.c)
	h.c = nil

	h.d.Stop()
	runtime.KeepAlive(h) // Ensures that cleanup is reachable across the call to Stop

	if c_ret != 0 {
		return ErrCloseHandle
	}
	return nil
}
