package simulation

type Node interface {
	ID() string
	OnMessage(msg *Message)
	Init(env Environment)
	Close() error
}
