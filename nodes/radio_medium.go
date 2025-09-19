package nodes

import (
	"encoding/json"
	"time"
)

func NewRadioMedium(propagationDelay time.Duration) *RadioMedium {
	return &RadioMedium{
		params: map[string]any{
			"propagationDelay": propagationDelay,
		},
	}
}

type RadioMedium struct {
	sim    *Simulation
	params map[string]any
}

func (r *RadioMedium) ID() string {
	return "radio"
}

func (r *RadioMedium) Param(name string) (any, bool) {
	if r.params == nil {
		return nil, false
	}

	ret, ok := r.params[name]
	return ret, ok
}

func (r *RadioMedium) SetParam(name string, value any) {
	if r.params == nil {
		r.params = make(map[string]any)
	}

	r.params[name] = value
}

func (r *RadioMedium) Init(sim *Simulation) {
	r.sim = sim
}

func (r *RadioMedium) HandleMessage(msg *Message, sim *Simulation, timestamp time.Time) {
	j, _ := json.Marshal(msg)
	sim.Log("radio medium, aka air received message: %s", j)

	//here you can handle geo positioning, frequency node state etc
	nodes := sim.Nodes()

	iprop, ok := r.Param("propagationDelay")
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
			newMsg.Timestamp = newMsg.Timestamp.Add(propagationDelay)
			sim.SendMessage(&newMsg)
		}
	}
}

func getFrequency(n Node) float64 {
	ifreq, ok := n.Param("radioFrequency")
	if !ok {
		return -1
	}
	return ifreq.(float64)
}
