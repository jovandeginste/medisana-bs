package main

import (
	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
)

var metric_chan chan *structs.PartialMetric

func main() {
	log.Println("[MAIN] Starting Bluetooth Scale monitor")

	plugins.Initialize(allPlugins)
	metric_chan = make(chan *structs.PartialMetric, 2)
	go func() {
		MetricParser()
	}()

	go func() {
		//StartBluetooth()
		FakeBluetooth()
	}()

	select {}
}
