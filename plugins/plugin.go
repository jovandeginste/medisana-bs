package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
	"reflect"
)

var allPlugins map[string]structs.Plugin

var pluginRegistry = map[string]reflect.Type{
	"mail": reflect.TypeOf(Mail{}),
	"csv":  reflect.TypeOf(Csv{}),
}

func Initialize(configuration interface{}) {
	ap := configuration.(map[string]interface{})
	allPlugins = make(map[string]structs.Plugin)

	log.Println("[PLUGIN] Initializing plugins")
	for name, plugin_config := range ap {
		log.Printf("[PLUGIN]  --> %s\n", name)
		plugin_type := pluginRegistry[name]
		plugin_builder := reflect.New(plugin_type).Elem().Interface().(structs.Plugin)
		allPlugins[name] = plugin_builder.Initialize(plugin_config)
		if allPlugins[name] != nil {
			log.Println("[PLUGIN]  *-> success")
		} else {
			log.Println("[PLUGIN]  !-> FAILED")
			os.Exit(1)
		}
	}
	log.Println("[PLUGIN] All plugins initialized.")
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
