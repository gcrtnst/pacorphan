package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import (
	"iter"
	"runtime"
	"unsafe"
)

type List[T any] struct {
	o owner
	c *C.alpm_list_t
	f func(c_data unsafe.Pointer) T
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
	defer runtime.KeepAlive(l)
	if !l.Alive() {
		return 0
	}

	c_count := C.alpm_list_count(l.c)
	return c2goSize(c_count)
}

func (l *List[T]) Front() *Elem[T] {
	return &Elem[T]{o: l, c: l.c}
}

func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		e := l.Front()
		for e != nil {
			if !yield(e.Data()) {
				break
			}
			e = e.Next()
		}
	}
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
	defer runtime.KeepAlive(e)
	if !e.Alive() {
		var zero T
		return zero
	}

	return e.o.f(e.c.data)
}

func (e *Elem[T]) Next() *Elem[T] {
	defer runtime.KeepAlive(e)
	if !e.Alive() {
		return nil
	}

	c_list := C.alpm_list_next(e.c)
	if c_list == nil {
		return nil
	}

	return &Elem[T]{o: e.o, c: c_list}
}

func (e *Elem[T]) Prev() *Elem[T] {
	defer runtime.KeepAlive(e)
	if !e.Alive() {
		return nil
	}

	c_list := C.alpm_list_previous(e.c)
	if c_list == nil {
		return nil
	}

	return &Elem[T]{o: e.o, c: c_list}
}

func newPkgList(h *Handle, d *DB, c_list *C.alpm_list_t) *List[*Pkg] {
	return &List[*Pkg]{
		o: d,
		c: c_list,
		f: func(c_data unsafe.Pointer) *Pkg {
			return newPkg(h, d, (*C.alpm_pkg_t)(c_data))
		},
	}
}

func newDepList(p *Pkg, c_list *C.alpm_list_t) *List[*Depend] {
	return &List[*Depend]{
		o: p,
		c: c_list,
		f: func(c_data unsafe.Pointer) *Depend {
			return newDep(p, (*C.alpm_depend_t)(c_data))
		},
	}
}

type owner interface {
	Alive() bool
}
