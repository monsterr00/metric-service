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

// Metrics возвращает набор метрик из мапы metrics.
func (api *app) Metrics() (map[string]models.Metric, error) {
	return api.metrics, nil
}

// LockRW устанавливает блокировки для горутины.
func (api *app) LockRW() {
	api.m.RLock()
}

// UnlockRW снимает блокировки для горутины.
func (api *app) UnlockRW() {
	api.m.RUnlock()
}

// GenMetricsStats формирует набор метрик с помощью пакета runtime.MemStats.
func (api *app) GenMetricsStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	api.updateMetric("Alloc", float64(memStats.Alloc)+rand.Float64())
	api.updateMetric("BuckHashSys", float64(memStats.BuckHashSys)+rand.Float64())
	api.updateMetric("Frees", float64(memStats.Frees)+rand.Float64())
	api.updateMetric("GCCPUFraction", memStats.GCCPUFraction+rand.Float64())
	api.updateMetric("GCSys", float64(memStats.GCSys)+rand.Float64())
	api.updateMetric("HeapAlloc", float64(memStats.HeapAlloc)+rand.Float64())
	api.updateMetric("HeapIdle", float64(memStats.HeapIdle)+rand.Float64())
	api.updateMetric("HeapInuse", float64(memStats.HeapInuse)+rand.Float64())
	api.updateMetric("HeapObjects", float64(memStats.HeapObjects)+rand.Float64())
	api.updateMetric("HeapReleased", float64(memStats.HeapReleased)+rand.Float64())
	api.updateMetric("HeapSys", float64(memStats.HeapSys)+rand.Float64())
	api.updateMetric("LastGC", float64(memStats.LastGC)+rand.Float64())
	api.updateMetric("Lookups", float64(memStats.Lookups)+rand.Float64())
	api.updateMetric("MCacheInuse", float64(memStats.MCacheInuse)+rand.Float64())
	api.updateMetric("MCacheSys", float64(memStats.MCacheSys)+rand.Float64())
	api.updateMetric("MSpanInuse", float64(memStats.MSpanInuse)+rand.Float64())
	api.updateMetric("MSpanSys", float64(memStats.MSpanSys)+rand.Float64())
	api.updateMetric("Mallocs", float64(memStats.Mallocs)+rand.Float64())
	api.updateMetric("NextGC", float64(memStats.NextGC)+rand.Float64())
	api.updateMetric("NumForcedGC", float64(memStats.NumForcedGC)+rand.Float64())
	api.updateMetric("NumGC", float64(memStats.NumGC)+rand.Float64())
	api.updateMetric("OtherSys", float64(memStats.OtherSys)+rand.Float64())
	api.updateMetric("PauseTotalNs", float64(memStats.PauseTotalNs)+rand.Float64())
	api.updateMetric("StackInuse", float64(memStats.StackInuse)+rand.Float64())
	api.updateMetric("StackSys", float64(memStats.StackSys)+rand.Float64())
	api.updateMetric("Sys", float64(memStats.Sys)+rand.Float64())
	api.updateMetric("TotalAlloc", float64(memStats.TotalAlloc)+rand.Float64())
	api.updateMetric("RandomValue", random.Float64())

	var metric models.Metric
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

// SetMetrics запускает процесс формирования метрик.
func (api *app) SetMetrics() {
	for {
		api.GenMetricsStats()
	}
}

// SetMetricsGOPSUTIL формирует набор метрик с помощью пакета runtime.MemStats.
func (api *app) SetMetricsGOPSUTIL() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	v, _ := mem.VirtualMemory()

	for {
		api.updateMetric("TotalMemory", float64(v.Total))
		api.updateMetric("FreeMemory", float64(v.Free))

		percent, _ := cpu.Percent(0, false)
		api.updateMetric("CPUutilization1", percent[0])

		time.Sleep(time.Duration(config.ClientOptions.PollInterval) * time.Second)
	}
}

// updateMetric добавляет метрику типа gauge в мапу с метриками.
func (api *app) updateMetric(id string, value float64) {
	metric := models.Metric{
		ID:    id,
		MType: "gauge",
		Value: &value,
	}
	api.m.Lock()
	api.metrics[id] = metric
	api.m.Unlock()
}
