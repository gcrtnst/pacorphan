package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

func Version() string {
	return C.GoString(C.alpm_version())
}
