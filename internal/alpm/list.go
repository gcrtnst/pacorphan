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
	o  owner
	c  **C.alpm_list_t
	fi func(T) unsafe.Pointer
	fo func(unsafe.Pointer) T
}

func (l *List[T]) alive() bool {
	if l == nil || l.c == nil {
		return false
	}
	if l.o != nil && !l.o.alive() {
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
}

func (e *Elem[T]) alive() bool {
	if e == nil || e.c == nil {
		return false
	}
	if !e.o.alive() {
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

func newPkgList(h *Handle, d *DB, c_list *C.alpm_list_t) *List[*Pkg] {
	return &List[*Pkg]{
		o:  d,
		c:  &c_list,
		fi: nil,
		fo: func(c_data unsafe.Pointer) *Pkg {
			return newPkg(h, d, (*C.alpm_pkg_t)(c_data))
		},
	}
}

func newDepList(p *Pkg, c_list *C.alpm_list_t) *List[*Depend] {
	return &List[*Depend]{
		o:  p,
		c:  &c_list,
		fi: nil,
		fo: func(c_data unsafe.Pointer) *Depend {
			return newDep(p, (*C.alpm_depend_t)(c_data))
		},
	}
}

func newDBList(h *Handle) *List[*DB] {
	c_list := (*C.alpm_list_t)(nil)
	l := &List[*DB]{
		c: &c_list,
		fi: func(db *DB) unsafe.Pointer {
			if !db.alive() || db.h != h {
				return nil
			}
			return unsafe.Pointer(db.c)
		},
		fo: func(c_data unsafe.Pointer) *DB {
			return newDB(h, (*C.alpm_db_t)(c_data))
		},
	}
	runtime.AddCleanup(l, func(c_list **C.alpm_list_t) {
		C.alpm_list_free(*c_list)
	}, &c_list)
	return l
}

type owner interface {
	alive() bool
}
