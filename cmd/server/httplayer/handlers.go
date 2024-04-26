package httplayer

import (
	"bytes"
	"context"
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

	ctx := req.Context()

	// считывем сохраненные метрики
	metrics, err := api.app.Metrics(ctx)
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

	ctx := req.Context()
	var savedMetric models.Metric

	switch metric.MType {
	case gaugeMetricType, counterMetricType:
		savedMetric, err = api.app.Metric(ctx, metric.ID, metric.MType)
		if err != nil {
			http.Error(res, "No metric", http.StatusNotFound)
			return
		}
	default:
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(savedMetric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		fmt.Printf("Server: error writing request body %s\n", err)
	}
}

func (api *httpAPI) getMetricNoJSON(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	splitPath := strings.Split(req.URL.Path, "/")
	if len(splitPath) > metricNamePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName := splitPath[3]

		switch memType {
		case gaugeMetricType:
			savedMetric, err := api.app.Metric(ctx, memName, memType)
			if err == nil {
				_, err = res.Write([]byte(strconv.FormatFloat(*savedMetric.Value, 'f', -1, 64)))
				if err != nil {
					fmt.Printf("Server: error writing request body %s\n", err)
				}
			} else {
				http.Error(res, "No metric", http.StatusNotFound)
				return
			}
		case counterMetricType:
			savedMetric, err := api.app.Metric(ctx, memName, memType)
			if err == nil {
				_, err = res.Write([]byte(fmt.Sprintf("%d", *savedMetric.Delta)))
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
	ctx := req.Context()

	if len(splitPath) > metricValuePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName := splitPath[3]
		// значение метрики
		memValue := splitPath[4]

		var err error
		var metric models.Metric

		switch memType {
		case gaugeMetricType:
			if len(splitPath) > metricNamePosition {
				metric.ID = memName
				metric.MType = gaugeMetricType
				metricValue, err := strconv.ParseFloat(memValue, 64)
				if err != nil {
					http.Error(res, "Wrong metric value", http.StatusBadRequest)
					return
				}
				metric.Value = &metricValue
				err = api.app.AddMetric(ctx, metric)
				if err != nil {
					http.Error(res, "Server: add metric error", http.StatusBadRequest)
					return
				}
			}
			_, err = res.Write([]byte(fmt.Sprintf("%f", *metric.Value)))
			if err != nil {
				fmt.Printf("Server: error writing request body %s\n", err)
			}
		case counterMetricType:
			metric.ID = memName
			metric.MType = gaugeMetricType
			counterValue, err := strconv.ParseInt(memValue, 10, 64)
			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
			}
			metric.Delta = &counterValue

			savedMetric, err := api.app.Metric(ctx, metric.ID, metric.MType)
			if err == nil {
				var counter int64

				if savedMetric.Delta == nil {
					counter = 0
				} else {
					counter = *savedMetric.Delta
				}

				if metric.Delta != nil {
					counter += *metric.Delta
					metric.Delta = &counter
				}
			}
			err = api.app.AddMetric(ctx, metric)
			if err != nil {
				http.Error(res, "Server: add metric error", http.StatusBadRequest)
				return
			}

			_, err = res.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
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

	ctx := req.Context()
	api.saveJSONMetric(ctx, res, metric)
}

func (api *httpAPI) pingDB(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	err := api.app.PingDB()
	if err != nil {
		http.Error(res, "Server: no DB connection", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (api *httpAPI) postMetrics(res http.ResponseWriter, req *http.Request) {
	var err error

	// читаем тело запроса
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	var metrics []models.Metric
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	for _, metric := range metrics {
		api.saveJSONMetric(ctx, res, metric)
	}
}

func (api *httpAPI) saveJSONMetric(ctx context.Context, res http.ResponseWriter, metric models.Metric) {
	var err error

	switch metric.MType {
	case gaugeMetricType:
		err = api.app.AddMetric(ctx, metric)
		if err != nil {
			http.Error(res, "Server: add metric error", http.StatusBadRequest)
			return
		}

	case counterMetricType:
		savedMetric, err := api.app.Metric(ctx, metric.ID, metric.MType)
		if err == nil {
			var counter int64

			// счмиываем значение счетчика метрики
			if savedMetric.Delta == nil {
				counter = 0
			} else {
				counter = *savedMetric.Delta
			}

			if metric.Delta != nil {
				counter += *metric.Delta
				metric.Delta = &counter
			}
		}
		err = api.app.AddMetric(ctx, metric)
		if err != nil {
			http.Error(res, "Server: add metric error", http.StatusBadRequest)
			return

		}
	default:
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	api.saveMetricsSync()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		fmt.Printf("Server: error writing request body %s\n", err)
	}
}
