package embedded

import (
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/Gordy96/evt-sim/simulation/message"
)

type bufferNodeWrapper struct {
	name string
	node simulation.Node
	src  *EmbeddedDevice
	buf  []byte
}

func (b bufferNodeWrapper) Name() string {
	return b.name
}

func (b bufferNodeWrapper) Read(buf []byte) (int, error) {
	if len(b.buf) == 0 {
		return 0, nil
	}

	n := copy(buf, b.buf)

	copy(b.buf, b.buf[n:])
	b.buf = b.buf[:len(b.buf)-n]

	return n, nil
}

func (b bufferNodeWrapper) Write(buf []byte) (n int, err error) {
	var c = make([]byte, len(buf))
	copy(c, buf)

	parameters := message.Parameters{}
	builder := message.Builder{}

	b.src.env.SendMessage(builder.
		WithSrc(b.src.ID()).
		WithDst(b.node.ID()).
		WithKind("wire/payload").
		WithParams(parameters.WithBytes("payload", c)).
		Build(), 0)

	return len(buf), nil
}
