package configuration

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name    string         `hcl:"name,label"`
	Type    hcl.Expression `hcl:"type"`
	Default cty.Value      `hcl:"default"`
}

func ParseVariables(body hcl.Body, ctx *hcl.EvalContext) (hcl.Body, hcl.Diagnostics) {
	decls, rest, diag := body.PartialContent(&hcl.BodySchema{
		Attributes: nil,
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "variable",
				LabelNames: []string{
					"name",
				},
			},
		},
	})

	if diag.HasErrors() {
		return nil, diag
	}

	for _, d := range decls.Blocks {
		var v Variable
		v.Name = d.Labels[0]
		diag = gohcl.DecodeBody(d.Body, ctx, &v)
		if diag.HasErrors() {
			return nil, diag
		}
		ctx.Variables[v.Name] = v.Default
	}

	return rest, nil
}
