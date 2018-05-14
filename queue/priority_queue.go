package queue

import "container/heap"

type (
	// Heapable reflects an interface that an item must implement in order to
	// be placed into a priority queue.
	Heapable interface {
		Priority(other interface{}) bool
	}

	// items is the internal container for a priority queue. It's implementation
	// should not be exposed publically.
	items []Heapable

	// PriorityQueue implements a priority queue of Heapable items. Each item
	// must implement the Heapable interface and is responsible for giving each
	// item priority respective to every other item. The main benefit of this
	// implementation is to abstract away any direct heap usage.
	PriorityQueue struct {
		queue *items
	}
)

// NewPriorityQueue returns a reference to a new initialized PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{queue: &items{}}

	heap.Init(pq.queue)
	return pq
}

// Push adds a Sortable item to the priority queue.
func (pq *PriorityQueue) Push(s Heapable) {
	heap.Push(pq.queue, s)
}

// Pop removes a Sortable item from the priority queue with the highest
// priority and returns it.
func (pq *PriorityQueue) Pop() Heapable {
	return heap.Pop(pq.queue).(Heapable)
}

// Size returns the size of the priority queue.
func (pq *PriorityQueue) Size() int {
	return pq.queue.Len()
}

// Len implements the sort.Interface.
func (pq items) Len() int {
	return len(pq)
}

// Less implements the sort.Interface.
func (pq items) Less(i, j int) bool {
	return pq[i].Priority(pq[j])
}

// Swap implements the sort.Interface.
func (pq items) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push implements the heap interface.
func (pq *items) Push(x interface{}) {
	item := x.(Heapable)
	*pq = append(*pq, item)
}

// Pop implements the heap interface.
func (pq *items) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
