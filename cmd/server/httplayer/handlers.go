package httplayer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

const (
	gaugeMetricType   = "gauge"
	counterMetricType = "counter"
)

func (api *httpAPI) getMainPage(res http.ResponseWriter, req *http.Request) {
	var body string
	var err error

	// считывем сохраненные метрики
	metrics, err := api.app.Metric()
	if err != nil {
		http.Error(res, "Server: error getting metrics %s\n", http.StatusNotFound)
		return
	}

	for k, v := range metrics {
		switch v.MType {
		case gaugeMetricType:
			body += fmt.Sprintf("%s: %v\r\n", k, *v.Value)
		case counterMetricType:
			body += fmt.Sprintf("%s: %v\r\n", k, *v.Delta)
		}
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte(body))
	if err != nil {
		fmt.Printf("Server: error writing request body %s\n", err)
	}
}

func (api *httpAPI) getMetric(res http.ResponseWriter, req *http.Request) {
	var err error

	// читаем тело запроса
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	var metric models.Metric
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// считывем сохраненные метрики
	metrics, err := api.app.Metric()
	if err != nil {
		fmt.Printf("Server: error getting metrics %s\n", err)
	}

	var resp []byte

	switch metric.MType {
	case gaugeMetricType, counterMetricType:
		_, isSet := metrics[metric.ID]
		if !isSet {
			http.Error(res, "No metric", http.StatusNotFound)
			return
		}
	default:
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	resp, err = json.Marshal(metrics[metric.ID])
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (api *httpAPI) postMetric(res http.ResponseWriter, req *http.Request) {
	var err error

	// читаем тело запроса
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	var metric models.Metric
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// считывем сохраненные метрики
	metrics, err := api.app.Metric()
	if err != nil {
		fmt.Printf("Server: error getting metrics %s\n", err)
	}

	var resp []byte

	switch metric.MType {
	case gaugeMetricType:
		metrics[metric.ID] = metric

	case counterMetricType:
		_, isSet := metrics[metric.ID]
		if isSet {
			var counter int64

			if metrics[metric.ID].Delta == nil {
				counter = 0
			} else {
				counter = *metrics[metric.ID].Delta
			}

			if metric.Delta != nil {
				metricCounter := metrics[metric.ID]
				counter += *metric.Delta
				metricCounter.Delta = &counter
				metrics[metric.ID] = metricCounter
			}
		} else {
			metrics[metric.ID] = metric
		}
	default:
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	resp, err = json.Marshal(metrics[metric.ID])
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	api.saveMetricsSync()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
