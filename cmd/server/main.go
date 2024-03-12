package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/handlers"
)

func init() {
	flag.StringVar(&config.ServerOptions.Host, "a", "localhost:8080", "server host")
	envAddress, isEnv := os.LookupEnv("ADDRESS")

	if isEnv {
		config.ServerOptions.Host = envAddress
	}
}

func main() {
	config.MemStorage.Gauge = make(map[string]float64)
	config.MemStorage.Counter = make(map[string]int64)

	flag.Parse()

	r := chi.NewRouter()
	r.Get("/", handlers.MainPage)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.PostMetric)
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetric)

	err := http.ListenAndServe(config.ServerOptions.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
