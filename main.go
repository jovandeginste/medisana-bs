package main

import (
	"encoding/binary"
	"log"
)

var metric_chan chan *PartialMetric

func main() {
	log.Println("Starting Bluetooth Scale monitor")

	metric_chan = make(chan *PartialMetric, 2)
	go func() {
		MetricParser()
	}()

	go func() {
		StartBluetooth()
		//FakeBluetooth()
	}()

	select {}
}

func updateData(person int, new_weights BodyMetrics) {
	cur_weights := ImportCsv(person)

	cur_weights = mergeSort(cur_weights, new_weights)

	ExportCsv(person, cur_weights)
}

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
