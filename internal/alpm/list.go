package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type List[T any] struct {
	o owner
	c *C.alpm_list_t
	f func(o owner, c_data unsafe.Pointer) T
}

func (l *List[T]) Alive() bool {
	if l == nil || l.c == nil {
		return false
	}
	if !l.o.Alive() {
		l.c = nil
		return false
	}
	return true
}

func (l *List[T]) Len() int {
	if !l.Alive() {
		return 0
	}

	c_count := C.alpm_list_count(l.c)
	runtime.KeepAlive(l)
	return c2goSize(c_count)
}

func (l *List[T]) Front() *Elem[T] {
	return &Elem[T]{o: l, c: l.c}
}

type Elem[T any] struct {
	o *List[T]
	c *C.alpm_list_t
}

func (e *Elem[T]) Alive() bool {
	if e == nil || e.c == nil {
		return false
	}
	if !e.o.Alive() {
		e.c = nil
		return false
	}
	return true
}

func (e *Elem[T]) Data() T {
	if !e.Alive() {
		var zero T
		return zero
	}

	return e.o.f(e.o.o, e.c.data)
}

func (e *Elem[T]) Next() *Elem[T] {
	if !e.Alive() {
		return nil
	}

	c_list := C.alpm_list_next(e.c)
	runtime.KeepAlive(e)
	if c_list == nil {
		return nil
	}

	return &Elem[T]{o: e.o, c: c_list}
}

func (e *Elem[T]) Prev() *Elem[T] {
	if !e.Alive() {
		return nil
	}

	c_list := C.alpm_list_previous(e.c)
	runtime.KeepAlive(e)
	if c_list == nil {
		return nil
	}

	return &Elem[T]{o: e.o, c: c_list}
}

func newPkgList(o *DB, c_list *C.alpm_list_t) *List[*Pkg] {
	return &List[*Pkg]{
		o: o,
		c: c_list,
		f: func(o owner, c_data unsafe.Pointer) *Pkg {
			return &Pkg{o: o.(*DB), c: (*C.alpm_pkg_t)(c_data)}
		},
	}
}

func newDepList(o *Pkg, c_list *C.alpm_list_t) *List[*Depend] {
	return &List[*Depend]{
		o: o,
		c: c_list,
		f: func(o owner, c_data unsafe.Pointer) *Depend {
			return &Depend{o: o.(*Pkg), c: (*C.alpm_depend_t)(c_data)}
		},
	}
}

type owner interface {
	Alive() bool
}
