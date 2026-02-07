package configuration

import (
	"github.com/Gordy96/evt-sim/modules/radio"
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"go.uber.org/zap"
)

type Config struct {
	Name        string       `hcl:"name"`
	Modules     []Module     `hcl:"module,block"`
	RadioMedium *radioMedium `hcl:"radio_medium,block"`
}

func (c *Config) Decode(ctx *hcl.EvalContext, l *zap.Logger) ([]simulation.Node, error) {
	var res []simulation.Node
	for _, module := range c.Modules {
		m, err := module.Decode(ctx.NewChild(), l)
		if err != nil {
			return nil, err
		}
		res = append(res, m...)
	}

	//TODO: add medium
	if v, ok := ctx.Variables["needs_radio"]; ok && v.True() || c.RadioMedium != nil {
		var opts []radio.Option

		if c.RadioMedium != nil && c.RadioMedium.BackgroundNoiseLevel != nil {
			opts = append(opts, radio.WithBackgroundNoiseLevel(*c.RadioMedium.BackgroundNoiseLevel))
		}

		res = append(res, radio.NewRadioMedium(l, opts...))
	}

	return res, nil
}

type radioMedium struct {
	BackgroundNoiseLevel *float64 `hcl:"background_noise_level,optional"`
}
