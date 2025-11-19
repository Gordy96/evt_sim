package embedded

import (
	"time"

	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/multierr"
)

var _ simulation.Node = (*EmbeddedDevice)(nil)

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
	b.src.env.SendMessage(&simulation.Message{
		ID:   "",
		Src:  b.src.ID(),
		Dst:  b.node.ID(),
		Kind: "wire/payload",
		Params: map[string]any{
			"payload": c,
		},
	}, 0)

	return len(buf), nil
}

type EmbeddedDevice struct {
	id         string
	env        simulation.Environment
	app        device.Application
	ports      map[string]*bufferNodeWrapper
	radios     []simulation.Node
	portLookup map[string]string
}

func (e *EmbeddedDevice) ID() string {
	return e.id
}

func (e *EmbeddedDevice) OnMessage(msg *simulation.Message) {
	switch msg.Kind {
	case "interrupt/delay":
		key := msg.Params["key"].(string)
		e.app.TriggerTimeInterrupt(key)
	case "interrupt/port":
		if portName, ok := e.portLookup[msg.Src]; ok {
			port := e.ports[portName]
			ipl, ok := msg.Params["payload"]
			if ok {
				payload := ipl.([]byte)
				port.buf = append(port.buf[:], payload...)
				e.app.TriggerPortInterrupt(portName)
			}
		}
	}
}

func (e *EmbeddedDevice) schedule(key string, timeMS int) {
	e.env.SendMessage(&simulation.Message{
		ID:   "",
		Src:  e.ID(),
		Dst:  e.ID(),
		Kind: "interrupt/delay",
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
	var err = e.app.Close()
	for _, node := range e.radios {
		err = multierr.Append(err, node.Close())
	}
	return err
}

func (e *EmbeddedDevice) Children() []simulation.Node {
	return e.radios
}

func New(id string, app device.Application, radios ...device.NamedConnection) *EmbeddedDevice {
	d := &EmbeddedDevice{
		id:         id,
		app:        app,
		ports:      make(map[string]*bufferNodeWrapper),
		radios:     make([]simulation.Node, 0, len(radios)),
		portLookup: make(map[string]string),
	}

	for _, r := range radios {
		d.radios = append(d.radios, r.Dst)
		d.portLookup[r.Dst.ID()] = r.Name
		d.ports[r.Name] = &bufferNodeWrapper{
			name: r.Name,
			node: r.Dst,
			src:  d,
		}
	}

	return d
}
