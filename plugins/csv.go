package plugins

import (
	"github.com/gocarina/gocsv"
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// Csv contains configuration for the Csv plugin
type Csv struct {
	Dir string
}

// Initialize the Csv plugin
func (plugin Csv) Initialize(c interface{}) structs.Plugin {
	newc := c.(map[string]interface{})
	plugin.Dir = newc["Dir"].(string)
	log.Println("[PLUGIN CSV] I am the CSV plugin")
	log.Printf("[PLUGIN CSV]   - Dir: %s\n", plugin.Dir)
	return plugin
}

// ParseData will parse new data for a given person
func (plugin Csv) ParseData(person *structs.PersonMetrics) bool {
	log.Println("[PLUGIN CSV] The csv plugin is parsing new data")
	personID := person.Person
	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))
	idx := 0
	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}
	sort.Sort(metrics)

	csvFile := plugin.Dir + "/" + strconv.Itoa(personID) + ".csv"
	log.Printf("[PLUGIN CSV] Writing to file '%s'.\n", csvFile)
	createCsvDir(csvFile)

	f, err := os.Create(csvFile)
	if err != nil {
		log.Printf("[PLUGIN CSV] %#v", err)
	}
	defer f.Close()

	err = gocsv.MarshalWithoutHeaders(&metrics, f)

	if err != nil {
		log.Printf("[PLUGIN CSV] %#v", err)
	}
	return true
}

func createCsvDir(file string) {
	path := filepath.Dir(file)
	mode := os.FileMode(0700)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, mode)
	}
}
