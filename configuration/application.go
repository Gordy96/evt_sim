package configuration

import (
	"errors"
	"math"
	"plugin"
	"strings"

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
	case "goplugin":
		var app GoPluginApplication
		diags := gohcl.DecodeBody(a.Rest, ctx, &app)
		if diags.HasErrors() {
			return nil, diags
		}

		return app.Decode(ctx.NewChild(), l)
	}

	return nil, errors.New("unknown application type " + a.Type)
}

type GoPluginApplication struct {
	Path       string    `hcl:"path"`
	Parameters cty.Value `hcl:"parameters,optional"`
}

func (a *GoPluginApplication) Decode(ctx *hcl.EvalContext, l *zap.Logger) (device.Application, error) {
	path := strings.Split(a.Path, "#")
	p, err := plugin.Open(path[0])
	if err != nil {
		return nil, err
	}

	cstr, err := p.Lookup(path[1])
	if err != nil {
		return nil, err
	}

	constr := cstr.(func(l *zap.Logger, params map[string]any) device.Application)

	params, err := finalize(ctx, a.Parameters)

	if err != nil {
		return nil, err
	}

	appLogger := l.Named("application")

	return constr(appLogger, params), nil
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

	params, err := finalize(ctx, a.Parameters)

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

func finalize(ctx *hcl.EvalContext, parameters cty.Value) (map[string]interface{}, error) {
	if parameters.IsNull() || !parameters.CanIterateElements() {
		return nil, nil
	}

	values := map[string]interface{}{}

	for name, v := range parameters.AsValueMap() {
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
