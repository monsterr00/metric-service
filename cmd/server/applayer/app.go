package applayer

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

type app struct {
	metrics map[string]models.Metric
	store   storelayer.Store
}

type App interface {
	Metrics(cttx context.Context) (map[string]models.Metric, error)
	Metric(ctx context.Context, id string, mtype string) (models.Metric, error)
	SaveMetricsFile() error
	LoadMetricsFile() error
	LoadMetricsDB() error
	AddMetric(ctx context.Context, metric models.Metric) error
	PingDB() error
	CloseDB() error
}

func New(storeLayer storelayer.Store) *app {
	return &app{
		metrics: make(map[string]models.Metric),
		store:   storeLayer,
	}
}

func (api *app) Metrics(ctx context.Context) (map[string]models.Metric, error) {
	return api.metrics, nil
}

func (api *app) Metric(ctx context.Context, id string, mtype string) (models.Metric, error) {
	metric, isSet := api.metrics[id]
	if isSet {
		return metric, nil
	}

	return metric, errors.New("server: no metric")
}

func (api *app) AddMetric(ctx context.Context, metric models.Metric) error {
	var err error

	if config.ServerOptions.Mode == config.DBMode {
		_, isSet := api.metrics[metric.ID]

		if isSet {
			err = api.updateMetric(ctx, metric)
		} else {
			err = api.createMetric(ctx, metric)
		}
	}

	if err != nil {
		return err
	}

	var m sync.RWMutex
	m.RLock()
	api.metrics[metric.ID] = metric
	m.RUnlock()
	return nil
}

func (api *app) SaveMetricsFile() error {
	file, err := os.OpenFile(config.ServerOptions.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	metrics, err := api.Metrics(context.TODO())
	if err != nil {
		return err
	}

	for _, v := range metrics {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		// записываем событие в буфер
		if _, err := writer.Write(data); err != nil {
			return err
		}

		// добавляем перенос строки
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}

		// записываем буфер в файл
		if err := writer.Flush(); err != nil {
			return err
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (api *app) LoadMetricsFile() error {
	file, err := os.OpenFile(config.ServerOptions.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// читаем данные из scanner
		data := scanner.Bytes()
		metric := models.Metric{}
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}

		err := api.AddMetric(context.TODO(), metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *app) LoadMetricsDB() error {
	var err error
	ctx := context.Background()

	api.metrics, err = api.fetchMetrics(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (api *app) PingDB() error {
	err := api.store.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (api *app) CloseDB() error {
	err := api.store.Close()
	if err != nil {
		return err
	}
	return nil
}

func (api *app) createMetric(ctx context.Context, metric models.Metric) error {
	return api.store.Create(ctx, metric)
}

func (api *app) updateMetric(ctx context.Context, metric models.Metric) error {
	return api.store.Update(ctx, metric)
}

func (api *app) fetchMetrics(ctx context.Context) (map[string]models.Metric, error) {
	return api.store.Fetch(ctx)
}
