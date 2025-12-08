package main

import (
	"github.com/Gordy96/evt-sim/configuration"
	"github.com/Gordy96/evt-sim/simulation"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
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
