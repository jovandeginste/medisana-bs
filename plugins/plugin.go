package plugins

import (
	"sync"

	log "github.com/sirupsen/logrus"

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

	log.Infoln("[PLUGIN] Initializing plugins")

	for name := range ap {
		log.Infof("[PLUGIN]  --> %s", name)

		pluginType, ok := pluginRegistry[name]
		if !ok {
			log.Infof("[PLUGIN]  *-> Unknown plugin: %s", name)
			continue
		}

		p, ok := pluginType.(structs.Plugin)
		if !ok {
			log.Infoln("[PLUGIN]  !-> FAILED")
			continue
		}

		allPlugins[name] = p.Initialize(configuration)

		log.Infoln("[PLUGIN]  *-> success")
	}

	log.Infoln("[PLUGIN] All plugins initialized.")
}

// ParseData will parse new data for a given person and send it to every configured plugin
func ParseData(person *structs.PersonMetrics) {
	log.Infoln("[PLUGIN] Sending data to all plugins")

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		log.Infof("[PLUGIN]  --> %s", name)

		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			if !p.ParseData(person) {
				log.Infof("[PLUGIN <%s>]  !-> FAILED", name)
				return
			}

			log.Infof("[PLUGIN <%s>]  *-> success", name)
		}(plugin, name)
	}

	wg.Wait()

	log.Infoln("[PLUGIN] All plugins parsed data.")
}

// InitializeData will send signal to all plugins that the data was initialized
func InitializeData(person *structs.PersonMetrics) {
	log.Infof("[PLUGIN] Sending initial data for %d (%s) to all plugins", person.Person, person.Name)

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		log.Infof("[PLUGIN]  --> %s", name)

		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			if !p.InitializeData(person) {
				log.Infof("[PLUGIN <%s>]  !-> FAILED", name)
				return
			}

			log.Infof("[PLUGIN <%s>]  *-> success", name)
		}(plugin, name)
	}

	wg.Wait()

	log.Infoln("[PLUGIN] All plugins parsed initial data.")
}
