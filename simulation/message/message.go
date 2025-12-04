package message

import (
	"time"
)

type Kind string

type Message struct {
	ID        string
	Src       string
	Dst       string
	Kind      Kind
	Timestamp time.Time
	Params    Parameters `builder:"passthrough"`
}

func (m *Message) Priority() int64 {
	return m.Timestamp.UnixNano()
}

func (m *Message) Builder() Builder {
	c := *m
	c.Params = Parameters{}
	c.Params.Merge(m.Params)
	return Builder{
		msg: c,
	}
}

type Builder struct {
	msg Message
}

func (m Builder) WithID(id string) Builder {
	m.msg.ID = id
	return m
}

func (m Builder) WithSrc(src string) Builder {
	m.msg.Src = src
	return m
}

func (m Builder) WithDst(dst string) Builder {
	m.msg.Dst = dst
	return m
}

func (m Builder) WithKind(kind Kind) Builder {
	m.msg.Kind = kind
	return m
}

func (m Builder) WithTimestamp(t time.Time) Builder {
	m.msg.Timestamp = t
	return m
}

func (m Builder) WithParams(p Parameters) Builder {
	m.msg.Params.Merge(p)

	return m
}

func (m Builder) Build() Message {
	return m.msg
}
