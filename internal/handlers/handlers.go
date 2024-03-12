package handlers

import (
	"fmt"

	"net/http"
	"strconv"
	"strings"

	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

func MainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		var body string

		for k, v := range config.MemStorage.Gauge {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}

		for k, v := range config.MemStorage.Counter {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(body))
	} else {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func GetMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		splitPath := strings.Split(req.URL.Path, "/")
		if len(splitPath) > 3 {
			// тип метрики
			memType := splitPath[2]
			// имя метрики
			memName := splitPath[3]

			switch memType {
			case "gauge":
				memValue, isSet := config.MemStorage.Gauge[memName]
				if isSet {
					res.Write([]byte(fmt.Sprintf("%.3f", memValue)))
				} else {
					http.Error(res, "No metric", http.StatusNotFound)
					return
				}
			case "counter":
				memValue, isSet := config.MemStorage.Counter[memName]
				if isSet {
					res.Write([]byte(fmt.Sprintf("%d", memValue)))
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
	} else {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func PostMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		splitPath := strings.Split(req.URL.Path, "/")

		var memName string
		var memValue string

		if len(splitPath) > 4 {
			// тип метрики
			memType := splitPath[2]
			// имя метрики
			memName = splitPath[3]
			// значение метрики
			memValue = splitPath[4]

			switch memType {
			case "gauge":
				if len(splitPath) > 3 {
					var err error
					config.MemStorage.Gauge[memName], err = strconv.ParseFloat(memValue, 64)

					if err != nil {
						http.Error(res, "Wrong metric value", http.StatusBadRequest)
						return
					}
				}
				res.Write([]byte(fmt.Sprintf("%f", config.MemStorage.Gauge[memName])))
			case "counter":
				counterValue, err := strconv.ParseInt(memValue, 10, 64)
				if err != nil {
					http.Error(res, "Wrong metric value", http.StatusBadRequest)
					return
				}

				config.MemStorage.Counter[memName] += counterValue
				res.Write([]byte(fmt.Sprintf("%d", config.MemStorage.Counter[memName])))
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
		return
	} else {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}
