package httplayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

const (
	gaugeMetricType     = "gauge"
	counterMetricType   = "counter"
	metricNamePosition  = 3
	metricValuePosition = 4
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

func (api *httpAPI) getMetricNoJSON(res http.ResponseWriter, req *http.Request) {
	// считывем сохраненные метрики
	metrics, err := api.app.Metric()
	if err != nil {
		fmt.Printf("Server: error getting metrics %s\n", err)
	}

	splitPath := strings.Split(req.URL.Path, "/")
	if len(splitPath) > metricNamePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName := splitPath[3]

		switch memType {
		case gaugeMetricType:
			memValue, isSet := metrics[memName]
			if isSet {
				_, err = res.Write([]byte(strconv.FormatFloat(*memValue.Value, 'f', -1, 64)))
				if err != nil {
					fmt.Printf("Server: error writing request body %s\n", err)
				}
			} else {
				http.Error(res, "No metric", http.StatusNotFound)
				return
			}
		case counterMetricType:
			memValue, isSet := metrics[memName]
			if isSet {
				_, err = res.Write([]byte(fmt.Sprintf("%d", *memValue.Delta)))
				if err != nil {
					fmt.Printf("Server: error writing request body %s\n", err)
				}
			} else {
				http.Error(res, "No metric", http.StatusNotFound)
				return
			}
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(res, "No metric type", http.StatusNotFound)
		return
	}
}

func (api *httpAPI) postMetricNoJSON(res http.ResponseWriter, req *http.Request) {
	splitPath := strings.Split(req.URL.Path, "/")

	metrics, err := api.app.Metric()
	if err != nil {
		fmt.Printf("Server: error getting metrics %s\n", err)
	}

	var memName string
	var memValue string

	if len(splitPath) > metricValuePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName = splitPath[3]
		// значение метрики
		memValue = splitPath[4]

		switch memType {
		case gaugeMetricType:
			if len(splitPath) > metricNamePosition {
				var err error
				var metric models.Metric

				metric.ID = memName
				metric.MType = gaugeMetricType
				metricValue, err := strconv.ParseFloat(memValue, 64)
				if err != nil {
					http.Error(res, "Wrong metric value", http.StatusBadRequest)
					return
				}
				metric.Value = &metricValue
				metrics[memName] = metric
			}
			_, err = res.Write([]byte(fmt.Sprintf("%f", *metrics[memName].Value)))
			if err != nil {
				fmt.Printf("Server: error writing request body %s\n", err)
			}
		case counterMetricType:
			var metric models.Metric

			metric.ID = memName
			metric.MType = gaugeMetricType
			counterValue, err := strconv.ParseInt(memValue, 10, 64)
			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
			}
			metric.Delta = &counterValue

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
					counter += counterValue
					metricCounter.Delta = &counter
					metrics[metric.ID] = metricCounter
				}
			} else {
				metrics[metric.ID] = metric
			}
			_, err = res.Write([]byte(fmt.Sprintf("%d", *metrics[memName].Delta)))
			if err != nil {
				fmt.Printf("Server: error writing request body %s\n", err)
			}
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(res, "No metric type or metric value", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
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
