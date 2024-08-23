package httplayer

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
)

func Example_getMetricNoJSON() {
	api := &httpAPI{
		router: chi.NewRouter(),
		app:    applayer.New(storelayer.New()),
	}

	http.HandleFunc("/value/{metricType}/{metricName}", api.getMetricNoJSON)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
