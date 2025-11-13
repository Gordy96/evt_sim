package main

import (
	"time"

	"github.com/Gordy96/evt-sim/modules/device"
	"github.com/Gordy96/evt-sim/modules/device/embedded"
	"github.com/Gordy96/evt-sim/modules/device/lora"
	"github.com/Gordy96/evt-sim/modules/radio"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

var _ device.Application = (*FakeApp)(nil)

type FakeApp struct {
	ports       map[string]device.Port
	scheduler   func(key string, timeMS int)
	initializer bool
	l           *zap.Logger
	port        string
}

func (f *FakeApp) Init(scheduler func(key string, timeMS int), ports ...device.Port) error {
	f.scheduler = scheduler

	if f.ports == nil {
		f.ports = make(map[string]device.Port)
	}

	for _, port := range ports {
		f.ports[port.Name()] = port
	}

	if f.initializer {
		f.ports[f.port].Write([]byte("ping"))
	}

	return nil
}

func (f *FakeApp) TriggerPortInterrupt(port string) error {
	var buf [128]byte
	n, err := f.ports[port].Read(buf[:])
	if err != nil {
		return err
	}

	f.l.Info("Received", zap.String("port", port), zap.ByteString("data", buf[:n]))

	if !f.initializer {
		f.scheduler("scheduled_answer", 10)
	}

	return nil
}

func (f *FakeApp) TriggerTimeInterrupt(key string) error {
	if key == "scheduled_answer" {
		f.ports[f.port].Write([]byte("answer"))
	}

	return nil
}

func main() {
	logger, _ := zap.NewDevelopment()

	sim, err := simulation.NewSimulation(logger, []simulation.Node{
		embedded.New(
			"first",
			&FakeApp{
				initializer: true,
				l:           logger.Named("first.app"),
				port:        "first.radio",
			},
			lora.New("first.radio", "first", 433.0, 20),
		),
		embedded.New(
			"second",
			&FakeApp{
				l:    logger.Named("second.app"),
				port: "second.radio",
			},
			lora.New("second.radio", "second", 433.0, 20),
		),
		//radio medium is also a node that can recieve messages
		//think of it as 'aether' anything that has radio can talk to it,
		//then it decides what simulation should recieve message (effectively duplicating messages)
		//based on node parameters (potentially simulation can have ports/interfaces, that would hold parameters/talk to 'aether')
		radio.NewRadioMedium(logger, 100*time.Millisecond),
	})

	if err != nil {
		logger.Fatal("Failed to create simulation", zap.Error(err))
	}

	sim.Run()
}
