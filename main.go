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

func main() {
	log.SetLevel(log.DebugLevel)

	log.Infoln("[MAIN] Initializing Bluetooth Scale monitor")

	config = ReadConfig("config.toml")

	plugins.Initialize(config)

	metricChan = make(chan *structs.PartialMetric, 2)

	go MetricParser()

	log.Infoln("[MAIN] Starting Bluetooth Scale monitor")

	if config.Fakeit {
		FakeBluetooth()
	} else {
		StartBluetooth()
	}
}
