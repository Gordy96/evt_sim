package message

import (
	"container/heap"
	"sync"
)

type queue []*Message

func (pq queue) Len() int { return len(pq) }

func (pq queue) Less(i, j int) bool {
	if pq[i].Priority() < pq[j].Priority() {
		return true
	}
	if pq[i].Priority() > pq[j].Priority() {
		return false
	}

	return pq[i].ID < pq[j].ID
}

func (pq queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *queue) Push(x any) {
	item := x.(*Message)
	*pq = append(*pq, item)
}

func (pq *queue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type PriorityQueue struct {
	q   queue
	mux *sync.Mutex
}

func (pq *PriorityQueue) Len() int {
	pq.mux.Lock()
	defer pq.mux.Unlock()
	return pq.q.Len()
}

func (pq *PriorityQueue) Push(item *Message) {
	pq.mux.Lock()
	defer pq.mux.Unlock()
	heap.Push(&pq.q, item)
}

func (pq *PriorityQueue) Pop() *Message {
	pq.mux.Lock()
	defer pq.mux.Unlock()
	return heap.Pop(&pq.q).(*Message)
}

func NewQueue() *PriorityQueue {
	q := &PriorityQueue{
		mux: new(sync.Mutex),
	}
	heap.Init(&q.q)
	return q
}
