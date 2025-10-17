package main

import (
	"time"

	"github.com/Gordy96/evt-sim/modules/radio"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

type BaseNode struct {
	simulation.ParameterBag
	l    *zap.Logger
	id   string
	env  simulation.Environment
	init func(*BaseNode, simulation.Environment)
}

func (b *BaseNode) ID() string {
	return b.id
}

func (b *BaseNode) Init(env simulation.Environment) {
	b.env = env
	if b.init != nil {
		b.init(b, env)
	}
}

func (b *BaseNode) Close() error {
	return nil
}

func (b *BaseNode) OnMessage(msg *simulation.Message) {
	switch msg.Kind {
	case simulation.KindDelay:
		b.l.Info("finished sleep, sending message over radio")
		_, ok := b.GetParam("onWakeDoNothing")
		if ok {
			b.l.Info("node won't do anything on wake")
			return
		}
		b.RemoveParam("busy")
		b.env.SendMessage(&simulation.Message{
			ID:   "some message",
			Src:  b.ID(),
			Dst:  "radio",
			Kind: simulation.KindMessage,
			Params: map[string]any{
				"payload": "hello world!!!",
			},
		}, 10*time.Millisecond)
	case simulation.KindMessage:
		b.l.Info("node received message", zap.String("node", b.id), zap.Any("message", msg))
	}
}

func (b *BaseNode) Delay(t time.Duration) {
	b.env.SendMessage(&simulation.Message{
		ID:   "some message",
		Src:  b.ID(),
		Dst:  b.ID(),
		Kind: simulation.KindDelay,
	}, t)
}

func baseNode(l *zap.Logger, id string, init func(self *BaseNode, env simulation.Environment), params map[string]any) *BaseNode {
	b := &BaseNode{
		id:   id,
		l:    l.Named("base:" + id),
		init: init,
	}

	for k, v := range params {
		b.SetParam(k, v)
	}

	return b
}

func main() {
	logger, _ := zap.NewDevelopment()
	sim := simulation.NewSimulation(logger, []simulation.Node{
		baseNode(
			logger,
			"first",
			func(self *BaseNode, env simulation.Environment) {
				self.Delay(1000 * time.Millisecond)
				self.SetParam("busy", true)
			},
			map[string]any{
				"radioFrequency": 433.0,
			},
		),
		baseNode(
			logger,
			"second",
			nil,
			map[string]any{
				"radioFrequency": 433.0,
			},
		),
		//going to miss because it's busy
		baseNode(
			logger,
			"fourth",
			func(self *BaseNode, env simulation.Environment) {
				self.Delay(5000 * time.Millisecond)
				self.SetParam("busy", true)
			},
			map[string]any{
				"radioFrequency": 433.0,
			},
		),
		//third one would never receive any messages
		baseNode(
			logger,
			"third",
			nil,
			map[string]any{
				"radioFrequency": 915.0,
			},
		),
		//radio medium is also a node that can recieve messages
		//think of it as 'aether' anything that has radio can talk to it,
		//then it decides what simulation should recieve message (effectively duplicating messages)
		//based on node parameters (potentially simulation can have ports/interfaces, that would hold parameters/talk to 'aether')
		radio.NewRadioMedium(logger, 100*time.Millisecond),
	})

	sim.Run()
}
