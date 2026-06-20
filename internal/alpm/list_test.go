package alpm

import (
	"slices"
	"testing"
)

func TestNilListFree(t *testing.T) {
	l := (*List[int])(nil)
	l.Free() // should not panic
}

func TestNilListLen(t *testing.T) {
	l := (*List[int])(nil)

	got := l.Len()
	want := 0
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestNilListFront(t *testing.T) {
	l := (*List[int])(nil)

	got := l.Front()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.Front() = %#v, want %#v", got, want)
	}
}

func TestNilListPushBack(t *testing.T) {
	l := (*List[int])(nil)

	const data = 0
	got := l.PushBack(data)
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.PushBack(%#v) = %#v, want %#v", data, got, want)
	}
}

func TestNilListAll(t *testing.T) {
	l := (*List[int])(nil)

	got := 0
	for range l.All() {
		got++
	}

	want := 0
	if got != want {
		t.Errorf("l.All() iteration count = %d, want %d", got, want)
	}
}

func TestZeroListFree(t *testing.T) {
	l := &List[int]{}
	l.Free()

	got := l.v
	want := uint(0)
	if got != want {
		t.Errorf("l.v = %d, want %d", got, want)
	}
}

func TestZeroListLen(t *testing.T) {
	l := &List[int]{}

	got := l.Len()
	want := 0
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestZeroListFront(t *testing.T) {
	l := &List[int]{}

	got := l.Front()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.Front() = %#v, want %#v", got, want)
	}
}

func TestZeroListPushBack(t *testing.T) {
	l := &List[int]{}

	const data = 0
	got := l.PushBack(data)
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.PushBack(%#v) = %#v, want %#v", data, got, want)
	}
}

func TestZeroListAll(t *testing.T) {
	l := &List[int]{}

	got := 0
	for range l.All() {
		got++
	}

	want := 0
	if got != want {
		t.Errorf("l.All() iteration count = %d, want %d", got, want)
	}
}

func TestNilElemData(t *testing.T) {
	e := (*Elem[int])(nil)

	got := e.Data()
	want := 0
	if got != want {
		t.Errorf("e.Data() = %d, want %d", got, want)
	}
}

func TestNilElemNext(t *testing.T) {
	e := (*Elem[int])(nil)

	got := e.Next()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Next() = %#v, want %#v", got, want)
	}
}

func TestNilElemPrev(t *testing.T) {
	e := (*Elem[int])(nil)

	got := e.Prev()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Prev() = %#v, want %#v", got, want)
	}
}

func TestZeroElemData(t *testing.T) {
	e := &Elem[int]{}

	got := e.Data()
	want := 0
	if got != want {
		t.Errorf("e.Data() = %d, want %d", got, want)
	}
}

func TestZeroElemNext(t *testing.T) {
	e := &Elem[int]{}

	got := e.Next()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Next() = %#v, want %#v", got, want)
	}
}

func TestZeroElemPrev(t *testing.T) {
	e := &Elem[int]{}

	got := e.Prev()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Prev() = %#v, want %#v", got, want)
	}
}

func TestBorrowedListFree(t *testing.T) {
	o := newTestListOwner(nil)
	defer o.Free()

	l := newTestBorrowedList(o)

	// Free on a borrowed list should be a no-op
	want := l.v
	l.Free()
	got := l.v
	if got != want {
		t.Errorf("after l.Free(): l.v = %d, want %d", got, want)
	}
}

func TestBorrowedListLen(t *testing.T) {
	s := []int{1, 2, 3}
	o := newTestListOwner(s)
	defer o.Free()

	l := newTestBorrowedList(o)

	got := l.Len()
	want := len(s)
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestBorrowedListLenEmpty(t *testing.T) {
	o := newTestListOwner(nil)
	defer o.Free()

	l := newTestBorrowedList(o)

	got := l.Len()
	want := 0
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestBorrowedListLenAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	o.Free()

	got := l.Len()
	want := 0
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestBorrowedListFront(t *testing.T) {
	s := []int{1, 2, 3}
	o := newTestListOwner(s)
	defer o.Free()

	l := newTestBorrowedList(o)
	e := l.Front()

	if got := e.Prev(); got != nil {
		t.Errorf("[0] e.Prev() = %#v, want nil", got)
	}

	wantPrev := 0
	for i, wantData := range s {
		gotData := e.Data()
		if gotData != wantData {
			t.Errorf("[%d] e.Data() = %d, want %d", i, gotData, wantData)
		}

		gotPrev := e.Prev().Data()
		if gotPrev != wantPrev {
			t.Errorf("[%d] e.Prev().Data() = %d, want %d", i, gotPrev, wantPrev)
		}

		e = e.Next()
		wantPrev = wantData
	}

	if e != nil {
		t.Errorf("[%d] e = %#v, want nil", len(s), e)
	}
}

func TestBorrowedListFrontEmpty(t *testing.T) {
	o := newTestListOwner(nil)
	defer o.Free()

	l := newTestBorrowedList(o)

	got := l.Front()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.Front() = %#v, want %#v", got, want)
	}
}

func TestBorrowedListFrontAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	o.Free()

	got := l.Front()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.Front() = %#v, want %#v", got, want)
	}
}

func TestBorrowedListPushBack(t *testing.T) {
	s := []int{1, 2, 3}
	o := newTestListOwner(s)
	defer o.Free()

	l := newTestBorrowedList(o)

	v := 4
	gotElem := l.PushBack(v)
	wantElem := (*Elem[int])(nil)
	if gotElem != wantElem {
		t.Errorf("l.PushBack(%d) = %#v, want %#v", v, gotElem, wantElem)
	}

	gotAll := slices.Collect(l.All())
	wantAll := s
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedListPushBackEmpty(t *testing.T) {
	o := newTestListOwner(nil)
	defer o.Free()

	l := newTestBorrowedList(o)

	v := 4
	gotElem := l.PushBack(v)
	wantElem := (*Elem[int])(nil)
	if gotElem != wantElem {
		t.Errorf("l.PushBack(%d) = %#v, want %#v", v, gotElem, wantElem)
	}

	gotAll := slices.Collect(l.All())
	wantAll := []int{}
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedListPushBackAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	o.Free()

	v := 4
	gotElem := l.PushBack(v)
	wantElem := (*Elem[int])(nil)
	if gotElem != wantElem {
		t.Errorf("l.PushBack(%d) = %#v, want %#v", v, gotElem, wantElem)
	}

	gotAll := slices.Collect(l.All())
	wantAll := []int{}
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedListAll(t *testing.T) {
	s := []int{1, 2, 3}
	o := newTestListOwner(s)
	defer o.Free()

	l := newTestBorrowedList(o)

	gotAll := slices.Collect(l.All())
	wantAll := s
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedListAllEmpty(t *testing.T) {
	o := newTestListOwner(nil)
	defer o.Free()

	l := newTestBorrowedList(o)

	gotAll := slices.Collect(l.All())
	wantAll := []int{}
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedListAllAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	o.Free()

	gotAll := slices.Collect(l.All())
	wantAll := []int{}
	if !slices.Equal(gotAll, wantAll) {
		t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
	}
}

func TestBorrowedElemDataAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	e := l.Front()
	o.Free()

	got := e.Data()
	want := 0
	if got != want {
		t.Errorf("e.Data() = %d, want %d", got, want)
	}
}

func TestBorrowedElemNextAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	e := l.Front()
	o.Free()

	got := e.Next()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Next() = %#v, want %#v", got, want)
	}
}

func TestBorrowedElemPrevAfterFree(t *testing.T) {
	o := newTestListOwner([]int{1, 2, 3})
	l := newTestBorrowedList(o)
	e := l.Front()
	o.Free()

	got := e.Prev()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("e.Prev() = %#v, want %#v", got, want)
	}
}

func TestOwnedListFree(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	l.PushBack(1)

	gotLen1 := l.Len()
	wantLen1 := 1
	if gotLen1 != wantLen1 {
		t.Errorf("before free: l.Len() = %d, want %d", gotLen1, wantLen1)
	}

	l.Free()

	gotLen2 := l.Len()
	wantLen2 := 0
	if gotLen2 != wantLen2 {
		t.Errorf("after free: l.Len() = %d, want %d", gotLen2, wantLen2)
	}
}

func TestOwnedListFreeEmpty(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	l.Free()

	gotLen := l.Len()
	wantLen := 0
	if gotLen != wantLen {
		t.Errorf("l.Len() = %d, want %d", gotLen, wantLen)
	}
}

func TestOwnedListFreeDouble(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	l.PushBack(1)

	gotLen1 := l.Len()
	wantLen1 := 1
	if gotLen1 != wantLen1 {
		t.Errorf("before free: l.Len() = %d, want %d", gotLen1, wantLen1)
	}

	l.Free()
	l.Free()

	gotLen2 := l.Len()
	wantLen2 := 0
	if gotLen2 != wantLen2 {
		t.Errorf("after free: l.Len() = %d, want %d", gotLen2, wantLen2)
	}
}

func TestOwnedListLen(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	for i := range 3 {
		gotLen := l.Len()
		wantLen := i
		if gotLen != wantLen {
			t.Errorf("l.Len() = %d, want %d", gotLen, wantLen)
		}

		l.PushBack(i)
	}

	// Test the behavior after Free
	l.Free()
	gotLen := l.Len()
	wantLen := 0
	if gotLen != wantLen {
		t.Errorf("l.Len() = %d, want %d", gotLen, wantLen)
	}
}

func TestOwnedListLenEmpty(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	got := l.Len()
	want := 0
	if got != want {
		t.Errorf("l.Len() = %d, want %d", got, want)
	}
}

func TestOwnedListFront(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	s := []int{1, 2, 3}
	for _, v := range s {
		l.PushBack(v)
	}

	e := l.Front()

	if got := e.Prev(); got != nil {
		t.Errorf("[0] e.Prev() = %#v, want nil", got)
	}

	wantPrev := 0
	for i, wantData := range s {
		gotData := e.Data()
		if gotData != wantData {
			t.Errorf("[%d] e.Data() = %d, want %d", i, gotData, wantData)
		}

		gotPrev := e.Prev().Data()
		if gotPrev != wantPrev {
			t.Errorf("[%d] e.Prev().Data() = %d, want %d", i, gotPrev, wantPrev)
		}

		e = e.Next()
		wantPrev = wantData
	}

	if e != nil {
		t.Errorf("[%d] e = %#v, want nil", len(s), e)
	}
}

func TestOwnedListFrontEmpty(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	got := l.Front()
	want := (*Elem[int])(nil)
	if got != want {
		t.Errorf("l.Front() = %#v, want %#v", got, want)
	}
}

func TestOwnedListFrontAfterFree(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	for _, v := range []int{1, 2, 3} {
		l.PushBack(v)
	}

	e := l.Front()
	if e == nil {
		t.Errorf("before free: l.Front() = %#v, want non-nil", e)
	}

	l.Free()

	gotFreeFront := l.Front()
	wantFreeFront := (*Elem[int])(nil)
	if gotFreeFront != wantFreeFront {
		t.Errorf("after free: l.Front() = %#v, want %#v", gotFreeFront, wantFreeFront)
	}

	gotFreeData := e.Data()
	wantFreeData := 0
	if gotFreeData != wantFreeData {
		t.Errorf("after free: e.Data() = %d, want %d", gotFreeData, wantFreeData)
	}

	gotFreePrev := e.Prev()
	wantFreePrev := (*Elem[int])(nil)
	if gotFreePrev != wantFreePrev {
		t.Errorf("after free: e.Prev() = %#v, want %#v", gotFreePrev, wantFreePrev)
	}

	gotFreeNext := e.Next()
	wantFreeNext := (*Elem[int])(nil)
	if gotFreeNext != wantFreeNext {
		t.Errorf("after free: e.Next() = %#v, want %#v", gotFreeNext, wantFreeNext)
	}

	s := []int{1, 2, 3}
	for _, v := range s {
		l.PushBack(v)
	}

	e = l.Front()

	if got := e.Prev(); got != nil {
		t.Errorf("[0] e.Prev() = %#v, want nil", got)
	}

	wantPrev := 0
	for i, wantData := range s {
		gotData := e.Data()
		if gotData != wantData {
			t.Errorf("[%d] e.Data() = %d, want %d", i, gotData, wantData)
		}

		gotPrev := e.Prev().Data()
		if gotPrev != wantPrev {
			t.Errorf("[%d] e.Prev().Data() = %d, want %d", i, gotPrev, wantPrev)
		}

		e = e.Next()
		wantPrev = wantData
	}

	if e != nil {
		t.Errorf("[%d] e = %#v, want nil", len(s), e)
	}
}

func TestOwnedListPushBack(t *testing.T) {
	s := []int{1, 2, 3}
	l := newTestOwnedList()
	defer l.Free()

	if got := slices.Collect(l.All()); len(got) > 0 {
		t.Errorf("l.All() = %#v, want empty", got)
	}

	for i := range s {
		e := l.PushBack(s[i])

		gotData := e.Data()
		wantData := s[i]
		if gotData != wantData {
			t.Errorf("l.PushBack(%d).Data() = %d, want %d", wantData, gotData, wantData)
		}

		gotAll := slices.Collect(l.All())
		wantAll := s[:i+1]
		if !slices.Equal(gotAll, wantAll) {
			t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
		}
	}

	// Test the behavior after Free
	l.Free()

	if got := slices.Collect(l.All()); len(got) > 0 {
		t.Errorf("l.All() = %#v, want empty", got)
	}

	for i := range s {
		e := l.PushBack(s[i])

		gotData := e.Data()
		wantData := s[i]
		if gotData != wantData {
			t.Errorf("l.PushBack(%d).Data() = %d, want %d", wantData, gotData, wantData)
		}

		gotAll := slices.Collect(l.All())
		wantAll := s[:i+1]
		if !slices.Equal(gotAll, wantAll) {
			t.Errorf("l.All() = %#v, want %#v", gotAll, wantAll)
		}
	}
}

func TestOwnedListAll(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	s := []int{1, 2, 3}
	for _, v := range s {
		l.PushBack(v)
	}

	got := make([]int, 0, len(s))
	l.All()(func(v int) bool {
		got = append(got, v)
		return true
	})

	want := s
	if !slices.Equal(got, want) {
		t.Errorf("l.All() = %#v, want %#v", got, want)
	}
}

func TestOwnedListAllStop(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	s := []int{1, 2, 3}
	for _, v := range s {
		l.PushBack(v)
	}

	i := 2
	got := make([]int, 0, len(s))
	l.All()(func(v int) bool {
		got = append(got, v)
		return len(got) < i
	})

	want := s[:i]
	if !slices.Equal(got, want) {
		t.Errorf("l.All()[:%d] = %#v, want %#v", i, got, want)
	}
}

func TestOwnedListAllPushBack(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	s := []int{1, 2, 3}
	l.PushBack(s[0])

	got := make([]int, 0, len(s))
	l.All()(func(v int) bool {
		got = append(got, v)

		if len(got) < len(s) {
			l.PushBack(s[len(got)])
		}

		return true
	})

	want := s
	if !slices.Equal(got, want) {
		t.Errorf("l.All() = %#v, want %#v", got, want)
	}
}

func TestOwnedListAllEmpty(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	l.All()(func(v int) bool {
		t.Errorf("l.All() yielded value %d, want no values (empty list)", v)
		return true
	})
}

func TestOwnedListAllFree(t *testing.T) {
	l := newTestOwnedList()
	defer l.Free()

	s := []int{1, 2, 3}
	for _, v := range s {
		l.PushBack(v)
	}

	i := 2
	got := make([]int, 0, len(s))
	l.All()(func(v int) bool {
		got = append(got, v)
		if len(got) >= i {
			l.Free()
		}
		return true
	})

	want := s[:i]
	if !slices.Equal(got, want) {
		t.Errorf("l.All()[:%d] = %#v, want %#v", i, got, want)
	}
}

func TestOwnedElemNextAfterPushBack(t *testing.T) {
	s := []int{1, 2, 3}
	l := newTestOwnedList()
	defer l.Free()

	e := l.PushBack(s[0])
	if got := e.Next(); got != nil {
		t.Errorf("e.Next() = %#v, want nil", got)
	}

	for i := 1; i < len(s); i++ {
		l.PushBack(s[i])
		e = e.Next()

		gotData := e.Data()
		wantData := s[i]
		if gotData != wantData {
			t.Errorf("e.Data() = %d, want %d", gotData, wantData)
		}
	}

	if got := e.Next(); got != nil {
		t.Errorf("e.Next() = %#v, want nil", got)
	}
}
