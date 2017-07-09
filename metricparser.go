package main

import (
	"log"
)

func MetricParser(metric_chan <-chan PartialMetric) {
	for {
		partial_metric := <-metric_chan
		log.Printf("Received partial metric: %+v\n", partial_metric)
	}
}
