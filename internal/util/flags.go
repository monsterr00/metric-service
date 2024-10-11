package util

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/helpers"
)

func init() {
	flag.StringVar(&config.ServerOptions.Host, "a", "localhost:8080", "server host")
	flag.Int64Var(&config.ServerOptions.StoreInterval, "i", 300, "server file store interval")
	flag.StringVar(&config.ServerOptions.FileStoragePath, "f", "tmp/metrics-db.json", "server metric storage path")
	flag.BoolVar(&config.ServerOptions.Restore, "r", true, "server read metrics on start")
	flag.StringVar(&config.ServerOptions.DBaddress, "d", "", "DB address")
	flag.StringVar(&config.ServerOptions.Key, "k", "", "secret key")
	flag.StringVar(&config.ServerOptions.PrivateKeyPath, "crypto-key", "/internal/config/private.key", "private key path")
	flag.StringVar(&config.ServerOptions.ConfigJSONPath, "c", "internal/config/server_config.json", "server config JSON path")
	flag.StringVar(&config.ServerOptions.ConfigJSONPath, "config", "internal/config/server_config.json", "server config JSON path")
}

// SetFlags инициализирует настройки программы.
func SetFlags() {
	var err error

	configJSON, isSet := os.LookupEnv("CONFIG")
	if isSet && configJSON != "" {
		config.ServerOptions.ConfigJSONPath = configJSON
	}
	jsonConfig, err := helpers.ReadConfigJSON(config.ServerOptions.ConfigJSONPath)
	if err != nil {
		log.Printf("Server: read error config json")
	}

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet && envAddress != "" {
		config.ServerOptions.Host = envAddress
	}
	if !isSet && !helpers.IsFlagPassed("a") && jsonConfig != nil {
		config.ServerOptions.Host = jsonConfig["address"]
	}

	envInterval, isSet := os.LookupEnv("STORE_INTERVAL")
	if isSet {
		config.ServerOptions.StoreInterval, err = strconv.ParseInt(envInterval, 10, 64)
		if err != nil {
			log.Printf("Server: wrong parametr type for STORE_INTERVAL")
		}
	}
	if !isSet && !helpers.IsFlagPassed("i") && jsonConfig != nil {
		config.ServerOptions.StoreInterval, err = strconv.ParseInt(jsonConfig["store_interval"], 10, 64)
		if err != nil {
			log.Printf("Server: wrong parametr type for STORE_INTERVAL")
		}
	}

	envFilePath, isSet := os.LookupEnv("FILE_STORAGE_PATH")
	if isSet {
		config.ServerOptions.FileStoragePath = envFilePath
	}
	if !isSet && !helpers.IsFlagPassed("f") && jsonConfig != nil {
		config.ServerOptions.FileStoragePath = jsonConfig["store_file"]
	}

	envRestore, isSet := os.LookupEnv("RESTORE")
	if isSet {
		config.ServerOptions.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			log.Printf("Server: wrong parametr type for RESTORE")
		}
	}
	if !isSet && !helpers.IsFlagPassed("r") && jsonConfig != nil {
		config.ServerOptions.Restore, err = strconv.ParseBool(jsonConfig["restore"])
		if err != nil {
			log.Printf("Server: wrong parametr type for RESTORE")
		}
	}

	dbAddress, isSet := os.LookupEnv("DATABASE_DSN")
	if isSet && dbAddress != "" {
		config.ServerOptions.DBaddress = dbAddress
	}
	if !isSet && !helpers.IsFlagPassed("d") && jsonConfig != nil {
		config.ServerOptions.DBaddress = jsonConfig["database_dsn"]
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
	if !isSet && !helpers.IsFlagPassed("k") && jsonConfig != nil {
		config.ServerOptions.Key = jsonConfig["secret_key"]
	}
	if config.ServerOptions.Key != "" {
		config.ServerOptions.SignMode = true
	}

	cryptoKey, isSet := os.LookupEnv("CRYPTO_KEY")
	if isSet {
		config.ServerOptions.PrivateKeyPath = cryptoKey
	}
	if !isSet && !helpers.IsFlagPassed("crypto-key") && jsonConfig != nil {
		config.ServerOptions.PrivateKeyPath = jsonConfig["crypto_key"]
	}

	config.ServerOptions.ReconnectCount = 3
	config.ServerOptions.ReconnectDelta = 2
}
