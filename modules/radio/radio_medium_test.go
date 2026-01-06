package radio

import (
	"testing"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/Gordy96/evt-sim/simulation/message"
	"github.com/stretchr/testify/assert"
)

type fakeNode simulation.Position

func (f fakeNode) ID() string { return f.Type }

func (f fakeNode) OnMessage(msg message.Message) { panic("implement me") }

func (f fakeNode) Init(env simulation.Environment) { panic("implement me") }

func (f fakeNode) Close() error { panic("implement me") }

func (f fakeNode) Parent() simulation.Node { return nil }

func (f fakeNode) Frequency() float64 { return 0.0 }

func (f fakeNode) Power() float64 { panic("implement me") }

func (f fakeNode) Reachable(dist float64) bool { return true }

func (f fakeNode) Position() simulation.Position { return simulation.Position(f) }

var _ radioNode = (*fakeNode)(nil)

func TestDistanceToOther(t *testing.T) {
	var first, second, third = fakeNode{Lat: 50.45624, Lon: 30.36545, Type: "1"}, fakeNode{Lat: 50.45422, Lon: 30.44862, Type: "2"}, fakeNode{Lat: 50.44812, Lon: 30.525}

	_, fts := distanceToOther(first, second)
	assert.InDelta(t, 5910, fts, 10)

	_, stt := distanceToOther(second, third)
	assert.InDelta(t, 5460, stt, 10)

	_, ftt := distanceToOther(first, third)
	assert.InDelta(t, 11360, ftt, 10)
}
