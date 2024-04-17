package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

func init() {
	flag.StringVar(&config.ServerOptions.Host, "a", "localhost:8080", "server host")
	flag.Int64Var(&config.ServerOptions.StoreInterval, "i", 300, "server file store interval")
	flag.StringVar(&config.ServerOptions.FileStoragePath, "f", "/Users/denis/metric-service/tmp/metrics-db.json", "server metric storage path")
	flag.BoolVar(&config.ServerOptions.Restore, "r", true, "server read metrics on start")

	var err error

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet {
		config.ServerOptions.Host = envAddress
	}

	envInterval, isSet := os.LookupEnv("STORE_INTERVAL")
	if isSet {
		config.ServerOptions.StoreInterval, err = strconv.ParseInt(envInterval, 10, 64)
		if err != nil {
			log.Printf("Servr: wrong parametr type for STORE_INTERVAL")
		}
	}

	envFilePath, isSet := os.LookupEnv("FILE_STORAGE_PATH")
	if isSet {
		config.ServerOptions.FileStoragePath = envFilePath
	}

	envRestore, isSet := os.LookupEnv("RESTORE")
	if isSet {
		config.ServerOptions.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			log.Printf("Servr: wrong parametr type for RESTORE")
		}
	}
}
