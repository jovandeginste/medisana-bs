package plugins

import (
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
)

type Mail struct {
	Server   string
	StartTLS bool
}

func (plugin Mail) Initialize() bool {
	log.Println("I am the Mail plugin")
	log.Printf("  - Server: %s\n", plugin.Server)
	log.Printf("  - StartTLS: %t\n", plugin.StartTLS)
	return true
}
func (plugin Mail) ParseData(person *structs.PersonMetrics) bool {
	log.Println("The mail plugin is parsing new data")
	return true
}
