package alpm

// #cgo pkg-config: libalpm
// #include <stdlib.h>
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
	return &Elem[T]{o: l, c: *l.c, v: l.v}
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
	return &Elem[T]{o: l, c: c_elem, v: l.v}
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

	return &Elem[T]{o: e.o, c: c_list, v: e.v}
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

	return &Elem[T]{o: e.o, c: c_list, v: e.v}
}

// ------------------------------------------------------------
// Since CGo is not supported within *_test.go files,
// we need to place things for testing here.

type testListOwner struct {
	a bool
	c *C.alpm_list_t
	d runtime.Cleanup
}

func newTestListOwner(s []int) *testListOwner {
	c_list := (*C.alpm_list_t)(nil)
	for _, v := range s {
		c_data := C.malloc(C.size_t(unsafe.Sizeof(v)))
		*(*int)(c_data) = v
		C.alpm_list_append(&c_list, c_data)
	}

	o := &testListOwner{a: true, c: c_list}
	o.d = runtime.AddCleanup(o, freeTestListOwner, c_list)
	return o
}

func (o *testListOwner) alive() bool {
	return o.a
}

func (o *testListOwner) Free() {
	freeTestListOwner(o.c)
	o.a = false
	o.c = nil
	o.d.Stop()
	runtime.KeepAlive(o)
}

func freeTestListOwner(c_list *C.alpm_list_t) {
	for c_elem := c_list; c_elem != nil; c_elem = C.alpm_list_next(c_elem) {
		C.free(c_elem.data)
	}
	C.alpm_list_free(c_list)
}

func newTestBorrowedList(o *testListOwner) *List[int] {
	return &List[int]{
		o: o,
		c: &o.c,
		fo: func(c_data unsafe.Pointer) int {
			return *(*int)(c_data)
		},
	}
}

func newTestOwnedList() *List[int] {
	c_list := (*C.alpm_list_t)(nil)
	l := &List[int]{c: &c_list}

	l.fi = func(v int) unsafe.Pointer {
		defer runtime.KeepAlive(l)

		c_data := C.malloc(C.size_t(unsafe.Sizeof(v)))
		runtime.AddCleanup(l, func(c_data unsafe.Pointer) {
			C.free(c_data)
		}, c_data)
		*(*int)(c_data) = v
		return c_data
	}

	l.fo = func(c_data unsafe.Pointer) int {
		defer runtime.KeepAlive(l)
		return *(*int)(c_data)
	}

	runtime.AddCleanup(l, func(c_list **C.alpm_list_t) {
		C.alpm_list_free(*c_list)
	}, &c_list)

	return l
}
