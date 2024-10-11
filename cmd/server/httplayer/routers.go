package httplayer

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/helpers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type httpAPI struct {
	router      *chi.Mux
	app         applayer.App
	sugarLogger *zap.SugaredLogger
	grpc        *grpc.Server
	srvHTTP     http.Server
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
		grpc:        grpc.NewServer(grpc.UnaryInterceptor(checkSubnetInterceptor)),
	}

	if config.ServerOptions.GrpcOn {
		api.startGrpc()
	} else {
		api.setupRoutes()
	}

	api.loadMetrics()
	go api.saveMetricsInterval()
	return api
}

// setupRoutes формирует связь путей и хэндлеров.
func (api *httpAPI) setupRoutes() {
	api.router.Get("/", api.WithLogging(api.GzipMiddleware(api.getMainPage)))
	api.router.Post("/update/{metricType}/{metricName}/{metricValue}", api.WithLogging(api.GzipMiddleware(api.postMetricNoJSON)))
	api.router.Post("/update/", api.WithLogging(api.GzipMiddleware(api.postMetric)))
	api.router.Post("/updates/", api.WithLogging(api.Decrypt(api.GzipMiddleware(api.SubNetMiddleware(api.postMetrics)))))
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
	api.generateCryptoKeys()

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		api.stopServer()
		close(idleConnsClosed)
	}()

	api.startServer()

	<-idleConnsClosed
	log.Printf("Server Shutdown gracefully start")
	api.stopServices()
	log.Printf("Server Shutdown gracefully end")
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

// startServer запускает http/grpc-серверы
func (api *httpAPI) startServer() {
	if config.ServerOptions.GrpcOn {
		listen, err := net.Listen("tcp", config.ServerOptions.GrpcHost)
		if err != nil {
			log.Fatal(err)
		}

		if err := api.grpc.Serve(listen); err != nil {
			log.Fatal(err)
		}
	} else {
		api.srvHTTP = http.Server{Addr: config.ServerOptions.Host, Handler: api.router}

		if err := api.srvHTTP.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}
}

// stopServer останавливает http/grpc-серверы
func (api *httpAPI) stopServer() {
	if config.ServerOptions.GrpcOn {
		api.grpc.Stop()
	} else {
		if err := api.srvHTTP.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}
}

// stopServer останавливает сервисы
func (api *httpAPI) stopServices() {
	if config.ServerOptions.Mode == config.FileMode {
		api.saveMetrics()
	}

	if config.ServerOptions.Mode == config.DBMode {
		api.closeDB()
	}
}
