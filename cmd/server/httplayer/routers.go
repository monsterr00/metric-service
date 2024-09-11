package httplayer

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/helpers"
	"go.uber.org/zap"
)

type httpAPI struct {
	router      *chi.Mux
	app         applayer.App
	sugarLogger *zap.SugaredLogger
}

// New инициализирует http-сервер и другие службы приложения.
func New(appLayer applayer.App) *httpAPI {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("Server: logger start error")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Server: logger sync error")
		}
	}()

	api := &httpAPI{
		router:      chi.NewRouter(),
		app:         appLayer,
		sugarLogger: logger.Sugar(),
	}

	api.setupRoutes()
	api.loadMetrics()
	go api.saveMetricsInterval()
	return api
}

// setupRoutes формирует связь путей и хэндлеров.
func (api *httpAPI) setupRoutes() {
	api.router.Get("/", api.WithLogging(api.GzipMiddleware(api.getMainPage)))
	api.router.Post("/update/{metricType}/{metricName}/{metricValue}", api.WithLogging(api.GzipMiddleware(api.postMetricNoJSON)))
	api.router.Post("/update/", api.WithLogging(api.GzipMiddleware(api.postMetric)))
	api.router.Post("/updates/", api.WithLogging(api.GzipMiddleware(api.postMetrics)))
	api.router.Get("/value/{metricType}/{metricName}", api.WithLogging(api.GzipMiddleware(api.getMetricNoJSON)))
	api.router.Post("/value/", api.WithLogging(api.GzipMiddleware(api.getMetric)))
	api.router.Get("/ping", api.WithLogging(api.GzipMiddleware(api.pingDB)))

	api.router.Get("/debug/pprof/", pprof.Index)
	api.router.Get("/debug/pprof/heap", pprof.Index)
	api.router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	api.router.Get("/debug/pprof/profile", pprof.Profile)
	api.router.Get("/debug/pprof/symbol", pprof.Symbol)
	api.router.Get("/debug/pprof/trace", pprof.Trace)
}

// Engage запускает http-сервер и другие службы приложения.
func (api *httpAPI) Engage() {
	helpers.PrintBuildInfo()

	err := http.ListenAndServe(config.ServerOptions.Host, api.router)
	if err != nil {
		log.Fatal(err)
	}
	api.saveMetrics()
	api.closeDB()
}

// saveMetrics сохраняет данные из мапы metrics в файл.
func (api *httpAPI) saveMetrics() {
	if config.ServerOptions.Mode == config.FileMode {
		err := api.app.SaveMetricsFile()
		if err != nil {
			log.Printf("Server: metrics file save on exit error, %s", err)
		}
	}
}

// saveMetricsInterval сохраняет данные из мапы metrics в файл с определенным интервалом сохранения.
func (api *httpAPI) saveMetricsInterval() {
	if config.ServerOptions.Mode == config.FileMode && config.ServerOptions.StoreInterval > 0 {
		for {
			err := api.app.SaveMetricsFile()
			if err != nil {
				log.Printf("Server: metrics file save interval error, %s", err)
			}
			time.Sleep(time.Duration(config.ServerOptions.StoreInterval) * time.Second)
		}
	}
}

// saveMetricsSync сохраняет данные из мапы metrics в файл в момент вызова функции, если не определен интервал сохранения.
func (api *httpAPI) saveMetricsSync() {
	if config.ServerOptions.Mode == config.FileMode && config.ServerOptions.StoreInterval == 0 {
		err := api.app.SaveMetricsFile()
		if err != nil {
			log.Printf("Server: metrics file save sync error, %s", err)
		}
	}
}

// loadMetrics загружает данные из БД или файла в мапу metrics.
func (api *httpAPI) loadMetrics() {
	if config.ServerOptions.Mode == config.FileMode && config.ServerOptions.Restore {
		err := api.app.LoadMetricsFile()
		if err != nil {
			log.Printf("Server: metrics file load error, %s", err)
		}
	}

	if config.ServerOptions.Mode == config.DBMode {
		err := api.app.LoadMetricsDB()
		if err != nil {
			log.Printf("Server: metrics file load error, %s", err)
		}
	}
}

// closeDB вызывает функцию закрытия БД.
func (api *httpAPI) closeDB() {
	err := api.app.CloseDB()
	if err != nil {
		log.Printf("Server: error closing db, %s", err)
	}
}
