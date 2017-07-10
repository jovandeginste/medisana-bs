package main

import (
	"encoding/binary"
	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
)

var metric_chan chan *structs.PartialMetric

func main() {
	log.Println("Starting Bluetooth Scale monitor")

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

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
