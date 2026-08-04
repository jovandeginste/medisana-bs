package plugins

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/jovandeginste/medisana-bs/structs"
)

// CSV contains configuration for the CSV plugin
type CSV struct {
	Dir string
}

func (plugin CSV) Name() string {
	return "CSV"
}

func (plugin CSV) Logger() *slog.Logger {
	return slog.With("plugin", plugin.Name())
}

// Initialize the Csv plugin
func (plugin CSV) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["csv"]
	plugin.Dir = newc.Dir

	plugin.Logger().Debug("I am the CSV plugin")
	plugin.Logger().Debug("Config", "dir", plugin.Dir)

	return plugin
}

// ParseData will parse new data for a given person
func (plugin CSV) ParseData(person *structs.PersonMetrics) bool {
	plugin.Logger().Info("The csv plugin is parsing new data")

	personID := person.Person
	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))
	idx := 0

	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}

	sort.Sort(metrics)

	csvFile := plugin.Dir + "/" + strconv.Itoa(personID) + ".csv"
	plugin.Logger().Info("Writing to file", "path", csvFile)

	if err := createCsvDir(csvFile); err != nil {
		plugin.Logger().Error("Could not create CSV dir", "err", err)
		os.Exit(1)
	}

	f, err := os.Create(csvFile)
	if err != nil {
		plugin.Logger().Error("Failed to create CSV file", "err", err, "path", csvFile)
		return false
	}
	defer f.Close()

	if err := gocsv.MarshalWithoutHeaders(&metrics, f); err != nil {
		plugin.Logger().Error("Failed to write CSV", "err", err, "path", csvFile)
		return false
	}

	return true
}

func createCsvDir(file string) error {
	path := filepath.Dir(file)
	mode := os.FileMode(0o700)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(path, mode); mkErr != nil {
			return mkErr
		}
	}

	return nil
}

func (plugin CSV) InitializeData(_ *structs.PersonMetrics) bool {
	return true
}
