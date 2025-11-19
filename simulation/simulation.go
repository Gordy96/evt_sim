package simulation

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Gordy96/evt-sim/internal/pq"
	"go.uber.org/zap"
)

type Simulation struct {
	l       *zap.Logger
	pq      *pq.PriorityQueue[*Message]
	nodes   map[string]Node
	init    []Node
	now     time.Time
	elapsed time.Duration
	indexer atomic.Uint64
}

func (s *Simulation) Nodes() map[string]Node {
	return s.nodes
}

func (s *Simulation) SendMessage(msg *Message, delay time.Duration) {
	msg.ID = strconv.FormatUint(s.indexer.Add(1), 10)
	msg.Timestamp = s.now.Add(delay)
	s.pq.Push(msg)
}

func (s *Simulation) Run() {
	s.now = time.Time{}
	start := time.Now()

	s.l.Info("start")

	for _, node := range s.init {
		node.Init(s)
	}

	for s.pq.Len() > 0 {
		msg := s.pq.Pop()
		s.l.Debug("Message", zap.Any("msg", msg))
		s.elapsed += msg.Timestamp.Sub(s.now)
		s.now = msg.Timestamp
		node := s.FindNode(msg.Dst)
		node.OnMessage(msg)
	}
	s.l.Info("finished", zap.Duration("elapsed", time.Since(start)), zap.Duration("simulation_time", s.now.Sub(time.Time{})))

	for _, node := range s.init {
		node.Close()
	}
}

func (s *Simulation) FindNode(id string) Node {
	if node, ok := s.nodes[id]; ok {
		return node
	}

	return nil
}

func (s *Simulation) Now() time.Time {
	return s.now
}

func (s *Simulation) addNode(n Node) error {
	if _, ok := s.nodes[n.ID()]; ok {
		return fmt.Errorf("node with ID %s already exists", n.ID())
	}

	s.nodes[n.ID()] = n

	if composite, ok := n.(CompositeNode); ok {
		for _, sub := range composite.Children() {
			if err := s.addNode(sub); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewSimulation(l *zap.Logger, nodes []Node) (*Simulation, error) {
	s := &Simulation{
		l:     l.Named("simulation"),
		pq:    pq.New[*Message](),
		nodes: make(map[string]Node),
		init:  make([]Node, 0),
	}

	for _, node := range nodes {
		if err := s.addNode(node); err != nil {
			return nil, err
		}

		s.init = append(s.init, node)
	}

	return s, nil
}
