package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort      string `envconfig:"HTTP_PORT" default:"8080"`
	StorageEndpoint  string `envconfig:"STORAGE_URL" default:"localhost:2379"`
	DBUrl   string `envconfig:"DB_URL"`
}

func LoadConfig() *Config {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &conf
}
