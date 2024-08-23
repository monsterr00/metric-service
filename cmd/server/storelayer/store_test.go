package storelayer_test

import (
	"context"
	"errors"
	"flag"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
)

var (
	value  float64 = 0.5634
	value2 int64   = 25
)

func TestPing(t *testing.T) {
	type want struct {
		err     error
		startDB bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				err:     errors.New("db: not started"),
				startDB: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag.Parse()
			util.SetFlags()

			config.SetMode(config.FileMode)

			if test.want.startDB {
				config.SetMode(config.DBMode)
			}

			err := storelayer.New().Ping()
			if err != nil && err.Error() != test.want.err.Error() {
				t.Errorf("Ping return error %s, want %s", err, test.want.err)
			}
		})
	}
}

func Test_store_Create(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		ctx    context.Context
		metric models.Metric
	}{
		{
			name: "negative test #1",
			err:  nil,
			ctx:  context.Background(),
			metric: models.Metric{
				ID:    "PauseTotalNs",
				MType: "gauge",
				Delta: &value2,
				Value: &value,
			},
		},
	}
	flag.Parse()
	util.SetFlags()

	storl := storelayer.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storl.Create(tt.ctx, tt.metric)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Ping return error %s, want %s", err, tt.err)
				t.Errorf("store.Create() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

/*
	err := storelayer.New().Ping()
	if err != nil && err.Error() != test.want.err.Error() {
		t.Errorf("Ping return error %s, want %s", err, test.want.err)
		t.Errorf("store.Create() error = %v, wantErr %v", err, tt.wantErr)
	}
*/
