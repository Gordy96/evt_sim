package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"go.uber.org/zap"
)

type Simulation struct {
	Nodes    []simulation.Node
	Realtime bool
	Name     string
}

func ParseFile(path string, l *zap.Logger) (*Simulation, error) {
	parser := hclparse.NewParser()
	file, diag := parser.ParseHCLFile(path)
	if diag.HasErrors() {
		return nil, diag
	}

	env := make(map[string]cty.Value)

	for _, kv := range os.Environ() {
		if eq := strings.IndexByte(kv, '='); eq != -1 {
			key := kv[:eq]
			val := kv[eq+1:]
			env[key] = cty.StringVal(val)
		}
	}

	ctx := &hcl.EvalContext{
		Variables: env,
		Functions: map[string]function.Function{
			"itoa": function.New(&function.Spec{
				Description:  "converts int to string",
				Params:       []function.Parameter{{Name: "i", Type: cty.Number}},
				VarParam:     nil,
				Type:         function.StaticReturnType(cty.String),
				RefineResult: nil,
				Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
					i, _ := args[0].AsBigFloat().Int64()
					return cty.StringVal(fmt.Sprintf("%d", i)), nil
				},
			}),
		},
	}

	rest, diag := ParseVariables(file.Body, ctx)

	var root Config

	diag = gohcl.DecodeBody(rest, ctx, &root)
	if diag.HasErrors() {
		return nil, diag
	}

	nodes, err := root.Decode(ctx, l)
	if err != nil {
		return nil, err
	}

	return &Simulation{
		Nodes:    nodes,
		Realtime: root.Realtime,
		Name:     root.Name,
	}, nil
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
