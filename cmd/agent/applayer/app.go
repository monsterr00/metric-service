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
		api.gauge["Alloc"] = float64(memStats.Alloc)
		api.gauge["BuckHashSys"] = float64(memStats.BuckHashSys)
		api.gauge["Frees"] = float64(memStats.Frees)
		api.gauge["GCCPUFraction"] = memStats.GCCPUFraction
		api.gauge["GCSys"] = float64(memStats.GCSys)
		api.gauge["HeapAlloc"] = float64(memStats.HeapAlloc)
		api.gauge["HeapIdle"] = float64(memStats.HeapIdle)
		api.gauge["HeapInuse"] = float64(memStats.HeapInuse)
		api.gauge["HeapObjects"] = float64(memStats.HeapObjects)
		api.gauge["HeapReleased"] = float64(memStats.HeapReleased)
		api.gauge["HeapSys"] = float64(memStats.HeapSys)
		api.gauge["LastGC"] = float64(memStats.LastGC)
		api.gauge["Lookups"] = float64(memStats.Lookups)
		api.gauge["MCacheInuse"] = float64(memStats.MCacheInuse)
		api.gauge["MCacheSys"] = float64(memStats.MCacheSys)
		api.gauge["MSpanInuse"] = float64(memStats.MSpanInuse)
		api.gauge["MSpanSys"] = float64(memStats.MSpanSys)
		api.gauge["Mallocs"] = float64(memStats.Mallocs)
		api.gauge["NextGC"] = float64(memStats.NextGC)
		api.gauge["NumForcedGC"] = float64(memStats.NumForcedGC)
		api.gauge["NumGC"] = float64(memStats.NumGC)
		api.gauge["OtherSys"] = float64(memStats.OtherSys)
		api.gauge["PauseTotalNs"] = float64(memStats.PauseTotalNs)
		api.gauge["StackInuse"] = float64(memStats.StackInuse)
		api.gauge["StackSys"] = float64(memStats.StackSys)
		api.gauge["Sys"] = float64(memStats.Sys)
		api.gauge["TotalAlloc"] = float64(memStats.TotalAlloc)

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		api.gauge["RandomValue"] = random.Float64()

		api.counter["PollCount"] += 1
		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}
