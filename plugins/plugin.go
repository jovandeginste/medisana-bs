package plugins

import (
	"log"
	"os"

	"github.com/jovandeginste/medisana-bs/structs"
)

var allPlugins map[string]structs.Plugin

var pluginRegistry = map[string]interface{}{
	"mail": Mail{},
	"csv":  Csv{},
}

// Initialize all plugins from configuration
func Initialize(configuration structs.Config) {
	ap := configuration.Plugins
	allPlugins = make(map[string]structs.Plugin)

	log.Println("[PLUGIN] Initializing plugins")
	for name := range ap {
		log.Printf("[PLUGIN]  --> %s\n", name)
		pluginType := pluginRegistry[name]
		allPlugins[name] = pluginType.(structs.Plugin).Initialize(configuration)
		if allPlugins[name] != nil {
			log.Println("[PLUGIN]  *-> success")
		} else {
			log.Println("[PLUGIN]  !-> FAILED")
			os.Exit(1)
		}
	}
	log.Println("[PLUGIN] All plugins initialized.")
}

// ParseData will parse new data for a given person and send it to every configured plugin
func ParseData(person *structs.PersonMetrics) {
	log.Println("[PLUGIN] Sending data to all plugins")
	for name, plugin := range allPlugins {
		log.Printf("[PLUGIN]  --> %s\n", name)
		result := plugin.ParseData(person)
		if result {
			log.Println("[PLUGIN]  *-> success")
		} else {
			log.Println("[PLUGIN]  !-> FAILED")
		}
	}
	log.Println("[PLUGIN] All plugins parsed data.")
}
