package httplayer

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type httpApi struct {
	client *resty.Client
	app    applayer.App
}

func New(appLayer applayer.App) *httpApi {
	api := &httpApi{
		client: resty.New(),
		app:    appLayer,
	}

	api.setupClient()
	return api
}

func (api *httpApi) setupClient() {
	api.client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second)
}

func (api *httpApi) Engage() {
	go api.app.SetMetrics()

	for {
		api.sendToServer()
		time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
	}

}

func (api *httpApi) sendToServer() {
	var err error
	gauge, err := api.app.GetGaugeMetrics()
	if err != nil {
		fmt.Printf("Client: error getting gauge metrics %s\n", err)
	}

	for k, v := range gauge {
		metricGaugeURL := fmt.Sprintf("/update/gauge/%s/%f", k, v)
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, metricGaugeURL)

		req, err := api.client.R().
			SetHeader("Content-Type", "text/plain").
			Post(requestURL)
		if err != nil {
			fmt.Printf("Client: error sending http-request: %s\n", err)
		}
		fmt.Printf("Status code: %d\n", req.StatusCode())
	}

	counter, err := api.app.GetCounterMetrics()
	if err != nil {
		fmt.Printf("Client: error getting counter metrics %s\n", err)
	}

	for k, v := range counter {
		metricCounterURL := fmt.Sprintf("/update/counter/%s/%d", k, v)
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, metricCounterURL)

		req, err := api.client.R().
			SetHeader("Content-Type", "text/plain").
			Post(requestURL)
		if err != nil {
			fmt.Printf("Client: error sending http-request: %s\n", err)
		}
		fmt.Printf("Status code: %d\n", req.StatusCode())
	}
}
