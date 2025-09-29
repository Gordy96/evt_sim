package main

import (
	"encoding/json"
	"time"

	"github.com/Gordy96/evt-sim/modules/radio"
	"github.com/Gordy96/evt-sim/simulation"
)

type BaseNode struct {
	id     string
	params map[string]any
	init   func(*BaseNode, *simulation.Simulation)
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

func (b *BaseNode) Init(sim *simulation.Simulation) {
	if b.init != nil {
		b.init(b, sim)
	}
}

func (b *BaseNode) HandleMessage(msg *simulation.Message, sim *simulation.Simulation, timestamp time.Time) {
	switch msg.Kind {
	case simulation.KindDelay:
		sim.Log("node '%s' finished sleep, sending message over radio", b.ID())
		_, ok := b.Param("onWakeDoNothing")
		if ok {
			sim.Log("node '%s' won't do anything on wake", b.ID())
			return
		}
		delete(b.params, "busy")
		sim.SendMessage(&simulation.Message{
			ID:        "some message",
			Src:       b.ID(),
			Dst:       "radio",
			Kind:      simulation.KindMessage,
			Timestamp: timestamp.Add(10 * time.Millisecond),
			Params: map[string]any{
				"payload": "hello world!!!",
			},
		})
	case simulation.KindMessage:
		j, _ := json.Marshal(msg)
		sim.Log("node '%s' received message: %s", b.ID(), j)
	}
}

func main() {
	sim := simulation.NewSimulation([]simulation.Node{
		&BaseNode{
			id: "first",
			init: func(self *BaseNode, sim *simulation.Simulation) {
				sim.Delay(self, 1000*time.Millisecond)
				self.SetParam("busy", true)
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
		//going to miss because it's busy
		&BaseNode{
			id: "fourth",
			init: func(self *BaseNode, sim *simulation.Simulation) {
				sim.Delay(self, 5000*time.Millisecond)
				self.SetParam("busy", true)
			},
			params: map[string]any{
				"radioFrequency": 433.0,
			},
		},
		//third one would never recieve any messages
		&BaseNode{
			id: "third",
			params: map[string]any{
				"radioFrequency": 915.0,
			},
		},
		//radio medium is also a node that can recieve messages
		//think of it as 'aether' anything that has radio can talk to it,
		//then it decides what simulation should recieve message (effectively duplicating messages)
		//based on node parameters (potentially simulation can have ports/interfaces, that would hold parameters/talk to 'aether')
		radio.NewRadioMedium(100 * time.Millisecond),
	})

	sim.Run()
}
