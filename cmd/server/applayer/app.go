package applayer

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

type app struct {
	metrics map[string]models.Metric
	store   storelayer.Store
}

type App interface {
	Metric() (map[string]models.Metric, error)
	SaveMetrics() error
	LoadMetrics() error
	AddMetrics(metric models.Metric)
}

func New(storeLayer storelayer.Store) *app {
	return &app{
		metrics: make(map[string]models.Metric),
		store:   storeLayer,
	}
}

func (api *app) Metric() (map[string]models.Metric, error) {
	return api.metrics, nil
}

func (api *app) AddMetrics(metric models.Metric) {
	api.metrics[metric.ID] = metric
}

func (api *app) SaveMetrics() error {
	file, err := os.OpenFile(config.ServerOptions.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	metrics, err := api.Metric()
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

func (api *app) LoadMetrics() error {
	file, err := os.OpenFile(config.ServerOptions.FileStoragePath, os.O_RDONLY|os.O_APPEND, 0666)
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

		api.AddMetrics(metric)
	}

	return nil
}
