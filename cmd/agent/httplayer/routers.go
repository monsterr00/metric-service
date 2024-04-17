package httplayer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type httpAPI struct {
	client *resty.Client
	app    applayer.App
}

func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		client: resty.New(),
		app:    appLayer,
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
	go api.app.SetMetrics()

	for {
		api.sendToServer()
		time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
	}
}

func (api *httpAPI) sendToServer() {
	var err error

	gauge, err := api.app.GetGaugeMetrics()
	if err != nil {
		log.Printf("Client: error getting gauge metrics %s\n", err)
	}

	requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, "/update/")

	for k, v := range gauge {
		originalBody := fmt.Sprintf(`{"id":"%s","type":"gauge","value":%f}`, k, v)
		compressedBody, err := compress(originalBody)
		if err != nil {
			log.Printf("Client: compress error: %s\n", err)
			continue
		}

		req, err := api.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(compressedBody).
			Post(requestURL)
		if err != nil {
			log.Printf("Client: error sending http-request: %s\n", err)
		}
		log.Printf("Before compress, %d, after compress, %d, status code: %d\n", len(originalBody), len(compressedBody), req.StatusCode())
	}

	counter, err := api.app.GetCounterMetrics()
	if err != nil {
		log.Printf("Client: error getting counter metrics %s\n", err)
	}

	for k, v := range counter {
		originalBody := fmt.Sprintf(`{"id":"%s","type":"counter","delta":%d}`, k, v)
		compressedBody, err := compress(originalBody)
		if err != nil {
			log.Printf("Client: compress error: %s\n", err)
			continue
		}

		req, err := api.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(compressedBody).
			Post(requestURL)
		if err != nil {
			log.Printf("Client: error sending http-request: %s\n", err)
		}
		log.Printf("Before compress, %d, after compress, %d, status code: %d\n", len(originalBody), len(compressedBody), req.StatusCode())
	}
}

func compress(body string) (string, error) {
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
