package plugins

import (
	//	"github.com/jovandeginste/medisana-bs/structs"
	"log"
)

type Mail struct {
	Server   string
	StartTLS bool
}

func (mail Mail) Initialize() bool {
	log.Println("I am the Mail plugin")
	log.Printf("  - Server: %s\n", mail.Server)
	log.Printf("  - StartTLS: %t\n", mail.StartTLS)
	return true
}
