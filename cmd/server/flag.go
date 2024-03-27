package main

import (
	"flag"
	"os"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

func init() {
	flag.StringVar(&config.ServerOptions.Host, "a", "localhost:8080", "server host")
	envAddress, isEnv := os.LookupEnv("ADDRESS")

	if isEnv {
		config.ServerOptions.Host = envAddress
	}
}
