package simulation

type Node interface {
	ID() string
	Param(name string) (any, bool)
	SetParam(name string, value any)
	OnMessage(msg *Message)
	Init(env Environment)
}
