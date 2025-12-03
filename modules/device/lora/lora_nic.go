package lora

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
)

type Option func(*loraOptions)

type loraOptions struct {
	frequency     float64
	power         uint64
	receiveDelay  time.Duration
	transmitDelay time.Duration
	parent        simulation.Node
}

func WithPower(power uint64) Option {
	return func(o *loraOptions) {
		o.power = power
	}
}

func WithReceiveDelay(receiveDelay time.Duration) Option {
	return func(o *loraOptions) {
		o.receiveDelay = receiveDelay
	}
}

func WithTransmitDelay(transmitDelay time.Duration) Option {
	return func(o *loraOptions) {
		o.transmitDelay = transmitDelay
	}
}

func WithParent(parent simulation.Node) Option {
	return func(o *loraOptions) {
		o.parent = parent
	}
}

func New(id string, frequency float64, options ...Option) *LoraNic {
	var o = loraOptions{
		frequency:     frequency,
		power:         20,
		receiveDelay:  10 * time.Millisecond,
		transmitDelay: 10 * time.Millisecond,
		parent:        nil,
	}

	for _, option := range options {
		option(&o)
	}

	return &LoraNic{
		id:      id,
		options: o,
	}
}

type state struct {
	receiving    bool
	transmitting bool
}

type LoraNic struct {
	id      string
	env     simulation.Environment
	options loraOptions
	state   state
}

func (l *LoraNic) SetParent(parent simulation.Node) {
	l.options.parent = parent
}

func (l *LoraNic) Parent() simulation.Node {
	return l.options.parent
}

func (l *LoraNic) Frequency() float64 {
	return l.options.frequency
}

func (l *LoraNic) Power() uint64 {
	return l.options.power
}

func (l *LoraNic) ID() string {
	return l.id
}

func (l *LoraNic) sendSelf(msg simulation.Message, delay time.Duration) {
	l.env.SendMessage(&simulation.Message{
		ID:     msg.ID,
		Src:    l.ID(),
		Dst:    l.ID(),
		Kind:   msg.Kind,
		Params: msg.Params,
	}, delay)
}

func (l *LoraNic) OnMessage(msg *simulation.Message) {
	switch msg.Kind {
	case "ota/start":
		//TODO: reject when already receiving and/or calculate SNR to drop messages as noise
		if !l.state.receiving {
			l.state.receiving = true
			l.sendSelf(simulation.Message{
				Kind:   "ota/finish",
				Params: msg.Params,
			}, l.options.receiveDelay)
		}
	case "ota/finish":
		l.state.receiving = false
		l.env.SendMessage(&simulation.Message{
			ID:     "",
			Src:    l.ID(),
			Dst:    l.options.parent.ID(),
			Kind:   "interrupt/port",
			Params: msg.Params,
		}, 0)
	case "wire/payload":
		//TODO: reject when already sending
		if !l.state.transmitting {
			l.state.transmitting = true
			l.sendSelf(simulation.Message{
				Kind:   "wire/finish",
				Params: msg.Params,
			}, l.options.transmitDelay)
		}
	case "wire/finish":
		l.state.transmitting = false
		l.env.SendMessage(&simulation.Message{
			ID:     "",
			Src:    l.id,
			Dst:    "radio",
			Kind:   "radio/message",
			Params: msg.Params,
		}, 0)
	}
}

func (l *LoraNic) Init(env simulation.Environment) {
	l.env = env
}

func (l *LoraNic) Close() error {
	return nil
}

var _ simulation.Node = (*LoraNic)(nil)
