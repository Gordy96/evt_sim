package radio

import (
	"math"

	"github.com/Gordy96/evt-sim/simulation"
)

type mediumOptions struct {
	backgroundNoiseLevel float64
}

type Option func(*mediumOptions)

func WithBackgroundNoiseLevel(l float64) Option {
	return func(o *mediumOptions) {
		o.backgroundNoiseLevel = l
	}
}

func NewRadioMedium(options ...Option) *RadioMedium {
	var o mediumOptions

	for _, opt := range options {
		opt(&o)
	}

	temp := &RadioMedium{
		mediumOptions: o,
	}

	return temp
}

var _ simulation.Node = (*RadioMedium)(nil)

type RadioMedium struct {
	env           simulation.Environment
	radios        []radioNode
	mediumOptions mediumOptions
}

func (r *RadioMedium) Parent() simulation.Node {
	return nil
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
	//here you can handle geo positioning, frequency node state etc
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
			newMsg.Kind = "ota/start"
			r.env.SendMessage(&newMsg, 0)
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
