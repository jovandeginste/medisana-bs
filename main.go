package main

import (
	"log/slog"
	"os"

	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
)

var (
	metricChan chan *structs.PartialMetric
	config     structs.Config
)

func mainLogger() *slog.Logger {
	return slog.With("component", "main")
}

func main() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	mainLogger().Info("Initializing Bluetooth Scale monitor")

	config = ReadConfig("config.toml")

	plugins.Initialize(config)

	metricChan = make(chan *structs.PartialMetric, 2)

	StartMetricParser()

	mainLogger().Info("Starting Bluetooth Scale monitor")

	runScanner()
}

func runScanner() {
	if config.Fakeit {
		FakeBluetooth()
		return
	}

	StartBluetooth()
}
