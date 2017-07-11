package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
)

var allPlugins structs.Plugins

var PluginMap = map[string](func(c interface{}) structs.Plugin){
	"mail": MailPlugin,
	"csv":  CsvPlugin,
}

func Initialize(configuration interface{}) {
	ap := configuration.(map[string]interface{})
	allPlugins = structs.Plugins{}

	log.Printf("%+v\n", ap)
	log.Println("[PLUGIN] Initializing plugins")
	for name, plugin_config := range ap {
		log.Printf("[PLUGIN]  --> %s\n", name)
		plugin_builder := PluginMap[name]
		plugin := plugin_builder(plugin_config)
		result := plugin.Initialize()
		allPlugins[name] = plugin
		if result {
			log.Println("[PLUGIN]  *-> success")
		} else {
			log.Println("[PLUGIN]  !-> FAILED")
			os.Exit(1)
		}
	}
	log.Println("[PLUGIN] All plugins initialized.")
	log.Printf("%+v\n", allPlugins)
}

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
