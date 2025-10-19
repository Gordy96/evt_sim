package lora

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
)

type LoraNic struct {
	simulation.ParameterBag
	id       string
	env      simulation.Environment
	parentID string
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
	case "start_receiving":
		//TODO: reject when already receiving and/or calculate SNR to drop messages as noise
		irec, ok := l.GetParam("receiving")
		if !ok || !irec.(bool) {
			l.SetParam("receiving", true)
			l.sendSelf(simulation.Message{
				Kind:   "finish_receiving",
				Params: msg.Params,
			}, time.Duration(10)*time.Nanosecond)
		}
	case "finish_receiving":
		l.RemoveParam("receiving")
		l.env.SendMessage(&simulation.Message{
			ID:     "",
			Src:    l.ID(),
			Dst:    l.parentID,
			Kind:   "received_radio_message",
			Params: msg.Params,
		}, 0)
	case "start_sending":
		//TODO: reject when already sending
		isnd, ok := l.GetParam("sending")
		if !ok || !isnd.(bool) {
			l.SetParam("sending", true)
			l.sendSelf(simulation.Message{
				Kind:   "finish_sending",
				Params: msg.Params,
			}, time.Duration(10)*time.Nanosecond)
		}
	case "finish_sending":
		l.RemoveParam("sending")
		l.env.SendMessage(&simulation.Message{
			ID:     "",
			Src:    l.id,
			Dst:    "radio",
			Kind:   "radio_message",
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
