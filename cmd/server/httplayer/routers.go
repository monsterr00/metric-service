package httplayer

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type httpApi struct {
	router *chi.Mux
	app    applayer.App
}

func New(appLayer applayer.App) *httpApi {
	api := &httpApi{
		router: chi.NewRouter(),
		app:    appLayer,
	}

	api.setupRoutes()
	return api
}

func (api *httpApi) setupRoutes() {
	api.router.Get("/", api.getMainPage)
	api.router.Post("/update/{metricType}/{metricName}/{metricValue}", api.postMetric)
	api.router.Get("/value/{metricType}/{metricName}", api.getMetric)
}

func (api *httpApi) Engage() {
	err := http.ListenAndServe(config.ServerOptions.Host, api.router)
	if err != nil {
		log.Fatal(err)
	}
}
