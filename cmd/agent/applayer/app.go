package applayer

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type app struct {
	gauge   map[string]float64
	counter map[string]int64
	store   storelayer.Store
}

type App interface {
	GetGaugeMetrics() (map[string]float64, error)
	GetCounterMetrics() (map[string]int64, error)
	SetMetrics()
}

func New(storeLayer storelayer.Store) *app {
	return &app{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
		store:   storeLayer,
	}
}

func (api *app) GetGaugeMetrics() (map[string]float64, error) {
	return api.gauge, nil
}

func (api *app) GetCounterMetrics() (map[string]int64, error) {
	return api.counter, nil
}

func (api *app) SetMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for {
		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)

		api.gauge["Alloc"] = float64(memStats.Alloc) + rand.Float64()
		api.gauge["BuckHashSys"] = float64(memStats.BuckHashSys) + rand.Float64()
		api.gauge["Frees"] = float64(memStats.Frees) + rand.Float64()
		api.gauge["GCCPUFraction"] = memStats.GCCPUFraction + rand.Float64()
		api.gauge["GCSys"] = float64(memStats.GCSys) + rand.Float64()
		api.gauge["HeapAlloc"] = float64(memStats.HeapAlloc) + rand.Float64()
		api.gauge["HeapIdle"] = float64(memStats.HeapIdle) + rand.Float64()
		api.gauge["HeapInuse"] = float64(memStats.HeapInuse) + rand.Float64()
		api.gauge["HeapObjects"] = float64(memStats.HeapObjects) + rand.Float64()
		api.gauge["HeapReleased"] = float64(memStats.HeapReleased) + rand.Float64()
		api.gauge["HeapSys"] = float64(memStats.HeapSys) + rand.Float64()
		api.gauge["LastGC"] = float64(memStats.LastGC) + rand.Float64()
		api.gauge["Lookups"] = float64(memStats.Lookups) + rand.Float64()
		api.gauge["MCacheInuse"] = float64(memStats.MCacheInuse) + rand.Float64()
		api.gauge["MCacheSys"] = float64(memStats.MCacheSys) + rand.Float64()
		api.gauge["MSpanInuse"] = float64(memStats.MSpanInuse) + rand.Float64()
		api.gauge["MSpanSys"] = float64(memStats.MSpanSys) + rand.Float64()
		api.gauge["Mallocs"] = float64(memStats.Mallocs) + rand.Float64()
		api.gauge["NextGC"] = float64(memStats.NextGC) + rand.Float64()
		api.gauge["NumForcedGC"] = float64(memStats.NumForcedGC) + rand.Float64()
		api.gauge["NumGC"] = float64(memStats.NumGC) + rand.Float64()
		api.gauge["OtherSys"] = float64(memStats.OtherSys) + rand.Float64()
		api.gauge["PauseTotalNs"] = float64(memStats.PauseTotalNs) + rand.Float64()
		api.gauge["StackInuse"] = float64(memStats.StackInuse) + rand.Float64()
		api.gauge["StackSys"] = float64(memStats.StackSys) + rand.Float64()
		api.gauge["Sys"] = float64(memStats.Sys) + rand.Float64()
		api.gauge["TotalAlloc"] = float64(memStats.TotalAlloc) + rand.Float64()

		api.gauge["RandomValue"] = random.Float64()

		api.counter["PollCount"] += 1
		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}
