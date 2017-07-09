package main

import (
	"encoding/binary"
	"log"
	"strconv"
)

var metric_chan chan *PartialMetric

func main() {
	log.Println("Starting Bluetooth Scale monitor")

	metric_chan = make(chan *PartialMetric, 2)
	go func() {
		MetricParser()
	}()
	/*
		for i := 1; i < 8; i++ {
			go func(i int) {
				new_weights := ImportCsv("csv/" + strconv.Itoa(i) + ".csv")
				updateData(i, new_weights)
			}(i)
		}
	*/

	go func() {
		//StartBluetooth()
		FakeBluetooth()
	}()

	select {}
}

func updateData(person int, new_weights BodyMetrics) {
	csvFile := "csv/" + strconv.Itoa(person) + ".csv"
	cur_weights := ImportCsv(csvFile)

	cur_weights = mergeSort(cur_weights, new_weights)

	ExportCsv(csvFile+".new", cur_weights)
}

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
