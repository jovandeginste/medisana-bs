package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
)

func Initialize(allPlugins structs.Plugins) {
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
