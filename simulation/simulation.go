package simulation

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Gordy96/evt-sim/simulation/message"
	"go.uber.org/zap"
)

type Simulation struct {
	l             *zap.Logger
	pq            *message.PriorityQueue
	nodes         map[string]Node
	init          []Node
	now           time.Time
	elapsed       time.Duration
	indexer       atomic.Uint64
	wg            sync.WaitGroup
	readyLock     sync.WaitGroup
	realtime      bool
	lastEventTime time.Time
}

func (s *Simulation) Nodes() map[string]Node {
	return s.nodes
}

func (s *Simulation) SendMessage(msg message.Message, delay time.Duration) {
	msg.ID = strconv.FormatUint(s.indexer.Add(1), 10)
	msg.Timestamp = s.now.Add(delay)
	if s.realtime {
		s.wg.Add(1)
		s.readyLock.Wait()
		time.AfterFunc(delay, func() {
			defer s.wg.Done()

			s.lastEventTime = time.Now()

			s.l.Debug("Message", zap.Any("msg", &msg))
			node := s.FindNode(msg.Dst)
			node.OnMessage(msg)
		})
	} else {
		s.pq.Push(&msg)
	}
}

func (s *Simulation) Run() {
	s.now = time.Time{}
	start := time.Now()

	s.l.Info("start")

	s.readyLock.Add(1)

	for _, node := range s.init {
		node.Init(s)
	}

	s.readyLock.Done()

	if !s.realtime {
		deadline := time.Now().Add(time.Second)
		s.lastEventTime = time.Now()

		for time.Now().Before(deadline) {
			if s.pq.Len() > 0 {
				msg := s.pq.Pop()
				s.l.Debug("Message", zap.Any("msg", msg))
				s.elapsed += msg.Timestamp.Sub(s.now)
				s.now = msg.Timestamp
				node := s.FindNode(msg.Dst)
				node.OnMessage(*msg)

				deadline = time.Now().Add(time.Second)
				s.lastEventTime = time.Now()
			}
		}
	}

	time.Sleep(time.Second)
	s.wg.Wait()

	s.l.Info("finished",
		zap.Duration("elapsed", s.lastEventTime.Sub(start)),
		zap.Duration("simulation_time", s.now.Sub(time.Time{})),
	)

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

func NewSimulation(l *zap.Logger, nodes []Node, realtime bool) (*Simulation, error) {
	s := &Simulation{
		l:        l.Named("simulation"),
		pq:       message.NewQueue(),
		nodes:    make(map[string]Node),
		init:     make([]Node, 0),
		realtime: realtime,
	}

	for _, node := range nodes {
		if err := s.addNode(node); err != nil {
			return nil, err
		}

		s.init = append(s.init, node)
	}

	return s, nil
}
