package configuration

import (
	"errors"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type Module struct {
	Type string   `hcl:"type,label"`
	ID   string   `hcl:"id,label"`
	Rest hcl.Body `hcl:",remain"`
}

func (m *Module) Decode(ctx *hcl.EvalContext) (simulation.Node, error) {
	switch m.Type {
	case "embedded":
		var e embeddedModule
		diags := gohcl.DecodeBody(m.Rest, ctx.NewChild(), &e)
		if diags.HasErrors() {
			return nil, diags
		}

		return e.Decode(ctx, m.ID)
	}

	return nil, errors.New("unknown module type " + m.Type)
}
