package httplayer

import (
	"fmt"

	"net/http"
	"strconv"
	"strings"
)

func (api *httpApi) getMainPage(res http.ResponseWriter, req *http.Request) {
	var body string

	gauge, err := api.app.GetGaugeMetrics()
	if err != nil {
		fmt.Printf("Client: error getting gauge metrics %s\n", err)
	}

	for k, v := range gauge {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}

	counter, err := api.app.GetCounterMetrics()
	if err != nil {
		fmt.Printf("Client: error getting counter metrics %s\n", err)
	}
	for k, v := range counter {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte(body))
	if err != nil {
		fmt.Printf("Server: error writing request body %s\n", err)
	}
}

func (api *httpApi) getMetric(res http.ResponseWriter, req *http.Request) {
	const gaugeMetricType = "gauge"
	const counterMetricType = "counter"
	const metricNamePosition = 3

	gauge, err := api.app.GetGaugeMetrics()
	if err != nil {
		fmt.Printf("Client: error getting gauge metrics %s\n", err)
	}
	counter, err := api.app.GetCounterMetrics()
	if err != nil {
		fmt.Printf("Client: error getting counter metrics %s\n", err)
	}

	splitPath := strings.Split(req.URL.Path, "/")
	if len(splitPath) > metricNamePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName := splitPath[3]

		switch memType {
		case gaugeMetricType:
			memValue, isSet := gauge[memName]
			if isSet {
				_, err = res.Write([]byte(fmt.Sprintf("%.3f", memValue)))
				if err != nil {
					fmt.Printf("Server: error writing request body %s\n", err)
				}
			} else {
				http.Error(res, "No metric", http.StatusNotFound)
				return
			}
		case counterMetricType:
			memValue, isSet := counter[memName]
			if isSet {
				_, err = res.Write([]byte(fmt.Sprintf("%d", memValue)))
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

func (api *httpApi) postMetric(res http.ResponseWriter, req *http.Request) {
	splitPath := strings.Split(req.URL.Path, "/")

	const metricNamePosition = 3
	const metricValuePosition = 4

	var memName string
	var memValue string

	gauge, err := api.app.GetGaugeMetrics()
	if err != nil {
		fmt.Printf("Client: error getting gauge metrics %s\n", err)
	}
	counter, err := api.app.GetCounterMetrics()
	if err != nil {
		fmt.Printf("Client: error getting counter metrics %s\n", err)
	}

	if len(splitPath) > metricValuePosition {
		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName = splitPath[3]
		// значение метрики
		memValue = splitPath[4]

		switch memType {
		case "gauge":
			if len(splitPath) > metricNamePosition {
				var err error
				gauge[memName], err = strconv.ParseFloat(memValue, 64)

				if err != nil {
					http.Error(res, "Wrong metric value", http.StatusBadRequest)
					return
				}
			}
			_, err = res.Write([]byte(fmt.Sprintf("%f", gauge[memName])))
			if err != nil {
				fmt.Printf("Server: error writing request body %s\n", err)
			}
		case "counter":
			counterValue, err := strconv.ParseInt(memValue, 10, 64)
			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
			}

			counter[memName] += counterValue
			_, err = res.Write([]byte(fmt.Sprintf("%d", counter[memName])))
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
