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

// Csv contains configuration for the Csv plugin
type Csv struct {
	Dir string
}

// Initialize the Csv plugin
func (plugin Csv) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["csv"]
	plugin.Dir = newc.Dir

	log.Debugln("[PLUGIN CSV] I am the CSV plugin")
	log.Debugf("[PLUGIN CSV]   - Dir: %s", plugin.Dir)

	return plugin
}

// ParseData will parse new data for a given person
func (plugin Csv) ParseData(person *structs.PersonMetrics) bool {
	log.Infoln("[PLUGIN CSV] The csv plugin is parsing new data")

	personID := person.Person
	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))
	idx := 0

	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}

	sort.Sort(metrics)

	csvFile := plugin.Dir + "/" + strconv.Itoa(personID) + ".csv"
	log.Infof("[PLUGIN CSV] Writing to file '%s'.", csvFile)
	createCsvDir(csvFile)

	f, err := os.Create(csvFile)
	if err != nil {
		log.Errorf("[PLUGIN CSV] %#v", err)
	}
	defer f.Close()

	if err := gocsv.MarshalWithoutHeaders(&metrics, f); err != nil {
		log.Errorf("[PLUGIN CSV] %#v", err)
	}

	return true
}

func createCsvDir(file string) {
	path := filepath.Dir(file)
	mode := os.FileMode(0o700)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(path, mode); mkErr != nil {
			log.Fatalf("[PLUGIN CSV] Could not create CSV dir: %s", mkErr)
		}
	}
}

func (plugin Csv) InitializeData(_ *structs.PersonMetrics) bool {
	return true
}
