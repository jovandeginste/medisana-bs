package main

import (
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jovandeginste/medisana-bs/structs"
)

// ReadConfig reads the file and converts it to a struct
func ReadConfig(configfile string) structs.Config {
	var cfg structs.Config

	_, err := os.Stat(configfile)
	if err != nil {
		slog.Error("Config file is missing", "path", configfile)
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(configfile, &cfg); err != nil {
		slog.Error("Failed to decode config", "err", err)
		os.Exit(1)
	}

	return cfg
}
