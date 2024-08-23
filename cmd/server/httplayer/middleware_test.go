package httplayer

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
	"go.uber.org/zap"
)

func TestWithLogging(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	flag.Parse()
	util.SetFlags()
	api := &httpAPI{
		router:      chi.NewRouter(),
		app:         applayer.New(storelayer.New()),
		sugarLogger: logger.Sugar(),
	}
	api.loadMetrics()

	r := httptest.NewRequest(http.MethodPost, "/update/gauge/PauseTotalNs/0.5634", nil)
	w := httptest.NewRecorder()

	api.postMetricNoJSON(w, r)

	wl := api.WithLogging(api.postMetricNoJSON)
	wl.ServeHTTP(w, r)
}

func TestGzipMiddleware(t *testing.T) {
	flag.Parse()
	util.SetFlags()
	api := &httpAPI{
		router: chi.NewRouter(),
		app:    applayer.New(storelayer.New()),
	}
	api.loadMetrics()

	r := httptest.NewRequest(http.MethodPost, "/update/gauge/PauseTotalNs/0.5634", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Content-Type", "application/json")
	api.postMetricNoJSON(w, r)

	wl := api.GzipMiddleware(api.postMetricNoJSON)
	wl.ServeHTTP(w, r)
}
