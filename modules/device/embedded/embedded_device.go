package embedded

import (
	"time"

	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/multierr"
)

type namedConnection struct {
	name string
	dst  simulation.Node
}

type deviceOptions struct {
	position    simulation.Position
	connections []namedConnection
}

type DeviceOption func(*deviceOptions)

func WithPosition(position simulation.Position) DeviceOption {
	return func(o *deviceOptions) {
		o.position = position
	}
}

func WithConnection(name string, dst simulation.Node) DeviceOption {
	return func(o *deviceOptions) {
		o.connections = append(o.connections, namedConnection{name: name, dst: dst})
	}
}

func New(id string, app device.Application, options ...DeviceOption) *EmbeddedDevice {
	var o deviceOptions
	for _, option := range options {
		option(&o)
	}

	d := &EmbeddedDevice{
		id:         id,
		app:        app,
		ports:      make(map[string]*bufferNodeWrapper),
		radios:     make([]simulation.Node, 0),
		portLookup: make(map[string]string),
		options:    o,
	}

	for _, conn := range o.connections {
		d.AddConnection(conn.dst, WithName(conn.name))
	}

	return d
}

var _ simulation.Node = (*EmbeddedDevice)(nil)
var _ simulation.CompositeNode = (*EmbeddedDevice)(nil)
var _ simulation.Positionable = (*EmbeddedDevice)(nil)

type EmbeddedDevice struct {
	id         string
	env        simulation.Environment
	app        device.Application
	ports      map[string]*bufferNodeWrapper
	radios     []simulation.Node
	portLookup map[string]string
	options    deviceOptions
}

func (e *EmbeddedDevice) Parent() simulation.Node {
	return nil
}

func (e *EmbeddedDevice) Position() simulation.Position {
	return e.options.position
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

func (e *EmbeddedDevice) AddConnection(node simulation.Node, options ...ConnectionOption) {
	var o connectionOptions
	for _, option := range options {
		option(&o)
	}

	var name string

	if o.name != "" {
		name = o.name
	} else {
		name = node.ID()
	}

	e.radios = append(e.radios, node)
	e.portLookup[node.ID()] = name
	e.ports[name] = &bufferNodeWrapper{
		name: name,
		node: node,
		src:  e,
	}

	type related interface {
		SetParent(simulation.Node)
		Parent() simulation.Node
	}

	if casted, ok := node.(related); ok {
		casted.SetParent(e)
	}
}

type connectionOptions struct {
	name string
}
type ConnectionOption func(options *connectionOptions)

func WithName(name string) ConnectionOption {
	return func(options *connectionOptions) {
		options.name = name
	}
}
