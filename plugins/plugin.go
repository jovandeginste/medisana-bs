package plugins

import (
	"log/slog"
	"sync"

	"github.com/jovandeginste/medisana-bs/structs"
)

var allPlugins map[string]structs.Plugin

var pluginRegistry = map[string]interface{}{
	"mail": Mail{},
	"csv":  CSV{},
	"mqtt": MQTT{},
}

func pluginLogger() *slog.Logger {
	return slog.With("component", "plugins")
}

// Initialize all plugins from configuration
func Initialize(configuration structs.Config) {
	ap := configuration.Plugins
	allPlugins = make(map[string]structs.Plugin)

	pluginLogger().Info("Initializing plugins")

	for name := range ap {
		pluginLogger().Info("Initializing plugin", "name", name)

		pluginType, ok := pluginRegistry[name]
		if !ok {
			pluginLogger().Warn("Unknown plugin", "name", name)
			continue
		}

		p, ok := pluginType.(structs.Plugin)
		if !ok {
			pluginLogger().Warn("Failed to cast plugin", "name", name)
			continue
		}

		allPlugins[name] = p.Initialize(configuration)

		pluginLogger().Info("Plugin initialized", "name", name)
	}

	pluginLogger().Info("All plugins initialized")
}

// ParseData will parse new data for a given person and send it to every configured plugin
func ParseData(person *structs.PersonMetrics) {
	pluginLogger().Info("Sending data to all plugins")

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		pluginLogger().Info("Parsing plugin data", "name", name)

		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			plog := pluginLogger().With("plugin", name)

			if !p.ParseData(person) {
				plog.Warn("Plugin parse failed")
				return
			}

			plog.Info("Plugin parse success")
		}(plugin, name)
	}

	wg.Wait()

	pluginLogger().Info("All plugins parsed data")
}

// InitializeData will send signal to all plugins that the data was initialized
func InitializeData(person *structs.PersonMetrics) {
	pluginLogger().Info("Sending initial data to all plugins", "person_id", person.Person, "name", person.Name)

	var wg sync.WaitGroup

	for name, plugin := range allPlugins {
		wg.Add(1)

		go func(p structs.Plugin, name string) {
			defer wg.Done()

			plog := pluginLogger().With("plugin", name)

			if !p.InitializeData(person) {
				plog.Warn("Plugin initial data failed")
				return
			}

			plog.Info("Plugin initial data success")
		}(plugin, name)
	}

	wg.Wait()

	pluginLogger().Info("All plugins parsed initial data")
}
