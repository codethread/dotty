package lib

import (
	"log"
	"os"
	"path"
)

type ImplicitConfig struct {
	Home           string
	ConfigLocation string
}

// GetImplicitConfig gathers config information based on user's system, these are assumed as defaults if not overridden
func GetImplicitConfig() ImplicitConfig {
	HOME, err := os.UserHomeDir()

	if err != nil {
		log.Fatal("Could not find a suitable HOME value")
	}

	ConfigLocation := path.Join(HOME, ".config", "dotty")

	return ImplicitConfig{
		Home:           HOME,
		ConfigLocation: ConfigLocation,
	}
}
