package simulation

import "time"

type Environment interface {
	SendMessage(msg *Message)
	Nodes() map[string]Node
	Now() time.Time
}
