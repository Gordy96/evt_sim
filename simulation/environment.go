package simulation

import "time"

type Environment interface {
	SendMessage(msg *Message, delay time.Duration)
	Nodes() map[string]Node
	FindNode(id string) Node
	Now() time.Time
}
