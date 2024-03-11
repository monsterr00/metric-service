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

		// тип метрики
		memType := splitPath[2]
		// имя метрики
		memName := splitPath[3]
		// значение метрики
		memValue := splitPath[4]

		switch memType {
		case "gauge":
			var err error
			memStorage.Gauge[memName], err = strconv.ParseFloat(memValue, 64)

			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
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

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		fmt.Printf("client: status code: %d\n", memStorage.Counter[memName])
		fmt.Printf("%s:%f\n", memName, memStorage.Gauge[memName])
		return
	} else {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
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

/*
ADDRESS отвечает за адрес эндпоинта HTTP-сервера.

export FILES=test1.txt:test2.txt
export TASK_DURATION=5s

Приоритет параметров должен быть таким:
Если указана переменная окружения, то используется она.
Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

Задание по треку «Сервис сбора метрик и алертинга»
+ Редиректы не поддерживаются.
+ Сервер должен быть доступен по адресу http://localhost:8080
+ Принимать метрики по протоколу HTTP методом POST.
+ Для хранения метрик объявите тип MemStorage. Рекомендуем использовать тип struct с полем-коллекцией внутри (slice или map). В будущем это позволит добавлять к объекту хранилища новые поля, например логер или мьютекс, чтобы можно было использовать их в методах.
+ Принимать данные в формате http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain.
+ При попытке передать запрос с некорректным типом метрики возвращать http.StatusBadRequest.
+ Принимать и хранить произвольные метрики двух типов:
+ Тип gauge, float64 — новое значение должно замещать предыдущее.
+ Тип counter, int64 — новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
+ При успешном приёме возвращать http.StatusOK.
+ При попытке передать запрос с некорректным значением возвращать http.StatusBadRequest.

- При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
- Опишите интерфейс для взаимодействия с этим хранилищем.
- Во всех заданиях нужно обрабатывать все пограничные случаи и негатив-кейсы. С каждым инкрементом автотесты будут становиться строже. Темплейты могут обновляться — с добавлением всё более строгих проверок.

Пример запроса к серверу:
POST /update/counter/someMetric/527 HTTP/1.1
Host: localhost:8080
Content-Length: 0
Content-Type: text/plain

Пример ответа от сервера:
HTTP/1.1 200 OK
Date: Tue, 21 Feb 2023 02:51:35 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

*/

// Мапа из мап
