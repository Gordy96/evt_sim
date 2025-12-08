package configuration

import (
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

func ParseFile(path string, l *zap.Logger) ([]simulation.Node, error) {
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

	return root.Decode(ctx, l)
}

func topContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	if ctx == nil {
		return nil
	}

	if ctx.Parent() == nil {
		return ctx
	}

	return topContext(ctx.Parent())
}

func bubbleUp(ctx *hcl.EvalContext, name string, val cty.Value) {
	if ctx == nil {
		return
	}

	ctx = topContext(ctx)

	if ctx.Variables == nil {
		ctx.Variables = make(map[string]cty.Value)
	}

	ctx.Variables[name] = val
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
