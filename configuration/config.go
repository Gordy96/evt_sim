package configuration

import (
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
)

type Config struct {
	Name    string   `hcl:"name"`
	Modules []Module `hcl:"module,block"`
}

func (c *Config) Decode(ctx *hcl.EvalContext) ([]simulation.Node, error) {
	res := make([]simulation.Node, len(c.Modules))
	for i, module := range c.Modules {
		m, err := module.Decode(ctx.NewChild())
		if err != nil {
			return nil, err
		}
		res[i] = m
	}

	if v, ok := ctx.Variables["needs_radio"]; ok && v.True() {
		res = append(res)
	}

	return res, nil
}
