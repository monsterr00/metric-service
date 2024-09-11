package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/helpers"
)

func init() {
	flag.StringVar(&config.ClientOptions.Host, "a", "localhost:8080", "server host")
	flag.Int64Var(&config.ClientOptions.ReportInterval, "r", 2, "reportInterval value")
	flag.Int64Var(&config.ClientOptions.PollInterval, "p", 10, "pollInterval value")
	flag.StringVar(&config.ClientOptions.Key, "k", "", "secret key")
	flag.Int64Var(&config.ClientOptions.RateLimit, "l", 100, "max request pool")
	flag.StringVar(&config.ClientOptions.PublicKeyPath, "crypto-key", "internal/config/public.key", "public key path")
	flag.StringVar(&config.ClientOptions.ConfigJSONPath, "c", "internal/config/client_config.json", "client config JSON path")
	flag.StringVar(&config.ClientOptions.ConfigJSONPath, "config", "internal/config/client_config.json", "client config JSON path")
}

func setFlags() {
	var err error

	configJSON, isSet := os.LookupEnv("CONFIG")
	if isSet && configJSON != "" {
		config.ClientOptions.ConfigJSONPath = configJSON
	}
	jsonConfig, err := helpers.ReadConfigJSON(config.ClientOptions.ConfigJSONPath)
	if err != nil {
		log.Printf("Client: read error config json")
	}

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet && envAddress != "" {
		config.ClientOptions.Host = envAddress
	}
	if !isSet && !helpers.IsFlagPassed("a") && jsonConfig != nil {
		config.ClientOptions.Host = jsonConfig["address"]
	}

	envRepInterval, isSet := os.LookupEnv("REPORT_INTERVAL")
	if isSet {
		config.ClientOptions.ReportInterval, err = strconv.ParseInt(envRepInterval, 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for REPORT_INTERVAL")
		}
	}
	if !isSet && !helpers.IsFlagPassed("r") && jsonConfig != nil {
		config.ClientOptions.ReportInterval, err = strconv.ParseInt(jsonConfig["report_interval"], 10, 64)
		if err != nil {
			log.Printf("Client: wrong parametr type for REPORT_INTERVAL")
		}
	}

	envPollInterval, isSet := os.LookupEnv("POLL_INTERVAL")
	if isSet {
		config.ClientOptions.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for POLL_INTERVAL")
		}
	}
	if !isSet && !helpers.IsFlagPassed("p") && jsonConfig != nil {
		config.ClientOptions.PollInterval, err = strconv.ParseInt(jsonConfig["poll_interval"], 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for POLL_INTERVAL")
		}
	}

	secretKey, isSet := os.LookupEnv("KEY")
	if isSet {
		config.ClientOptions.Key = secretKey
	}
	if !isSet && !helpers.IsFlagPassed("k") && jsonConfig != nil {
		config.ClientOptions.Key = jsonConfig["secret_key"]
	}
	if config.ClientOptions.Key != "" {
		config.ClientOptions.SignMode = true
	}

	rateLimit, isSet := os.LookupEnv("RATE_LIMIT")
	if isSet {
		config.ClientOptions.RateLimit, err = strconv.ParseInt(rateLimit, 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for RATE_LIMIT")
		}
	}
	if !isSet && !helpers.IsFlagPassed("l") && jsonConfig != nil {
		config.ClientOptions.RateLimit, err = strconv.ParseInt(jsonConfig["rate_limit"], 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for RATE_LIMIT")
		}
	}

	cryptoKey, isSet := os.LookupEnv("CRYPTO_KEY")
	if isSet {
		config.ClientOptions.PublicKeyPath = cryptoKey
	}
	if !isSet && !helpers.IsFlagPassed("crypto-key") && jsonConfig != nil {
		config.ClientOptions.PublicKeyPath = jsonConfig["crypto_key"]
	}

	config.ClientOptions.BatchSize = 5
	config.ClientOptions.PoolWorkers = 10
	config.ClientOptions.PrivateKeyPath = "internal/config/private.key"
}
