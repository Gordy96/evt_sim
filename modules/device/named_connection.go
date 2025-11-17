package device

import "github.com/Gordy96/evt-sim/simulation"

type NamedConnection struct {
	Name string
	Dst  simulation.Node
}
