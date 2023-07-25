package main

import (
	"log"

	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
)

var (
	metricChan chan *structs.PartialMetric
	config     structs.Config
)

func main() {
	log.Println("[MAIN] Initializing Bluetooth Scale monitor")

	config = ReadConfig("config.toml")

	plugins.Initialize(config)
	metricChan = make(chan *structs.PartialMetric, 2)
	go MetricParser()

	log.Println("[MAIN] Starting Bluetooth Scale monitor")

	if config.Fakeit {
		FakeBluetooth()
	} else {
		StartBluetooth()
	}
}
