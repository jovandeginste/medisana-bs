package plugins

import (
	"log"
	"sync"

	"github.com/jovandeginste/medisana-bs/structs"
)

var allPlugins map[string]structs.Plugin

var pluginRegistry = map[string]interface{}{
	"mail": Mail{},
	"csv":  Csv{},
	"mqtt": MQTT{},
}

// Initialize all plugins from configuration
func Initialize(configuration structs.Config) {
	ap := configuration.Plugins
	allPlugins = make(map[string]structs.Plugin)

	log.Println("[PLUGIN] Initializing plugins")

	for name := range ap {
		log.Printf("[PLUGIN]  --> %s\n", name)

		pluginType, ok := pluginRegistry[name]
		if !ok {
			log.Printf("[PLUGIN]  *-> Unknown plugin: %s", name)
			continue
		}

		p, ok := pluginType.(structs.Plugin)
		if !ok {
			log.Println("[PLUGIN]  !-> FAILED")
			continue
		}

		allPlugins[name] = p.Initialize(configuration)

		log.Println("[PLUGIN]  *-> success")
	}

	log.Println("[PLUGIN] All plugins initialized.")
}

// ParseData will parse new data for a given person and send it to every configured plugin
func ParseData(person *structs.PersonMetrics) {
	log.Println("[PLUGIN] Sending data to all plugins")

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		log.Printf("[PLUGIN]  --> %s\n", name)

		wg.Add(1)

		go func(p structs.Plugin) {
			defer wg.Done()

			if !p.ParseData(person) {
				log.Println("[PLUGIN]  !-> FAILED")
				return
			}

			log.Println("[PLUGIN]  *-> success")
		}(plugin)
	}

	wg.Wait()

	log.Println("[PLUGIN] All plugins parsed data.")
}
