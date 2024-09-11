package applayer

import (
	"testing"

	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

var (
	value  float64 = 0.5634
	value2 int64   = 25
)

func Test_app_GenMetricsStats(t *testing.T) {
	type fields struct {
		metrics map[string]models.Metric
		store   storelayer.Store
	}
	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "positive test #1",
			fields: fields{
				metrics: map[string]models.Metric{
					"PollCount": {
						ID:    "PollCount",
						MType: "counterTT",
						Delta: &value2,
						Value: &value,
					},
				},
				store: storelayer.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &app{
				metrics: tt.fields.metrics,
				store:   tt.fields.store,
			}
			api.GenMetricsStats()
		})
	}
}
