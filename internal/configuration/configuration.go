package configuration

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"RUN_ADDRESS" envDefault:":3333"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./users.json"`
}

// New creates new Config
func New() *Config {
	var config = Config{}
	var err = env.Parse(&config)
	if err != nil {
		log.Printf("Error occurred when parsing config: %v", err)
	}

	flag.StringVar(&config.Address, "a", config.Address, "Launch address")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "Storage file path")
	flag.Parse()

	return &config
}
