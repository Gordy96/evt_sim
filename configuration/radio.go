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
)

type radioModule struct {
	Type string   `hcl:"type,label"`
	ID   string   `hcl:"id,optional"`
	Rest hcl.Body `hcl:",remain"`
}

func (r *radioModule) Decode(ctx *hcl.EvalContext, parent simulation.Node) (simulation.Node, error) {
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

		return radio.Decode(r.ID, parent), nil
	}

	return nil, errors.New("unknown radio type " + r.Type)
}

type LoRaNIC struct {
	Frequency     float64       `hcl:"frequency"`
	Power         uint64        `hcl:"power"`
	ReceiveDelay  time.Duration `hcl:"receive_delay,optional"`
	TransmitDelay time.Duration `hcl:"transmit_delay,optional"`
}

func (l *LoRaNIC) Decode(id string, parent simulation.Node) *lora.LoraNic {
	return lora.New(
		filepath.Join(parent.ID(), id),
		l.Frequency,
		lora.WithPower(l.Power),
		lora.WithReceiveDelay(l.ReceiveDelay),
		lora.WithTransmitDelay(l.TransmitDelay),
		lora.WithParent(parent),
	)
}
