package radio

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

func NewRadioMedium(propagationDelay time.Duration, l *zap.Logger) *RadioMedium {
	temp := &RadioMedium{}

	temp.l = l
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

func (r *RadioMedium) OnMessage(msg *simulation.Message) {
	r.l.Info("radio medium, aka air received message", zap.Any("message", msg))

	//here you can handle geo positioning, frequency node state etc
	nodes := r.env.Nodes()

	iprop, ok := r.Param("propagationDelay")
	if !ok {
		panic("radio must have propagationDelay")
	}

	propagationDelay := iprop.(time.Duration)

	srcFreq := getFrequency(nodes[msg.Src])

	for _, node := range nodes {
		_, busy := node.Param("busy")
		if getFrequency(node) == srcFreq && msg.Src != node.ID() && !busy {
			newMsg := *msg
			newMsg.Dst = node.ID()
			newMsg.Src = "radio"
			newMsg.Timestamp = newMsg.Timestamp.Add(propagationDelay)
			r.env.SendMessage(&newMsg)
		}
	}
}

func getFrequency(n simulation.Node) float64 {
	ifreq, ok := n.Param("radioFrequency")
	if !ok {
		return -1
	}
	return ifreq.(float64)
}
