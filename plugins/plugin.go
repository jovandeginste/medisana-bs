package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
)

var allPlugins structs.Plugins

func Initialize(ap structs.Plugins) {
	allPlugins = ap
	log.Println("[PLUGIN] Initializing plugins")
	for name, plugin := range allPlugins {
		log.Printf("[PLUGIN]  --> %s\n", name)
		result := plugin.Initialize()
		if result {
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
