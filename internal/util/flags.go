package util

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
	flag.StringVar(&config.ServerOptions.FileStoragePath, "f", "tmp/metrics-db.json", "server metric storage path")
	flag.BoolVar(&config.ServerOptions.Restore, "r", true, "server read metrics on start")
	flag.StringVar(&config.ServerOptions.DBaddress, "d", "host=postgres user=postgres password=postgres dbname=praktikum sslmode=disablee", "DB address")
	//flag.StringVar(&config.ServerOptions.DBaddress, "d", "host=localhost user=postgres password=postgres1 dbname=metrics sslmode=disable", "DB address")
	flag.StringVar(&config.ServerOptions.Key, "k", "", "secret key")
}

// SetFlags инициализирует настройки программы.
func SetFlags() {
	var err error

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet && envAddress != "" {
		config.ServerOptions.Host = envAddress
	}

	envInterval, isSet := os.LookupEnv("STORE_INTERVAL")
	if isSet {
		config.ServerOptions.StoreInterval, err = strconv.ParseInt(envInterval, 10, 64)
		if err != nil {
			log.Printf("Server: wrong parametr type for STORE_INTERVAL")
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
			log.Printf("Server: wrong parametr type for RESTORE")
		}
	}

	dbAddress, isSet := os.LookupEnv("DATABASE_DSN")
	if isSet && dbAddress != "" {
		config.ServerOptions.DBaddress = dbAddress
	}

	if config.ServerOptions.DBaddress != "" {
		config.ServerOptions.Mode = config.DBMode
	} else if config.ServerOptions.FileStoragePath != "" && config.ServerOptions.Restore {
		config.ServerOptions.Mode = config.FileMode
	} else {
		config.ServerOptions.Mode = config.MemoryMode
	}

	secretKey, isSet := os.LookupEnv("KEY")
	if isSet {
		config.ServerOptions.Key = secretKey
	}

	if config.ServerOptions.Key != "" {
		config.ServerOptions.SignMode = true
	}

	config.ServerOptions.ReconnectCount = 3
	config.ServerOptions.ReconnectDelta = 2
}
