package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
)

var allPlugins structs.Plugins

func Initialize(ap structs.Plugins) {
	allPlugins = ap
	log.Println("Initializing plugins")
	for name, plugin := range allPlugins {
		log.Printf(" --> %s\n", name)
		result := plugin.Initialize()
		if result {
			log.Println(" *-> success")
		} else {
			log.Println(" !-> FAILED")
			os.Exit(1)
		}
	}
	log.Println("All plugins initialized.")
}

func ParseData(person *structs.PersonMetrics) {
	log.Println("Sending data to all plugins")
	for name, plugin := range allPlugins {
		log.Printf(" --> %s\n", name)
		result := plugin.ParseData(person)
		if result {
			log.Println(" *-> success")
		} else {
			log.Println(" !-> FAILED")
		}
	}
	log.Println("All plugins parsed data.")
}
