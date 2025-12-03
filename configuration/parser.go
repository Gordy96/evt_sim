package configuration

import (
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

func ParseFile(path string) ([]simulation.Node, error) {
	parser := hclparse.NewParser()
	file, diag := parser.ParseHCLFile(path)
	if diag.HasErrors() {
		return nil, diag
	}

	ctx := &hcl.EvalContext{}

	var root Config

	diag = gohcl.DecodeBody(file.Body, ctx, &root)
	if diag.HasErrors() {
		return nil, diag
	}

	return root.Decode(ctx)
}

func bubbleUp(ctx *hcl.EvalContext, name string, val cty.Value) {
	if ctx == nil {
		return
	}

	ctx.Variables[name] = val

	bubbleUp(ctx.Parent(), name, val)
}

func ctxGet(ctx *hcl.EvalContext, name string) (cty.Value, bool) {
	v, ok := ctx.Variables[name]
	if ok {
		return v, true
	}

	if ctx.Parent() != nil {
		return ctxGet(ctx.Parent(), name)
	}

	return cty.NilVal, false
}
