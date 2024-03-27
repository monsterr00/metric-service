package httplayer

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type httpAPI struct {
	router *chi.Mux
	app    applayer.App
}

func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		router: chi.NewRouter(),
		app:    appLayer,
	}

	api.setupRoutes()
	return api
}

func (api *httpAPI) setupRoutes() {
	api.router.Get("/", api.getMainPage)
	api.router.Post("/update/{metricType}/{metricName}/{metricValue}", api.postMetric)
	api.router.Get("/value/{metricType}/{metricName}", api.getMetric)
}

func (api *httpAPI) Engage() {
	err := http.ListenAndServe(config.ServerOptions.Host, api.router)
	if err != nil {
		log.Fatal(err)
	}
}
