package utils

import (
	"time"

	"github.com/pelletier/go-toml"
)

//Configuration - configuration structure
type Configuration struct {
	LogFilePath  string
	Port         string
	ServedURL    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

//LoadConfiguration - load configuration file
func LoadConfiguration(path string) *Configuration {
	config, err := toml.LoadFile(path)

	if err != nil {
		panic(err)
	}

	configuration := Configuration{}
	config.Unmarshal(&configuration)

	return &configuration
}
