package embedded

import (
	"io"
	"slices"
	"time"

	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/Gordy96/evt-sim/simulation"
)

var _ simulation.Node = (*EmbeddedDevice)(nil)

type bufferNodeWrapper struct {
	node simulation.Node
	src  *EmbeddedDevice
	buf  []byte
}

func (b bufferNodeWrapper) Name() string {
	return b.node.ID()
}

func (b bufferNodeWrapper) Read(buf []byte) (n int, err error) {
	//TODO: destructive read?
	return copy(buf, b.buf), nil
}

func (b bufferNodeWrapper) Write(buf []byte) (n int, err error) {
	b.src.env.SendMessage(&simulation.Message{
		ID:   "",
		Src:  b.src.ID(),
		Dst:  b.node.ID(),
		Kind: "start_sending",
		Params: map[string]any{
			"payload": buf,
		},
	}, 0)

	return len(buf), nil
}

type EmbeddedDevice struct {
	simulation.ParameterBag
	id     string
	env    simulation.Environment
	app    device.Application
	ports  map[string]bufferNodeWrapper
	radios []simulation.Node
}

func (e *EmbeddedDevice) ID() string {
	return e.id
}

func (e *EmbeddedDevice) OnMessage(msg *simulation.Message) {
	switch msg.Kind {
	case "delay_interrupt":
		key := msg.Params["key"].(string)
		e.app.TriggerTimeInterrupt(key)
	case "received_radio_message":
		if port, ok := e.ports[msg.Src]; ok {
			ipl, ok := msg.Params["payload"]
			if ok {
				payload := ipl.([]byte)
				if len(payload) > cap(port.buf) {
					//since port wrapper is value - update buffer pointer in map
					port.buf = slices.Grow(port.buf, len(payload)-cap(port.buf))
					e.ports[msg.Src] = port
				}

				copy(port.buf, payload)
				e.app.TriggerPortInterrupt(msg.Src)
			}
		}
	}
}

func (e *EmbeddedDevice) schedule(key string, timeMS int) {
	e.env.SendMessage(&simulation.Message{
		ID:   "",
		Src:  e.ID(),
		Dst:  e.ID(),
		Kind: "delay_interrupt",
		Params: map[string]any{
			"key": key,
		},
	}, time.Duration(timeMS)*time.Millisecond)
}

func (e *EmbeddedDevice) Init(env simulation.Environment) {
	e.env = env

	for _, node := range e.radios {
		node.Init(env)
	}

	var ports = make([]device.Port, 0, len(e.ports))
	for _, v := range e.ports {
		ports = append(ports, v)
	}

	e.app.Init(e.schedule, ports...)
}

func (e *EmbeddedDevice) Close() error {
	if cc, ok := e.app.(io.Closer); ok {
		return cc.Close()
	}
	return nil
}

func (e *EmbeddedDevice) Children() []simulation.Node {
	return e.radios
}

func NewEmbeddedDevice(id string, app device.Application, radios []simulation.Node) *EmbeddedDevice {
	d := &EmbeddedDevice{
		id:     id,
		app:    app,
		ports:  make(map[string]bufferNodeWrapper),
		radios: radios,
	}

	for _, r := range d.radios {
		d.ports[r.ID()] = bufferNodeWrapper{
			node: r,
			src:  d,
		}
	}

	return d
}
