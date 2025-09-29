package simulation

import (
	"time"

	"github.com/Gordy96/evt-sim/internal/pq"
	"go.uber.org/zap"
)

type Simulation struct {
	l       *zap.Logger
	pq      *pq.PriorityQueue[*Message]
	nodes   map[string]Node
	now     time.Time
	elapsed time.Duration
}

func (s *Simulation) Nodes() map[string]Node {
	return s.nodes
}

func (s *Simulation) SendMessage(msg *Message) {
	s.pq.Push(msg)
}

func (s *Simulation) Run() {
	s.now = time.Now()

	start := s.now

	s.l.Info("start")

	for _, node := range s.nodes {
		node.Init(s)
	}

	for s.pq.Len() > 0 {
		msg := s.pq.Pop()
		s.elapsed += msg.Timestamp.Sub(s.now)
		s.now = msg.Timestamp
		node := s.nodes[msg.Dst]
		node.OnMessage(msg)
	}
	s.l.Info("finished", zap.Duration("elapsed", time.Since(start)))
}

func (s *Simulation) Now() time.Time {
	return s.now
}

func NewSimulation(nodes []Node) *Simulation {
	n := make(map[string]Node)

	for _, node := range nodes {
		n[node.ID()] = node
	}

	return &Simulation{
		pq:    pq.New[*Message](),
		nodes: n,
	}
}
