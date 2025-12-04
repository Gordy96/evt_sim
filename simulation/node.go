package simulation

import "github.com/Gordy96/evt-sim/simulation/message"

type Node interface {
	ID() string
	OnMessage(msg message.Message)
	Init(env Environment)
	Close() error
	Parent() Node
}
