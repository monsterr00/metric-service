package httplayer

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"go.uber.org/zap"
)

type httpAPI struct {
	router      *chi.Mux
	app         applayer.App
	sugarLogger *zap.SugaredLogger
}

func New(appLayer applayer.App) *httpAPI {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("Server: logger start error")
	}
	defer logger.Sync()

	api := &httpAPI{
		router:      chi.NewRouter(),
		app:         appLayer,
		sugarLogger: logger.Sugar(),
	}

	api.setupRoutes()
	api.loadMetricStart()
	go api.saveMetricsInterval()
	return api
}

func (api *httpAPI) setupRoutes() {
	api.router.Get("/", api.WithLogging(api.GzipMiddleware(api.getMainPage)))
	api.router.Post("/update/{metricType}/{metricName}/{metricValue}", api.WithLogging(api.postMetricNoJSON))
	api.router.Post("/update/", api.WithLogging(api.GzipMiddleware(api.postMetric)))
	api.router.Get("/value/{metricType}/{metricName}", api.WithLogging(api.getMetricNoJSON))
	api.router.Post("/value/", api.WithLogging(api.GzipMiddleware(api.getMetric)))
}

func (api *httpAPI) Engage() {
	err := http.ListenAndServe(config.ServerOptions.Host, api.router)
	if err != nil {
		log.Fatal(err)
	}
	api.saveMetricsExit()
}

func (api *httpAPI) saveMetricsExit() {
	if config.ServerOptions.FileStoragePath != "" {
		err := api.app.SaveMetrics()
		if err != nil {
			log.Printf("Server: metrics file save on exit error, %s", err)
		}
	}
}

func (api *httpAPI) saveMetricsInterval() {
	if config.ServerOptions.FileStoragePath != "" && config.ServerOptions.StoreInterval > 0 {
		for {
			err := api.app.SaveMetrics()
			if err != nil {
				log.Printf("Server: metrics file save interval error, %s", err)
			}
			time.Sleep(time.Duration(config.ServerOptions.StoreInterval) * time.Second)
		}
	}
}

func (api *httpAPI) saveMetricsSync() {
	if config.ServerOptions.FileStoragePath != "" && config.ServerOptions.StoreInterval == 0 {
		err := api.app.SaveMetrics()
		if err != nil {
			log.Printf("Server: metrics file save sync error, %s", err)
		}
	}
}

func (api *httpAPI) loadMetricStart() {
	if config.ServerOptions.FileStoragePath != "" && config.ServerOptions.Restore {
		err := api.app.LoadMetrics()
		if err != nil {
			log.Printf("Server: metrics file load error, %s", err)
		}
	}
}
