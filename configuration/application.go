package configuration

import (
	"errors"

	"github.com/Gordy96/evt-sim/modules/adapter"
	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type applicationModule struct {
	Type string   `hcl:"type,label"`
	Rest hcl.Body `hcl:",remain"`
}

func (a *applicationModule) Decode(ctx *hcl.EvalContext) (device.Application, error) {
	switch a.Type {
	case "shared":
		var app SharedCApplication
		diags := gohcl.DecodeBody(a.Rest, ctx, &app)
		if diags.HasErrors() {
			return nil, diags
		}

		return app.Decode()
	}

	return nil, errors.New("unknown application type " + a.Type)
}

type SharedCApplication struct {
	Path   string   `hcl:"path"`
	Extras hcl.Body `hcl:",remain"`
}

func (a *SharedCApplication) Decode() (device.Application, error) {
	lib, err := adapter.OpenLib(a.Path)
	if err != nil {
		return nil, err
	}

	params, err := a.finalize()

	if err != nil {
		return nil, err
	}

	return adapter.New(lib, params)
}

func (a *SharedCApplication) finalize() (map[string]interface{}, error) {
	if a.Extras == nil {
		return nil, nil
	}

	// Decode attributes only
	attrs, diags := a.Extras.JustAttributes()
	if diags.HasErrors() {
		return nil, diags
	}

	values := map[string]interface{}{}

	for name, attr := range attrs {
		v, diag := attr.Expr.Value(nil)
		if diag.HasErrors() {
			return nil, diag
		}

		switch v.Type().FriendlyName() {
		case "string":
			values[name] = v.AsString()
		case "number":
			f := v.AsBigFloat()
			if f.IsInt() {
				values[name], _ = f.Int64()
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
