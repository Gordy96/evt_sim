package lora

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/Gordy96/evt-sim/simulation/message"
)

func New(id string, frequencyHz float64, options ...Option) *LoraNic {
	var o = loraOptions{
		fHz:           frequencyHz,
		Pt:            20,
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
	return l.options.fHz
}

func (l *LoraNic) Power() float64 {
	return l.options.Pt
}

func (l *LoraNic) ID() string {
	return l.id
}

func (l *LoraNic) sendSelf(kind message.Kind, params message.Parameters, delay time.Duration) {
	l.env.SendMessage(message.Builder{}.
		WithDst(l.ID()).
		WithSrc(l.ID()).
		WithKind(kind).
		WithParams(params).
		Build(),
		delay,
	)
}

func (l *LoraNic) OnMessage(msg message.Message) {
	switch msg.Kind {
	case "ota/start":
		l.sendSelf("ota/finish", msg.Params, l.options.receiveDelay)
	case "ota/finish":
		l.env.SendMessage(message.Builder{}.
			WithSrc(l.ID()).
			WithDst(l.options.parent.ID()).
			WithKind("interrupt/port").
			WithParams(msg.Params).
			Build(), 0)
	case "wire/payload":
		l.sendSelf("wire/finish", msg.Params, l.options.transmitDelay)
	case "wire/finish":
		l.env.SendMessage(message.Builder{}.
			WithSrc(l.ID()).
			WithDst("radio").
			WithKind("radio/message").
			WithParams(msg.Params).
			Build(), 0)
	}
}

func (l *LoraNic) Init(env simulation.Environment) {
	l.env = env
}

func (l *LoraNic) Close() error {
	return nil
}

var _ simulation.Node = (*LoraNic)(nil)

func (l *LoraNic) Reachable(dist float64) bool {
	return dist < l.options.fadeMargin
}
