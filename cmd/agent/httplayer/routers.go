package httplayer

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

const (
	gaugeMetricType   = "gauge"
	counterMetricType = "counter"
)

type httpAPI struct {
	client      *resty.Client
	app         applayer.App
	workersPool *Pool
}

func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		client:      resty.New(),
		app:         appLayer,
		workersPool: NewPool(),
	}

	api.setupClient()
	return api
}

func (api *httpAPI) setupClient() {
	api.client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second)
}

func (api *httpAPI) Engage() {
	// запускаем сбор метрик
	go api.app.SetMetrics()

	// подготовка запросов к отправке
	go api.prepBatch()

	//api.sendToServer()

	// отправка запросов
	api.workersPool.Run()
}

/*
func (api *httpAPI) sendToServer() {
	metrics, err := api.app.Metrics()
	if err != nil {
		log.Printf("Client: error getting gauge metrics %s\n", err)
	}

	requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, "/update/")
	var originalBody string

	for _, v := range metrics {
		switch v.MType {
		case gaugeMetricType:
			originalBody = fmt.Sprintf(`{"id":"%s","type":"%s","value":%f}`, v.ID, v.MType, *v.Value)
		case counterMetricType:
			originalBody = fmt.Sprintf(`{"id":"%s","type":"%s","delta":%d}`, v.ID, v.MType, *v.Delta)
		}

		api.sendReq(originalBody, requestURL)

		time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
	}
}
*/

func (api *httpAPI) compress(body string) (string, error) {
	var err error
	var buf bytes.Buffer
	b := []byte(body)

	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if _, err = gz.Write(b); err != nil {
		log.Printf("Client: compress error: %s\n", err)
		return buf.String(), err
	}
	if err = gz.Close(); err != nil {
		log.Printf("Client: gzip close error: %s\n", err)
		return buf.String(), err
	}

	return buf.String(), nil
}

func (api *httpAPI) prepBatch() {
	for {
		api.app.LockRW()
		metrics, err := api.app.Metrics()
		api.app.UnlockRW()

		if err != nil {
			log.Printf("Client: error getting gauge metrics %s\n", err)
		}

		var body = "["
		var counter int64

		for _, v := range metrics {
			switch v.MType {
			case gaugeMetricType:
				body += fmt.Sprintf(`{"id":"%s","type":"%s","value":%f},`, v.ID, v.MType, *v.Value)
				counter += 1
			case counterMetricType:
				body += fmt.Sprintf(`{"id":"%s","type":"%s","delta":%d},`, v.ID, v.MType, *v.Delta)
				counter += 1
			}

			if counter == config.ClientOptions.BatchSize {
				// отправляем запрос
				if len(body) > 1 {
					originalBody := body[:len(body)-1]
					originalBody += "]"
					api.sendReqToChan(originalBody)
					counter = 0
					body = "["
				}
			}
		}

		//отправляем остатки
		if len(body) > 1 {
			originalBody := body[:len(body)-1]
			originalBody += "]"
			api.sendReqToChan(originalBody)
		}

		time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
	}
}

func (api *httpAPI) sendReqToChan(originalBody string) {
	compressedBody, err := api.compress(originalBody)
	if err != nil {
		log.Printf("Client: compress error: %s\n", err)
	}
	req := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("HashSHA256", api.signBody(originalBody)).
		SetBody(compressedBody)

	api.workersPool.Add(req)
}

func (api *httpAPI) signBody(body string) string {
	// подписываем алгоритмом HMAC, используя SHA-256
	if config.ClientOptions.SignMode {
		h := hmac.New(sha256.New, []byte(config.ClientOptions.Key))
		h.Write([]byte(body))
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}
