package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Gordy96/evt-sim/nodes"
)

type BaseNode struct {
	id     string
	params map[string]any
	init   func(*BaseNode, *nodes.Simulation)
}

func (b *BaseNode) ID() string {
	return b.id
}

func (b *BaseNode) Param(name string) (any, bool) {
	if b.params == nil {
		return nil, false
	}

	r, ok := b.params[name]
	return r, ok
}

func (b *BaseNode) SetParam(name string, value any) {
	if b.params == nil {
		b.params = make(map[string]any)
	}

	b.params[name] = value
}

func (b *BaseNode) Init(sim *nodes.Simulation) {
	if b.init != nil {
		b.init(b, sim)
	}
}

func (b *BaseNode) HandleMessage(msg *nodes.Message, sim *nodes.Simulation, timestamp time.Time) {
	switch msg.Kind {
	case nodes.KindDelay:
		fmt.Printf("[%s] node '%s' finished sleep, sending message over radio\n", time.Now().Format(time.RFC3339Nano), b.ID())
		sim.SendMessage(&nodes.Message{
			ID:        "some message",
			Src:       b.ID(),
			Dst:       "radio",
			Kind:      nodes.KindMessage,
			Timestamp: timestamp.Add(10 * time.Millisecond),
			Params: map[string]any{
				"payload": "hello world!!!",
			},
		})
	case nodes.KindMessage:
		j, _ := json.Marshal(msg)
		fmt.Printf("[%s] node '%s' received message: %s\n", time.Now().Format(time.RFC3339Nano), b.ID(), j)
	}
}

func main() {
	sim := nodes.NewSimulation([]nodes.Node{
		&BaseNode{
			id: "first",
			init: func(self *BaseNode, sim *nodes.Simulation) {
				sim.Delay(self, 1000*time.Millisecond)
			},
			params: map[string]any{
				"radioFrequency": 433.0,
			},
		},
		&BaseNode{
			id: "second",
			params: map[string]any{
				"radioFrequency": 433.0,
			},
		},
		//third one would never recieve any messages
		&BaseNode{
			id: "third",
			params: map[string]any{
				"radioFrequency": 015.0,
			},
		},
		nodes.NewRadioMedium(100 * time.Millisecond),
	})

	sim.Run()
}
