package queue

import (
	"sort"
	"testing"
)

type testSortable struct {
	priority int
}

func (ts *testSortable) Priority(other interface{}) bool {
	if t, ok := other.(*testSortable); ok {
		return ts.priority > t.priority
	}

	return false
}

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()
	l := 10
	e := make([]int, 0, l)

	for i := 0; i < l; i++ {
		p := (i + 1) * 5
		e = append(e, p)

		pq.Push(&testSortable{priority: p})
	}

	sort.Ints(e)

	for i := l - 1; i >= 0; i-- {
		r := pq.Pop()

		if r.(*testSortable).priority != e[i] {
			t.Errorf("incorrect result: expected: %v, got: %v", e[i], r.(*testSortable).priority)
		}
	}

	if pq.Size() != 0 {
		t.Errorf("incorrect result: expected: %v, got: %v", 0, pq.Size())
	}
}
