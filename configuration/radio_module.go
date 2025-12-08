package configuration

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/Gordy96/evt-sim/modules/device/lora"
	"github.com/Gordy96/evt-sim/simulation"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

type radioModule struct {
	Type string   `hcl:"type,label"`
	ID   string   `hcl:"id,optional"`
	Rest hcl.Body `hcl:",remain"`
}

func (r *radioModule) Decode(ctx *hcl.EvalContext, parent simulation.Node, l *zap.Logger) (simulation.Node, error) {
	if r.ID == "" {
		r.ID = "radio"
	}

	bubbleUp(ctx, "needs_radio", cty.True)

	switch r.Type {
	case "lora":
		var radio LoRaNIC
		diags := gohcl.DecodeBody(r.Rest, ctx, &radio)
		if diags.HasErrors() {
			return nil, diags
		}

		return radio.Decode(r.ID, parent, l), nil
	}

	return nil, errors.New("unknown radio type " + r.Type)
}

type LoRaNIC struct {
	FrequencyHZ   float64       `hcl:"frequency_hz"`
	Power         float64       `hcl:"power,optional"`
	FadeMargin    float64       `hcl:"fade_margin,optional"`
	ReceiveDelay  time.Duration `hcl:"receive_delay,optional"`
	TransmitDelay time.Duration `hcl:"transmit_delay,optional"`
}

func (l *LoRaNIC) Decode(id string, parent simulation.Node, logger *zap.Logger) *lora.LoraNic {
	opts := []lora.Option{
		lora.WithPower(l.Power),
		lora.WithReceiveDelay(l.ReceiveDelay),
		lora.WithTransmitDelay(l.TransmitDelay),
		lora.WithParent(parent),
	}

	if l.FadeMargin > 0 {
		opts = append(opts, lora.WithFadeMargin(l.FadeMargin))
	}

	return lora.New(
		filepath.Join(parent.ID(), id),
		l.FrequencyHZ,
		opts...,
	)
}
