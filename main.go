package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
)

var (
	metricChan chan *structs.PartialMetric
	config     structs.Config
)

func mainLogger() log.FieldLogger {
	return log.WithField("component", "main")
}

func main() {
	log.SetLevel(log.DebugLevel)

	mainLogger().Infoln("Initializing Bluetooth Scale monitor")

	config = ReadConfig("config.toml")

	plugins.Initialize(config)

	metricChan = make(chan *structs.PartialMetric, 2)

	StartMetricParser()

	mainLogger().Infoln("Starting Bluetooth Scale monitor")

	runScanner()
}

func runScanner() {
	if config.Fakeit {
		FakeBluetooth()
		return
	}

	StartBluetooth()
}
