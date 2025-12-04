package simulation

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation/message"
)

type Environment interface {
	SendMessage(msg message.Message, delay time.Duration)
	Nodes() map[string]Node
	FindNode(id string) Node
	Now() time.Time
}
