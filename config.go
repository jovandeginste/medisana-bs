package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/jovandeginste/medisana-bs/structs"
)

// ReadConfig reads the file and converts it to a struct
func ReadConfig(configfile string) structs.Config {
	var cfg structs.Config

	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
