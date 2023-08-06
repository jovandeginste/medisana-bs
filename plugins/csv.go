package plugins

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"

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

func (plugin CSV) Logger() log.FieldLogger {
	return log.WithField("plugin", plugin.Name())
}

// Initialize the Csv plugin
func (plugin CSV) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["csv"]
	plugin.Dir = newc.Dir

	plugin.Logger().Debugln("I am the CSV plugin")
	plugin.Logger().Debugf("  - Dir: %s", plugin.Dir)

	return plugin
}

// ParseData will parse new data for a given person
func (plugin CSV) ParseData(person *structs.PersonMetrics) bool {
	plugin.Logger().Infoln("The csv plugin is parsing new data")

	personID := person.Person
	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))
	idx := 0

	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}

	sort.Sort(metrics)

	csvFile := plugin.Dir + "/" + strconv.Itoa(personID) + ".csv"
	plugin.Logger().Infof("Writing to file '%s'.", csvFile)

	if err := createCsvDir(csvFile); err != nil {
		plugin.Logger().Fatalf("Could not create CSV dir: %s", err)
	}

	f, err := os.Create(csvFile)
	if err != nil {
		plugin.Logger().Errorf("%#v", err)
	}
	defer f.Close()

	if err := gocsv.MarshalWithoutHeaders(&metrics, f); err != nil {
		plugin.Logger().Errorf("%#v", err)
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
