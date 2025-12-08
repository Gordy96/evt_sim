package configuration

import (
	"errors"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"go.uber.org/zap"
)

type Module struct {
	Type string   `hcl:"type,label"`
	ID   string   `hcl:"id,label"`
	Rest hcl.Body `hcl:",remain"`
}

func (m *Module) Decode(ctx *hcl.EvalContext, l *zap.Logger) (simulation.Node, error) {
	var subctx = ctx.NewChild()
	switch m.Type {
	case "embedded":
		var e embeddedModule
		diags := gohcl.DecodeBody(m.Rest, subctx, &e)
		if diags.HasErrors() {
			return nil, diags
		}

		return e.Decode(subctx, m.ID, l)
	}

	return nil, errors.New("unknown module type " + m.Type)
}
