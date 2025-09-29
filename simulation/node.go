package simulation

type Node interface {
	ID() string
	GetParam(name string) (any, bool)
	OnMessage(msg *Message)
	Init(env Environment)
}
