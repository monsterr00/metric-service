package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

func init() {
	flag.StringVar(&config.ClientOptions.Host, "a", "localhost:8080", "server host")
	flag.Int64Var(&config.ClientOptions.ReportInterval, "r", 2, "reportInterval value")
	flag.Int64Var(&config.ClientOptions.PollInterval, "p", 10, "pollInterval value")

	var err error

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet && envAddress != "" {
		config.ClientOptions.Host = envAddress
	}

	envRepInterval, isSet := os.LookupEnv("REPORT_INTERVAL")
	if isSet {
		config.ClientOptions.ReportInterval, err = strconv.ParseInt(envRepInterval, 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for REPORT_INTERVAL")
		}
	}

	envPollInterval, isSet := os.LookupEnv("POLL_INTERVAL")
	if isSet {
		config.ClientOptions.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Printf("Wrong parametr type for POLL_INTERVAL")
		}
	}

	config.ClientOptions.BatchSize = 5
}
