package configuration

import (
	"errors"
	"math"

	"github.com/Gordy96/evt-sim/modules/adapter"
	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type applicationModule struct {
	Type string   `hcl:"type,label"`
	Rest hcl.Body `hcl:",remain"`
}

func (a *applicationModule) Decode(ctx *hcl.EvalContext, l *zap.Logger) (device.Application, error) {
	switch a.Type {
	case "shared":
		var app SharedCApplication
		diags := gohcl.DecodeBody(a.Rest, ctx, &app)
		if diags.HasErrors() {
			return nil, diags
		}

		return app.Decode(ctx.NewChild(), l)
	}

	return nil, errors.New("unknown application type " + a.Type)
}

type SharedCApplication struct {
	Path        string    `hcl:"path"`
	Concurrent  bool      `hcl:"concurrent,optional"`
	DumpPackets bool      `hcl:"dump_packets,optional"`
	Parameters  cty.Value `hcl:"parameters,optional"`
}

func (a *SharedCApplication) Decode(ctx *hcl.EvalContext, l *zap.Logger) (device.Application, error) {
	lib, err := adapter.OpenLib(a.Path)
	if err != nil {
		return nil, err
	}

	params, err := a.finalize(ctx)

	if err != nil {
		return nil, err
	}

	appLogger := l.Named("application")

	return adapter.New(
		lib,
		adapter.WithParams(params),
		adapter.WithLogger(func(level int, line string) {
			appLogger.Log(zapcore.Level(level), line)
		}),
		adapter.WithConcurrency(a.Concurrent),
		adapter.WithDumpPackets(a.DumpPackets),
	)
}

func (a *SharedCApplication) finalize(ctx *hcl.EvalContext) (map[string]interface{}, error) {
	if a.Parameters.IsNull() || !a.Parameters.CanIterateElements() {
		return nil, nil
	}

	values := map[string]interface{}{}

	for name, v := range a.Parameters.AsValueMap() {
		switch v.Type().FriendlyName() {
		case "string":
			values[name] = v.AsString()
		case "number":
			f := v.AsBigFloat()
			if f.IsInt() {
				i, _ := f.Int64()
				if i > math.MaxInt64 {
					values[name], _ = f.Uint64()
				} else {
					values[name] = i
				}
			} else {
				values[name], _ = f.Float64()
			}
		case "bool":
			values[name] = v.True()
		default:
		}
	}

	return values, nil
}
