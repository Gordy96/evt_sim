package radio

import (
	"math"
	"time"

	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

func NewRadioMedium(l *zap.Logger, propagationDelay time.Duration) *RadioMedium {
	temp := &RadioMedium{
		l: l.Named("radio_medium"),
	}

	temp.SetParam("propagationDelay", propagationDelay)

	return temp
}

var _ simulation.Node = (*RadioMedium)(nil)

type RadioMedium struct {
	simulation.ParameterBag
	l      *zap.Logger
	env    simulation.Environment
	radios []radioNode
}

func (r *RadioMedium) ID() string {
	return "radio"
}

func (r *RadioMedium) Init(env simulation.Environment) {
	r.env = env
}

func (r *RadioMedium) Close() error {
	return nil
}

func (r *RadioMedium) OnMessage(msg *simulation.Message) {
	r.l.Debug("radio medium, aka air received message", zap.Any("message", msg))

	//here you can handle geo positioning, frequency node state etc
	iprop, ok := r.GetParam("propagationDelay")
	if !ok {
		panic("radio must have propagationDelay")
	}

	propagationDelay := iprop.(time.Duration)

	srcFreq := r.env.FindNode(msg.Src).(radioNode).Frequency()

	if len(r.radios) == 0 {
		nodes := make([]simulation.Node, 0)
		for _, n := range r.env.Nodes() {
			nodes = append(nodes, n)
		}
		r.cacheRadioNodes(nodes)
	}

	for _, node := range r.radios {
		if matchingFrequencies(srcFreq, node.Frequency(), 0.1) && msg.Src != node.ID() {
			newMsg := *msg
			newMsg.Dst = node.ID()
			newMsg.Src = "radio"
			newMsg.Kind = "start_receiving"
			r.env.SendMessage(&newMsg, propagationDelay)
		}
	}
}

func (r *RadioMedium) cacheRadioNodes(nodes []simulation.Node) {
	for _, node := range nodes {
		if f, ok := node.(radioNode); ok {
			r.radios = append(r.radios, f)
		}
	}
}

type radioNode interface {
	simulation.Node
	Frequency() float64
	Power() uint64
}

func matchingFrequencies(a, b float64, threshold float64) bool {
	return math.Abs(a-b) < threshold
}
