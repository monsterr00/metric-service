package httplayer

import (
	"flag"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
	"go.uber.org/zap"
)

func Test_httpAPI_saveMetrics(t *testing.T) {
	flag.Parse()
	util.SetFlags()
	logger, _ := zap.NewDevelopment()

	type fields struct {
		router      *chi.Mux
		app         applayer.App
		sugarLogger *zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				router:      chi.NewRouter(),
				app:         applayer.New(storelayer.New()),
				sugarLogger: logger.Sugar(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &httpAPI{
				router:      tt.fields.router,
				app:         tt.fields.app,
				sugarLogger: tt.fields.sugarLogger,
			}
			api.saveMetrics()
		})
	}
}

func Test_httpAPI_setupRoutes(t *testing.T) {
	flag.Parse()
	util.SetFlags()
	logger, _ := zap.NewDevelopment()
	type fields struct {
		router      *chi.Mux
		app         applayer.App
		sugarLogger *zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				router:      chi.NewRouter(),
				app:         applayer.New(storelayer.New()),
				sugarLogger: logger.Sugar(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &httpAPI{
				router:      tt.fields.router,
				app:         tt.fields.app,
				sugarLogger: tt.fields.sugarLogger,
			}
			api.setupRoutes()
		})
	}
}
