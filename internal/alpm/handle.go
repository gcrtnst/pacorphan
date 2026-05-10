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

var ErrCloseHandle = errors.New("failed to close handle")

type Handle struct {
	c       *C.alpm_handle_t
	cleanup runtime.Cleanup
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
	h.cleanup = runtime.AddCleanup(h, func(h *C.alpm_handle_t) {
		_ = C.alpm_release(h)
	}, c_handle)
	return h, nil
}

func (h *Handle) Close() error {
	c_ret := C.alpm_release(h.c)

	h.cleanup.Stop()
	runtime.KeepAlive(h) // Ensures that cleanup is reachable across the call to Stop

	if c_ret != 0 {
		return ErrCloseHandle
	}
	return nil
}
