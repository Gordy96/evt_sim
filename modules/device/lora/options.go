package lora

import (
	"time"

	"github.com/Gordy96/evt-sim/simulation"
)

type Option func(*loraOptions)

type loraOptions struct {
	Pt         float64
	Gt         float64
	Gr         float64
	Lsys       float64
	fHz        float64
	BW         float64
	NF         float64
	SF         uint64
	SNR        map[uint64]float64
	FadeMargin float64

	receiveDelay  time.Duration
	transmitDelay time.Duration
	parent        simulation.Node
}

func WithPower(power float64) Option {
	return func(o *loraOptions) {
		o.Pt = power
	}
}

func WithGains(Gt, Gr float64) Option {
	return func(o *loraOptions) {
		o.Gt = Gt
		o.Gr = Gr
	}
}

func WithSystemLoss(lsys float64) Option {
	return func(o *loraOptions) {
		o.Lsys = lsys
	}
}

func WithFrequencyHz(fHz float64) Option {
	return func(o *loraOptions) {
		o.fHz = fHz
	}
}

func WithNoiseFigure(nf float64) Option {
	return func(o *loraOptions) {
		o.NF = nf
	}
}

func WithBW(BW float64) Option {
	return func(o *loraOptions) {
		o.BW = BW
	}
}

func WithSF(SF uint64) Option {
	return func(o *loraOptions) {
		o.SF = SF
	}
}

func WithNSR(SF uint64, SNR float64) Option {
	return func(o *loraOptions) {
		if o.SNR == nil {
			o.SNR = make(map[uint64]float64)
		}
		o.SNR[SF] = SNR
	}
}

func WithFadeMargin(fm float64) Option {
	return func(o *loraOptions) {
		o.FadeMargin = fm
	}
}

func WithReceiveDelay(receiveDelay time.Duration) Option {
	return func(o *loraOptions) {
		o.receiveDelay = receiveDelay
	}
}

func WithTransmitDelay(transmitDelay time.Duration) Option {
	return func(o *loraOptions) {
		o.transmitDelay = transmitDelay
	}
}

func WithParent(parent simulation.Node) Option {
	return func(o *loraOptions) {
		o.parent = parent
	}
}
