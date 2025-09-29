package simulation

import "time"

type MessageKind uint64

const (
	KindDelay MessageKind = iota
	KindMessage
)

type Message struct {
	ID        string
	Src       string
	Dst       string
	Kind      MessageKind
	Timestamp time.Time
	Params    map[string]any
}

func (m *Message) Priority() int64 {
	return m.Timestamp.UnixNano()
}
