package device

import "github.com/Gordy96/evt-sim/simulation"

type LoraNode struct {
	simulation.ParameterBag
	id string
}

func (l *LoraNode) ID() string {
	return l.id
}

func (l *LoraNode) OnMessage(msg *simulation.Message) {
	//TODO implement me
	panic("implement me")
}

func (l *LoraNode) Init(env simulation.Environment) {
	//TODO implement me
	panic("implement me")
}

func (l *LoraNode) Close() error {
	return nil
}

var _ simulation.Node = (*LoraNode)(nil)
