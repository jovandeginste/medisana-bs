package plugins

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/jovandeginste/medisana-bs/structs"
)

var allPlugins map[string]structs.Plugin

var pluginRegistry = map[string]interface{}{
	"mail": Mail{},
	"csv":  CSV{},
	"mqtt": MQTT{},
}

func pluginLogger() log.FieldLogger {
	return log.WithField("component", "plugins")
}

// Initialize all plugins from configuration
func Initialize(configuration structs.Config) {
	ap := configuration.Plugins
	allPlugins = make(map[string]structs.Plugin)

	pluginLogger().Infoln("Initializing plugins")

	for name := range ap {
		pluginLogger().Infof(" --> %s", name)

		pluginType, ok := pluginRegistry[name]
		if !ok {
			pluginLogger().Infof(" *-> Unknown plugin: %s", name)
			continue
		}

		p, ok := pluginType.(structs.Plugin)
		if !ok {
			pluginLogger().Infoln(" !-> FAILED")
			continue
		}

		allPlugins[name] = p.Initialize(configuration)

		pluginLogger().Infoln(" *-> success")
	}

	pluginLogger().Infoln("All plugins initialized.")
}

// ParseData will parse new data for a given person and send it to every configured plugin
func ParseData(person *structs.PersonMetrics) {
	pluginLogger().Infoln("Sending data to all plugins")

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		pluginLogger().Infof(" --> %s", name)

		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			if !p.ParseData(person) {
				pluginLogger().Infof(" !-> '%s' FAILED", name)
				return
			}

			pluginLogger().Infof(" *-> '%s' success", name)
		}(plugin, name)
	}

	wg.Wait()

	pluginLogger().Infoln("All plugins parsed data.")
}

// InitializeData will send signal to all plugins that the data was initialized
func InitializeData(person *structs.PersonMetrics) {
	pluginLogger().Infof("Sending initial data for %d (%s) to all plugins", person.Person, person.Name)

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			if !p.InitializeData(person) {
				pluginLogger().Infof(" !-> '%s' FAILED", name)
				return
			}

			pluginLogger().Infof(" *-> '%s' success", name)
		}(plugin, name)
	}

	wg.Wait()

	pluginLogger().Infoln("All plugins parsed initial data.")
}
