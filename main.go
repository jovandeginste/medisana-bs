package main

import (
	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
)

var metric_chan chan *structs.PartialMetric
var config structs.Config

func main() {
	log.Println("[MAIN] Initializing Bluetooth Scale monitor")

	config = ReadConfig("config.toml")

	log.Printf("Configuration: %+v\n", config)

	plugins.Initialize(config.Plugins)
	metric_chan = make(chan *structs.PartialMetric, 2)
	go func() {
		MetricParser()
	}()

	log.Println("[MAIN] Starting Bluetooth Scale monitor")

	go func() {
		if config.Fakeit {
			FakeBluetooth()
		} else {
			StartBluetooth()
		}
	}()

	select {}
}
