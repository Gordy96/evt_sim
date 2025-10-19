package radio

import (
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
	l   *zap.Logger
	env simulation.Environment
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
	r.l.Info("radio medium, aka air received message", zap.Any("message", msg))

	//here you can handle geo positioning, frequency node state etc
	nodes := r.env.Nodes()

	iprop, ok := r.GetParam("propagationDelay")
	if !ok {
		panic("radio must have propagationDelay")
	}

	propagationDelay := iprop.(time.Duration)

	srcFreq := getFrequency(nodes[msg.Src])

	for _, node := range nodes {
		if getFrequency(node) == srcFreq && msg.Src != node.ID() {
			newMsg := *msg
			newMsg.Dst = node.ID()
			newMsg.Src = "radio"
			newMsg.Kind = "start_receiving"
			r.env.SendMessage(&newMsg, propagationDelay)
		}
	}
}

func getFrequency(n simulation.Node) float64 {
	ifreq, ok := n.GetParam("radioFrequency")
	if !ok {
		return -1
	}
	return ifreq.(float64)
}
