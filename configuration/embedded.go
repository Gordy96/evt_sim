package configuration

import (
	"github.com/Gordy96/evt-sim/modules/device/embedded"
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"go.uber.org/zap"
)

type embeddedModule struct {
	Radios      []radioModule     `hcl:"radio,block"`
	Application applicationModule `hcl:"application,block"`
	Position    position          `hcl:"position,block"`
}

func (e *embeddedModule) Decode(ctx *hcl.EvalContext, id string, l *zap.Logger) (simulation.Node, error) {
	app, err := e.Application.Decode(ctx.NewChild(), l.Named(id))
	if err != nil {
		return nil, err
	}

	var ops []embedded.DeviceOption

	var zeropos position

	if zeropos != e.Position {
		ops = append(ops, embedded.WithPosition(simulation.Position{
			Type: e.Position.Type,
			Lat:  e.Position.Lat,
			Lon:  e.Position.Lon,
			Elev: e.Position.Elev,
		}))
	}

	dev := embedded.New(id, app, ops...)

	for _, radio := range e.Radios {
		r, err := radio.Decode(ctx.NewChild(), dev, l)
		if err != nil {
			return nil, err
		}

		dev.AddConnection(r, embedded.WithName(radio.ID))
	}

	return dev, nil
}
