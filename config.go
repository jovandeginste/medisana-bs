package main

import (
	"github.com/BurntSushi/toml"
	"github.com/jovandeginste/medisana-bs/structs"
	"log"
	"os"
)

const MaxInt = 2147483647

func ReadConfig(configfile string) structs.Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config structs.Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
