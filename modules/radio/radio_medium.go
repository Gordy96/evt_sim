package radio

import (
	"math"
	"slices"
	"time"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/Gordy96/evt-sim/simulation/message"
	"github.com/tidwall/geodesic"
	"go.uber.org/zap"
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

func NewRadioMedium(l *zap.Logger, options ...Option) *RadioMedium {
	var o mediumOptions

	for _, opt := range options {
		opt(&o)
	}

	temp := &RadioMedium{
		mediumOptions: o,
		l:             l,
	}

	return temp
}

var _ simulation.Node = (*RadioMedium)(nil)

type RadioMedium struct {
	env           simulation.Environment
	radios        []radioNode
	mediumOptions mediumOptions
	l             *zap.Logger
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

func (r *RadioMedium) OnMessage(msg message.Message) {
	//here you can handle geo positioning, frequency node state etc
	src := r.env.FindNode(msg.Src).(radioNode)

	if len(r.radios) == 0 {
		nodes := make([]simulation.Node, 0)
		for _, n := range r.env.Nodes() {
			nodes = append(nodes, n)
		}
		r.cacheRadioNodes(nodes)
	}

	var cache []selectionSetEntry
	for _, node := range r.radios {
		reachable, dist := distanceToOther(src, node)

		if reachable {
			cache = append(cache, selectionSetEntry{
				node: node,
				dist: dist,
			})
		}
	}

	slices.SortFunc(cache, func(a, b selectionSetEntry) int {
		if a.dist > b.dist {
			return -1
		}
		return 1
	})

	for _, c := range cache {
		ttf := timeOfFlightAirNsInt(c.dist)
		parameters := message.Parameters{}

		mb := msg.Builder().
			WithParams(parameters.WithString("origin", msg.Src)).
			WithDst(c.node.ID()).
			WithSrc("radio").
			WithKind("ota/start")
		r.env.SendMessage(mb.Build(), ttf)
	}

}

type selectionSetEntry struct {
	node simulation.Node
	dist float64
}

type radioNode interface {
	simulation.Node
	Frequency() float64
	Power() float64
	Reachable(dist float64) bool
}

func (r *RadioMedium) cacheRadioNodes(nodes []simulation.Node) {
	for _, node := range nodes {
		if f, ok := node.(radioNode); ok {
			r.radios = append(r.radios, f)
		}
	}
}

func distanceToOther(from simulation.Node, to simulation.Node) (bool, float64) {
	if from.ID() == to.ID() {
		return false, -1
	}
	f, ok := from.(radioNode)
	if !ok {
		return false, -1
	}

	t, ok := to.(radioNode)
	if !ok {
		return false, -1
	}

	if !matchingFrequencies(t.Frequency(), f.Frequency(), 0.1) {
		return false, -1
	}

	ps := findPositionableNode(t)
	pd := findPositionableNode(f)

	if ps != nil && pd != nil {
		var p1 = ps.Position()
		var p2 = pd.Position()
		var dist float64
		var zeropos simulation.Position
		if zeropos != p1 && zeropos != p2 {
			geodesic.WGS84.Inverse(p1.Lat, p1.Lon, p2.Lat, p2.Lon, &dist, nil, nil)
			if !t.Reachable(dist) {
				return false, dist
			}

			return true, dist
		}
	}

	return false, 0
}

func matchingFrequencies(a, b float64, threshold float64) bool {
	return math.Abs(a-b) < threshold
}

func findPositionableNode(src simulation.Node) simulation.Positionable {
	if src == nil {
		return nil
	}

	if p, ok := src.(simulation.Positionable); ok {
		return p
	}

	return findPositionableNode(src.Parent())
}

const (
	speedOfLightInt = int64(299_792_458) // m/s
	airCoefPPM      = int64(999_700)     // 0.9997 expressed in ppm
)

func timeOfFlightAirNsInt(distanceMeters float64) time.Duration {
	speedInAir := speedOfLightInt * airCoefPPM / 1_000_000
	t := int64(distanceMeters*1_000_000_000) / speedInAir
	return time.Duration(t) * time.Nanosecond
}
