package applayer

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type app struct {
	metrics map[string]models.Metric
	store   storelayer.Store
	m       sync.RWMutex
}

type App interface {
	Metrics() (map[string]models.Metric, error)
	SetMetrics()
	SetMetricsGOPSUTIL()
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

		metric.MType = "gauge"

		metric.ID = "Alloc"
		metricValue = float64(memStats.Alloc) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["Alloc"] = metric
		api.m.Unlock()

		metric.ID = "BuckHashSys"
		metricValue = float64(memStats.BuckHashSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["BuckHashSys"] = metric
		api.m.Unlock()

		metric.ID = "Frees"
		metricValue = float64(memStats.Frees) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["Frees"] = metric
		api.m.Unlock()

		metric.ID = "GCCPUFraction"
		metricValue = memStats.GCCPUFraction + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["GCCPUFraction"] = metric
		api.m.Unlock()

		metric.ID = "GCSys"
		metricValue = float64(memStats.GCSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["GCSys"] = metric
		api.m.Unlock()

		metric.ID = "HeapAlloc"
		metricValue = float64(memStats.HeapAlloc) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapAlloc"] = metric
		api.m.Unlock()

		metric.ID = "HeapIdle"
		metricValue = float64(memStats.HeapIdle) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapIdle"] = metric
		api.m.Unlock()

		metric.ID = "HeapInuse"
		metricValue = float64(memStats.HeapInuse) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapInuse"] = metric
		api.m.Unlock()

		metric.ID = "HeapObjects"
		metricValue = float64(memStats.HeapObjects) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapObjects"] = metric
		api.m.Unlock()

		metric.ID = "HeapReleased"
		metricValue = float64(memStats.HeapReleased) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapReleased"] = metric
		api.m.Unlock()

		metric.ID = "HeapSys"
		metricValue = float64(memStats.HeapSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["HeapSys"] = metric
		api.m.Unlock()

		metric.ID = "LastGC"
		metricValue = float64(memStats.LastGC) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["LastGC"] = metric
		api.m.Unlock()

		metric.ID = "Lookups"
		metricValue = float64(memStats.Lookups) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["Lookups"] = metric
		api.m.Unlock()

		metric.ID = "MCacheInuse"
		metricValue = float64(memStats.MCacheInuse) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["MCacheInuse"] = metric
		api.m.Unlock()

		metric.ID = "MCacheSys"
		metricValue = float64(memStats.MCacheSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["MCacheSys"] = metric
		api.m.Unlock()

		metric.ID = "MSpanInuse"
		metricValue = float64(memStats.MSpanInuse) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["MSpanInuse"] = metric
		api.m.Unlock()

		metric.ID = "MSpanSys"
		metricValue = float64(memStats.MSpanSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["MSpanSys"] = metric
		api.m.Unlock()

		metric.ID = "Mallocs"
		metricValue = float64(memStats.Mallocs) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["Mallocs"] = metric
		api.m.Unlock()

		metric.ID = "NextGC"
		metricValue = float64(memStats.NextGC) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["NextGC"] = metric
		api.m.Unlock()

		metric.ID = "NumForcedGC"
		metricValue = float64(memStats.NumForcedGC) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["NumForcedGC"] = metric
		api.m.Unlock()

		metric.ID = "NumGC"
		metricValue = float64(memStats.NumGC) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["NumGC"] = metric
		api.m.Unlock()

		metric.ID = "OtherSys"
		metricValue = float64(memStats.OtherSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["OtherSys"] = metric
		api.m.Unlock()

		metric.ID = "PauseTotalNs"
		metricValue = float64(memStats.PauseTotalNs) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["PauseTotalNs"] = metric
		api.m.Unlock()

		metric.ID = "StackInuse"
		metricValue = float64(memStats.StackInuse) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["StackInuse"] = metric
		api.m.Unlock()

		metric.ID = "StackSys"
		metricValue = float64(memStats.StackSys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["StackSys"] = metric
		api.m.Unlock()

		metric.ID = "Sys"
		metricValue = float64(memStats.Sys) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["Sys"] = metric
		api.m.Unlock()

		metric.ID = "TotalAlloc"
		metricValue = float64(memStats.TotalAlloc) + rand.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["TotalAlloc"] = metric
		api.m.Unlock()

		metric.ID = "RandomValue"
		metricValue = random.Float64()
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["RandomValue"] = metric
		api.m.Unlock()

		metric.MType = "counter"
		metric.ID = "PollCount"

		var counter int64
		api.m.RLock()
		_, isSet := api.metrics["PollCount"]
		api.m.RUnlock()
		if isSet {
			api.m.RLock()
			counter = *api.metrics["PollCount"].Delta
			api.m.RUnlock()
		} else {
			counter = 0
		}

		counter += 1
		metric.Delta = &counter
		api.m.Lock()
		api.metrics["PollCount"] = metric
		api.m.Unlock()

		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}

func (api *app) SetMetricsGOPSUTIL() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	v, _ := mem.VirtualMemory()

	for {
		var metric models.Metric
		var metricValue float64

		metric.MType = "gauge"

		metric.ID = "TotalMemory"
		metricValue = float64(v.Total)
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["TotalMemory"] = metric
		api.m.Unlock()

		metric.ID = "FreeMemory"
		metricValue = float64(v.Free)
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["FreeMemory"] = metric
		api.m.Unlock()

		percent, _ := cpu.Percent(0, false)

		metric.ID = "CPUutilization1"
		metricValue = percent[0]
		metric.Value = &metricValue
		api.m.Lock()
		api.metrics["CPUutilization1"] = metric
		api.m.Unlock()

		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}
