package configuration

import (
	"errors"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

type Module struct {
	Type    string    `hcl:"type,label"`
	ForEach cty.Value `hcl:"for_each,optional"`
	Rest    hcl.Body  `hcl:",remain"`
}

func (m *Module) Decode(ctx *hcl.EvalContext, l *zap.Logger) ([]simulation.Node, error) {
	switch m.Type {
	case "embedded":
		if m.ForEach.IsNull() || !m.ForEach.CanIterateElements() {
			e, err := m.decodeOne(ctx.NewChild(), l)
			if err != nil {
				return nil, err
			}

			return []simulation.Node{e}, nil
		}

		var res []simulation.Node

		for i, v := range m.ForEach.Elements() {
			subctx := ctx.NewChild()
			subctx.Variables = map[string]cty.Value{
				"each": v,
				"iter": i,
			}
			e, err := m.decodeOne(subctx, l)
			if err != nil {
				return nil, err
			}
			res = append(res, e)
		}

		return res, nil
	}

	return nil, errors.New("unknown module type " + m.Type)
}

func (m *Module) decodeOne(ctx *hcl.EvalContext, l *zap.Logger) (simulation.Node, error) {
	var e embeddedModule
	diags := gohcl.DecodeBody(m.Rest, ctx, &e)
	if diags.HasErrors() {
		return nil, diags
	}

	parsed, err := e.Decode(ctx, l)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}
