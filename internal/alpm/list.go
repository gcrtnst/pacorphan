package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"
import (
	"iter"
	"runtime"
	"unsafe"
)

type owner interface {
	alive() bool
}

type List[T any] struct {
	o  owner
	c  **C.alpm_list_t
	v  uint
	fi func(T) unsafe.Pointer
	fo func(unsafe.Pointer) T
}

func (l *List[T]) alive() bool {
	if l == nil || l.c == nil {
		return false
	}
	if l.o != nil && !l.o.alive() {
		*l.c = nil
		l.c = nil
		return false
	}
	return true
}

func (l *List[T]) Free() {
	defer runtime.KeepAlive(l)
	if !l.alive() || l.o != nil {
		return
	}

	l.v++
	C.alpm_list_free(*l.c)
	*l.c = nil
}

func (l *List[T]) Len() int {
	defer runtime.KeepAlive(l)
	if !l.alive() {
		return 0
	}

	c_count := C.alpm_list_count(*l.c)
	return c2goSize(c_count)
}

func (l *List[T]) Front() *Elem[T] {
	if !l.alive() || *l.c == nil {
		return nil
	}
	return &Elem[T]{o: l, c: *l.c}
}

func (l *List[T]) PushBack(data T) *Elem[T] {
	defer runtime.KeepAlive(l)
	defer runtime.KeepAlive(data)
	if !l.alive() || l.o != nil || l.fi == nil {
		return nil
	}

	c_data := l.fi(data)
	if c_data == nil {
		return nil
	}

	c_elem := C.alpm_list_append(l.c, c_data)
	return &Elem[T]{o: l, c: c_elem}
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
	v uint
}

func (e *Elem[T]) alive() bool {
	if e == nil || e.c == nil {
		return false
	}
	if !e.o.alive() || e.o.v != e.v {
		e.c = nil
		return false
	}
	return true
}

func (e *Elem[T]) Data() T {
	defer runtime.KeepAlive(e)
	if !e.alive() {
		var zero T
		return zero
	}

	return e.o.fo(e.c.data)
}

func (e *Elem[T]) Next() *Elem[T] {
	defer runtime.KeepAlive(e)
	if !e.alive() {
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
	if !e.alive() {
		return nil
	}

	c_list := C.alpm_list_previous(e.c)
	if c_list == nil {
		return nil
	}

	return &Elem[T]{o: e.o, c: c_list}
}
