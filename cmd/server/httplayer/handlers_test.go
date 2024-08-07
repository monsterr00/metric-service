package httplayer

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	value  float64 = 0.5634
	value2 int64   = 25
)

func TestGetMainPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/html",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			api := &httpAPI{
				router: chi.NewRouter(),
				app:    applayer.New(storelayer.New()),
			}

			api.getMainPage(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGetMetric(t *testing.T) {
	type want struct {
		code        int
		contentType string
		metric      models.Metric
	}
	tests := []struct {
		name     string
		want     want
		request  string
		post     string
		metricID string
		metrics  map[string]models.Metric
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "application/json",
				metric: models.Metric{
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: nil,
					Value: &value,
				},
			},
			request:  "/value/",
			post:     "/update/gauge/PauseTotalNs/0.5634",
			metricID: "PauseTotalNs",
			metrics: map[string]models.Metric{
				"PauseTotalNs": {
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: nil,
					Value: nil,
				},
			},
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
				metric: models.Metric{
					ID:    "PauseTotalNs",
					MType: "counter",
					Delta: nil,
					Value: &value,
				},
			},
			request:  "/value/",
			post:     "/update/counter/PauseTotalNs/1",
			metricID: "PauseTotalNs",
			metrics: map[string]models.Metric{
				"PauseTotalNs": {
					ID:    "PauseTotalNs",
					MType: "counterTT",
					Delta: nil,
					Value: nil,
				},
			},
		},
		{
			name: "negative test #2",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
				metric: models.Metric{
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: nil,
					Value: &value,
				},
			},
			request:  "/value/",
			post:     "/update/gauge/PauseTotalNs/1.1",
			metricID: "PauseTotalNs",
			metrics: map[string]models.Metric{
				"PauseTotalNs": {
					ID:    "PauseTotalNsTT",
					MType: "gauge",
					Delta: nil,
					Value: nil,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := &httpAPI{
				router: chi.NewRouter(),
				app:    applayer.New(storelayer.New()),
			}

			r := httptest.NewRequest(http.MethodPost, test.post, nil)
			w := httptest.NewRecorder()

			api.postMetricNoJSON(w, r)

			res := w.Result()
			if res.StatusCode != http.StatusOK {
				t.Error("error in post request")
			}
			res.Body.Close()

			metric := test.metrics[test.metricID]
			body := fmt.Sprintf(`{"id":"%s","type":"%s"}`, metric.ID, metric.MType)
			r = httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(body))
			w = httptest.NewRecorder()

			api.getMetric(w, r)
			res = w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if res.StatusCode == http.StatusOK {
				err = json.Unmarshal(resBody, &metric)
				require.NoError(t, err)
				assert.Equal(t, test.want.metric, metric)
			}
		})
	}
}
func TestGetMetricNoJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		want    want
		request string
		post    string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    "0.296303",
				contentType: "text/plain",
			},
			request: "/value/gauge/LastGC",
			post:    "/update/gauge/LastGC/0.296303",
		},
		{
			name: "positive test #2",
			want: want{
				code:        200,
				response:    "0.13",
				contentType: "text/plain",
			},
			request: "/value/gauge/LastGC",
			post:    "/update/gauge/LastGC/0.13",
		},
		{
			name: "positive test #3",
			want: want{
				code:        200,
				response:    "8",
				contentType: "text/plain",
			},
			request: "/value/counter/PollCount",
			post:    "/update/counter/PollCount/8",
		},
		{
			name: "negative test #1",
			want: want{
				code:        404,
				response:    "No metric\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/value/gauge/LastGCCVV",
			post:    "/update/gauge/LastGC/0.13",
		},
		{
			name: "negative test #2",
			want: want{
				code:        404,
				response:    "No metric\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/value/counter/PollCount1",
			post:    "/update/counter/PollCount/13",
		},
		{
			name: "negative test #3",
			want: want{
				code:        400,
				response:    "Wrong metric type\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/value/newType/PollCount1",
			post:    "/update/counter/PollCount/13",
		},
		{
			name: "negative test #4",
			want: want{
				code:        404,
				response:    "No metric type\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/value/newType",
			post:    "/update/counter/PollCount/13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := &httpAPI{
				router: chi.NewRouter(),
				app:    applayer.New(storelayer.New()),
			}

			r := httptest.NewRequest(http.MethodPost, test.post, nil)
			w := httptest.NewRecorder()

			api.postMetricNoJSON(w, r)

			res := w.Result()
			if res.StatusCode != http.StatusOK {
				t.Error("error in post request")
			}
			res.Body.Close()

			r = httptest.NewRequest(http.MethodGet, test.request, nil)
			w = httptest.NewRecorder()

			api.getMetricNoJSON(w, r)
			res = w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
func TestPostMetricNoJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    "0.296303",
				contentType: "text/plain",
			},
			request: "/update/gauge/NumForcedGC/0.296303",
		},
		{
			name: "positive test #2",
			want: want{
				code:        200,
				response:    "5",
				contentType: "text/plain",
			},
			request: "/update/counter/PollCount/5",
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/counter/PollCount/5.6",
		},
		{
			name: "negative test #2",
			want: want{
				code:        400,
				response:    "Wrong metric type\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/counterQQ/PollCount/5.6",
		},
		{
			name: "negative test #3",
			want: want{
				code:        404,
				response:    "No metric type or metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/counter",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := &httpAPI{
				router: chi.NewRouter(),
				app:    applayer.New(storelayer.New()),
			}

			r := httptest.NewRequest(http.MethodPost, test.request, nil)
			w := httptest.NewRecorder()

			api.postMetricNoJSON(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
func TestPostMetric(t *testing.T) {
	type want struct {
		code        int
		contentType string
		metric      models.Metric
	}
	tests := []struct {
		name     string
		want     want
		request  string
		metricID string
		metrics  map[string]models.Metric
		flag     bool
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "application/json",
				metric: models.Metric{
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: &value2,
					Value: &value,
				},
			},
			request:  "/update/",
			metricID: "PauseTotalNs",
			metrics: map[string]models.Metric{
				"PauseTotalNs": {
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: &value2,
					Value: &value,
				},
			},
			flag: true,
		},
		{
			name: "positive test #2",
			want: want{
				code:        200,
				contentType: "application/json",
				metric: models.Metric{
					ID:    "PollCount",
					MType: "counter",
					Delta: &value2,
					Value: &value,
				},
			},
			request:  "/update/",
			metricID: "PollCount",
			metrics: map[string]models.Metric{
				"PollCount": {
					ID:    "PollCount",
					MType: "counter",
					Delta: &value2,
					Value: &value,
				},
			},
			flag: false,
		},
		{
			name: "positive test #3",
			want: want{
				code:        200,
				contentType: "application/json",
				metric: models.Metric{
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: &value2,
					Value: &value,
				},
			},
			request:  "/update/",
			metricID: "PauseTotalNs",
			metrics: map[string]models.Metric{
				"PauseTotalNs": {
					ID:    "PauseTotalNs",
					MType: "gauge",
					Delta: &value2,
					Value: &value,
				},
			},
			flag: true,
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
				metric: models.Metric{
					ID:    "PollCount",
					MType: "counter",
					Delta: &value2,
					Value: &value,
				},
			},
			request:  "/update/",
			metricID: "PollCount",
			metrics: map[string]models.Metric{
				"PollCount": {
					ID:    "PollCount",
					MType: "counterTT",
					Delta: &value2,
					Value: &value,
				},
			},
			flag: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag.Parse()
			util.SetFlags()

			config.SetMode(config.FileMode)

			if test.flag {
				config.SetMode(config.DBMode)
			}
			api := &httpAPI{
				router: chi.NewRouter(),
				app:    applayer.New(storelayer.New()),
			}
			api.loadMetrics()

			metric := test.metrics[test.metricID]
			body := fmt.Sprintf(`{"id":"%s","type":"%s","delta":%d,"value":%f}`, metric.ID, metric.MType, *metric.Delta, *metric.Value)
			r := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(body))
			w := httptest.NewRecorder()

			api.postMetric(w, r)
			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if res.StatusCode == http.StatusOK {
				err = json.Unmarshal(resBody, &metric)
				require.NoError(t, err)
				assert.Equal(t, test.want.metric, metric)
			}
		})
	}
}

func BenchmarkCheckSign(b *testing.B) {
	b.StopTimer()
	storeLayer := storelayer.New()
	appLayer := applayer.New(storeLayer)
	apiLayer := New(appLayer)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		apiLayer.checkSign([]byte("dqwd1e1e32ed"), "secret key")
	}
}
