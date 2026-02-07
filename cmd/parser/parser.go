package main

import (
	"github.com/Gordy96/evt-sim/configuration"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

func main() {
	logCfg := zap.NewProductionConfig()
	logCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logCfg.EncoderConfig.MessageKey = "line"

	logger, _ := logCfg.Build()
	nodes, err := configuration.ParseFile("cmd/parser/config.hcl", logger)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	sim, err := simulation.NewSimulation(logger, nodes)

	if err != nil {
		logger.Sugar().Fatal(err)
	}

	sim.Run()
}
