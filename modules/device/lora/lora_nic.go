package lora

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
)

func New(id string, parentID string, options Options) *LoraNic {
	return &LoraNic{
		id:       id,
		parentID: parentID,
		options:  options,
	}
}

type Options struct {
	Frequency     float64
	Power         uint64
	ReceiveDelay  time.Duration
	TransmitDelay time.Duration
}

type state struct {
	receiving    bool
	transmitting bool
}

type LoraNic struct {
	id       string
	env      simulation.Environment
	parentID string
	options  Options
	state    state
}

func (l *LoraNic) Frequency() float64 {
	return l.options.Frequency
}

func (l *LoraNic) Power() uint64 {
	return l.options.Power
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
			}, l.options.ReceiveDelay)
		}
	case "ota/finish":
		l.state.receiving = false
		l.env.SendMessage(&simulation.Message{
			ID:     "",
			Src:    l.ID(),
			Dst:    l.parentID,
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
			}, l.options.TransmitDelay)
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
