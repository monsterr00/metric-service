package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

var gauge map[string]float64
var counter map[string]int64

func sendToServer(client *resty.Client) {
	for k, v := range gauge {
		metricGaugeURL := fmt.Sprintf("/update/gauge/%s/%f", k, v)
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, metricGaugeURL)

		req, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(requestURL)
		if err != nil {
			fmt.Printf("Client: error sending http-request: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Status code: %d\n", req.StatusCode())
	}
	for k, v := range counter {
		metricCounterURL := fmt.Sprintf("/update/counter/%s/%d", k, v)
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, metricCounterURL)

		req, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(requestURL)
		if err != nil {
			fmt.Printf("Client: error sending http-request: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Status code: %d\n", req.StatusCode())
	}
}

func gaugeMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for {
		gauge["Alloc"] = float64(memStats.Alloc)
		gauge["BuckHashSys"] = float64(memStats.BuckHashSys)
		gauge["Frees"] = float64(memStats.Frees)
		gauge["GCCPUFraction"] = memStats.GCCPUFraction
		gauge["GCSys"] = float64(memStats.GCSys)
		gauge["HeapAlloc"] = float64(memStats.HeapAlloc)
		gauge["HeapIdle"] = float64(memStats.HeapIdle)
		gauge["HeapInuse"] = float64(memStats.HeapInuse)
		gauge["HeapObjects"] = float64(memStats.HeapObjects)
		gauge["HeapReleased"] = float64(memStats.HeapReleased)
		gauge["HeapSys"] = float64(memStats.HeapSys)
		gauge["LastGC"] = float64(memStats.LastGC)
		gauge["Lookups"] = float64(memStats.Lookups)
		gauge["MCacheInuse"] = float64(memStats.MCacheInuse)
		gauge["MCacheSys"] = float64(memStats.MCacheSys)
		gauge["MSpanInuse"] = float64(memStats.MSpanInuse)
		gauge["MSpanSys"] = float64(memStats.MSpanSys)
		gauge["Mallocs"] = float64(memStats.Mallocs)
		gauge["NextGC"] = float64(memStats.NextGC)
		gauge["NumForcedGC"] = float64(memStats.NumForcedGC)
		gauge["NumGC"] = float64(memStats.NumGC)
		gauge["OtherSys"] = float64(memStats.OtherSys)
		gauge["PauseTotalNs"] = float64(memStats.PauseTotalNs)
		gauge["StackInuse"] = float64(memStats.StackInuse)
		gauge["StackSys"] = float64(memStats.StackSys)
		gauge["Sys"] = float64(memStats.Sys)
		gauge["TotalAlloc"] = float64(memStats.TotalAlloc)

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		gauge["RandomValue"] = random.Float64()

		counter["PollCount"] += 1
		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}

func init() {
	flag.StringVar(&config.ClientOptions.Host, "a", "localhost:8080", "server host")
	flag.Int64Var(&config.ClientOptions.ReportInterval, "r", 2, "reportInterval value")
	flag.Int64Var(&config.ClientOptions.PollInterval, "p", 10, "pollInterval value")

	var err error

	envAddress, isSet := os.LookupEnv("ADDRESS")
	if isSet {
		config.ClientOptions.Host = envAddress
	}

	envRepInterval, isSet := os.LookupEnv("REPORT_INTERVAL")
	if isSet {
		config.ClientOptions.ReportInterval, err = strconv.ParseInt(envRepInterval, 10, 64)
		if err != nil {
			fmt.Printf("Wrong parametr type for REPORT_INTERVAL")
			os.Exit(1)
		}
	}

	envPollInterval, isSet := os.LookupEnv("POLL_INTERVAL")
	if isSet {
		config.ClientOptions.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			fmt.Printf("wrong parametr type for POLL_INTERVAL")
			os.Exit(1)
		}
	}
}

func main() {
	gauge = make(map[string]float64)
	counter = make(map[string]int64)

	flag.Parse()

	go gaugeMetrics()

	client := resty.New()

	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second)

	for {
		sendToServer(client)
		time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
	}
}
