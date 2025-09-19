package nodes

import (
	"container/heap"
	"fmt"
	"time"

	"github.com/Gordy96/evt-sim/internal"
)

type MessageKind uint64

const (
	KindDelay MessageKind = iota
	KindMessage
)

type Message struct {
	ID        string
	Src       string
	Dst       string
	Kind      MessageKind
	Timestamp time.Time
	Params    map[string]any
}

func (m *Message) Priority() int64 {
	return m.Timestamp.UnixNano()
}

type Node interface {
	ID() string
	Param(name string) (any, bool)
	SetParam(name string, value any)
	HandleMessage(msg *Message, sim *Simulation, timestamp time.Time)
	Init(sim *Simulation)
}

type Simulation struct {
	pq      internal.PriorityQueue[*Message]
	nodes   map[string]Node
	simTime time.Time
}

func (s *Simulation) Nodes() map[string]Node {
	return s.nodes
}

func (s *Simulation) SendMessage(msg *Message) {
	heap.Push(&s.pq, msg)
}

func (s *Simulation) Log(format string, a ...any) {
	args := []any{time.Now().Format(time.RFC3339Nano), s.simTime.Format(time.RFC3339Nano)}
	args = append(args, a...)
	fmt.Printf("[%s | %s] "+format+"\n", args...)
}

func (s *Simulation) Delay(node Node, delay time.Duration) {
	s.Log("node '%s' goes to sleep for %s", node.ID(), delay)
	heap.Push(&s.pq, &Message{
		ID:        "",
		Src:       node.ID(),
		Dst:       node.ID(),
		Kind:      KindDelay,
		Timestamp: s.simTime.Add(delay),
	})
}

func (s *Simulation) Run() {
	heap.Init(&s.pq)

	s.simTime = time.Now()

	s.Log("start")

	for _, node := range s.nodes {
		node.Init(s)
	}

	for s.pq.Len() > 0 {
		msg := heap.Pop(&s.pq).(*Message)
		s.simTime = msg.Timestamp
		node := s.nodes[msg.Dst]
		node.HandleMessage(msg, s, msg.Timestamp)
	}
	s.Log("finish")
}

func (s *Simulation) Time() time.Time {
	return s.simTime
}

func NewSimulation(nodes []Node) *Simulation {
	n := make(map[string]Node)

	for _, node := range nodes {
		n[node.ID()] = node
	}

	return &Simulation{
		pq:    internal.PriorityQueue[*Message]{},
		nodes: n,
	}
}
