package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort          string `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort          string `envconfig:"GRPC_PORT" default:"9090"`
	StorageEndpoint   string `envconfig:"STORAGE_URL" default:"localhost:2379"`
	PostgresURL       string `envconfig:"POSTGRES_URL"`
	FlagServicePrefix string `envconfig:"FLAG_SERVICE_PREFIX" default:"/featureflags/"`
	APIVersion        string `envconfig:"API_VERSION" default:"v1"`
}

func LoadConfig() *Config {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &conf
}
