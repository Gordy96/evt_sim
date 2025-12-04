package simulation

type Position struct {
	Lat  float64
	Lon  float64
	Elev float64
	Type string
}

type Positionable interface {
	Node
	Position() Position
}
