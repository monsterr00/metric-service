package applayer

import (
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
)

type app struct {
	gauge   map[string]float64
	counter map[string]int64
	store   storelayer.Store
}

type App interface {
	GetGaugeMetrics() (map[string]float64, error)
	GetCounterMetrics() (map[string]int64, error)
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
