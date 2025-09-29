package pq

import "container/heap"

type Item interface {
	Priority() int64
}

type PointerConstraint[T any] interface {
	*T
	Item
}

type queue[T Item] []Item

func (pq queue[T]) Len() int { return len(pq) }

func (pq queue[T]) Less(i, j int) bool {
	return pq[i].Priority() < pq[j].Priority()
}

func (pq queue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *queue[T]) Push(x any) {
	item := x.(Item)
	*pq = append(*pq, item)
}

func (pq *queue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type PriorityQueue[T Item] struct {
	q queue[T]
}

func (pq *PriorityQueue[T]) Len() int { return pq.q.Len() }

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(&pq.q, item)
}

func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(&pq.q).(T)
}

func New[T Item]() *PriorityQueue[T] {
	q := &PriorityQueue[T]{}
	heap.Init(&q.q)
	return q
}
