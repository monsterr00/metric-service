package httplayer

import (
	"bytes"
	"compress/gzip"
	"context"
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

// New инициализирует уровень app
func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		client:      resty.New(),
		app:         appLayer,
		workersPool: NewPool(),
	}

	api.setupClient()
	return api
}

// setupClient устанавливает настройки http-клиента.
func (api *httpAPI) setupClient() {
	api.client.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)
}

// Engage запускает сбор метрик и другие службы приложения.
func (api *httpAPI) Engage() {
	go api.app.SetMetrics()
	go api.app.SetMetricsGOPSUTIL()

	go api.prepBatch()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api.workersPool.Run(ctx)
	api.workersPool.Stop()
}

// compress сжимает тело запроса.
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

// prepBatch разбивает массив отправляемых данных по метрикам на пакеты.
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

// sendReqToChan подготовалливает post-запрос и отправляет его в фабрику.
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

// signBody подписывает тело запроса алгоритмом HMAC, используя SHA-256.
func (api *httpAPI) signBody(body string) string {
	if config.ClientOptions.SignMode {
		h := hmac.New(sha256.New, []byte(config.ClientOptions.Key))
		h.Write([]byte(body))
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}
