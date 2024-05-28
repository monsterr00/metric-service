package applayer

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

type app struct {
	metrics map[string]models.Metric
	store   storelayer.Store
	m       sync.RWMutex
}

type App interface {
	Metrics() (map[string]models.Metric, error)
	SetMetrics()
	LockRW()
	UnlockRW()
}

func New(storeLayer storelayer.Store) *app {
	return &app{
		metrics: make(map[string]models.Metric),
		store:   storeLayer,
	}
}

func (api *app) Metrics() (map[string]models.Metric, error) {
	return api.metrics, nil
}

func (api *app) LockRW() {
	api.m.RLock()
}

func (api *app) UnlockRW() {
	api.m.RUnlock()
}

func (api *app) SetMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for {
		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)

		var metric models.Metric
		var metricValue float64

		api.m.Lock()

		metric.MType = "gauge"

		metric.ID = "Alloc"
		metricValue = float64(memStats.Alloc) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["Alloc"] = metric

		metric.ID = "BuckHashSys"
		metricValue = float64(memStats.BuckHashSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["BuckHashSys"] = metric

		metric.ID = "Frees"
		metricValue = float64(memStats.Frees) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["Frees"] = metric

		metric.ID = "GCCPUFraction"
		metricValue = memStats.GCCPUFraction + rand.Float64()
		metric.Value = &metricValue
		api.metrics["GCCPUFraction"] = metric

		metric.ID = "GCSys"
		metricValue = float64(memStats.GCSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["GCSys"] = metric

		metric.ID = "HeapAlloc"
		metricValue = float64(memStats.HeapAlloc) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapAlloc"] = metric

		metric.ID = "HeapIdle"
		metricValue = float64(memStats.HeapIdle) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapIdle"] = metric

		metric.ID = "HeapInuse"
		metricValue = float64(memStats.HeapInuse) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapInuse"] = metric

		metric.ID = "HeapObjects"
		metricValue = float64(memStats.HeapObjects) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapObjects"] = metric

		metric.ID = "HeapReleased"
		metricValue = float64(memStats.HeapReleased) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapReleased"] = metric

		metric.ID = "HeapSys"
		metricValue = float64(memStats.HeapSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["HeapSys"] = metric

		metric.ID = "LastGC"
		metricValue = float64(memStats.LastGC) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["LastGC"] = metric

		metric.ID = "Lookups"
		metricValue = float64(memStats.Lookups) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["Lookups"] = metric

		metric.ID = "MCacheInuse"
		metricValue = float64(memStats.MCacheInuse) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["MCacheInuse"] = metric

		metric.ID = "MCacheSys"
		metricValue = float64(memStats.MCacheSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["MCacheSys"] = metric

		metric.ID = "MSpanInuse"
		metricValue = float64(memStats.MSpanInuse) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["MSpanInuse"] = metric

		metric.ID = "MSpanSys"
		metricValue = float64(memStats.MSpanSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["MSpanSys"] = metric

		metric.ID = "Mallocs"
		metricValue = float64(memStats.Mallocs) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["Mallocs"] = metric

		metric.ID = "NextGC"
		metricValue = float64(memStats.NextGC) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["NextGC"] = metric

		metric.ID = "NumForcedGC"
		metricValue = float64(memStats.NumForcedGC) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["NumForcedGC"] = metric

		metric.ID = "NumGC"
		metricValue = float64(memStats.NumGC) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["NumGC"] = metric

		metric.ID = "OtherSys"
		metricValue = float64(memStats.OtherSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["OtherSys"] = metric

		metric.ID = "PauseTotalNs"
		metricValue = float64(memStats.PauseTotalNs) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["PauseTotalNs"] = metric

		metric.ID = "StackInuse"
		metricValue = float64(memStats.StackInuse) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["StackInuse"] = metric

		metric.ID = "StackSys"
		metricValue = float64(memStats.StackSys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["StackSys"] = metric

		metric.ID = "Sys"
		metricValue = float64(memStats.Sys) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["Sys"] = metric

		metric.ID = "TotalAlloc"
		metricValue = float64(memStats.TotalAlloc) + rand.Float64()
		metric.Value = &metricValue
		api.metrics["TotalAlloc"] = metric

		metric.ID = "RandomValue"
		metricValue = random.Float64()
		metric.Value = &metricValue
		api.metrics["RandomValue"] = metric

		metric.MType = "counter"
		metric.ID = "PollCount"

		var counter int64
		_, isSet := api.metrics["PollCount"]
		if isSet {
			counter = *api.metrics["PollCount"].Delta
		} else {
			counter = 0
		}

		counter += 1
		metric.Delta = &counter
		api.metrics["PollCount"] = metric

		api.m.Unlock()

		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}
