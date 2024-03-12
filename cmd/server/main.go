package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

var memStorage MemStorage

var options struct {
	host string
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Для обновления параметров используйте ссылку в формате http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>", http.StatusBadRequest)
}

func updatePage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		splitPath := strings.Split(req.URL.Path, "/")

		var memName string
		var memValue string

		// проверяем возможное наличие типа метрики
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
					memStorage.Gauge[memName], err = strconv.ParseFloat(memValue, 64)

					if err != nil {
						http.Error(res, "Wrong metric value", http.StatusBadRequest)
						return
					}
				}

				res.Write([]byte(fmt.Sprintf("%f", memStorage.Gauge[memName])))
			case "counter":
				counterValue, err := strconv.ParseInt(memValue, 10, 64)

				if err != nil {
					http.Error(res, "Wrong metric value", http.StatusBadRequest)
					return
				}

				memStorage.Counter[memName] += counterValue

				res.Write([]byte(fmt.Sprintf("%d", memStorage.Counter[memName])))
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
		fmt.Printf("client: status code: %d\n", memStorage.Counter[memName])
		fmt.Printf("%s:%f\n", memName, memStorage.Gauge[memName])
		return
	} else {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func init() {
	flag.StringVar(&options.host, "a", "localhost:8080", "server host")
	envAddress, isEnv := os.LookupEnv("ADDRESS")

	if isEnv {
		options.host = envAddress
	}
}

func main() {
	memStorage.Gauge = make(map[string]float64)
	memStorage.Counter = make(map[string]int64)

	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/`, updatePage)

	err := http.ListenAndServe(options.host, mux)
	if err != nil {
		log.Fatal(err)
	}
}
