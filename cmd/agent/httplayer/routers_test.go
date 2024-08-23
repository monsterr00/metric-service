package httplayer

import (
	"flag"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
)

func Test_httpAPI_signBody(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		want     string
		signMode bool
	}{
		{
			name:     "positive test #1",
			body:     "this body need to sign",
			want:     "ca1c23608c89ebac1a0bc0b65fb5495aa7f0bd1d2003391e3eb1585366053ac3",
			signMode: true,
		},
		{
			name:     "positive test #2",
			body:     "this body need to sign",
			want:     "",
			signMode: false,
		},
	}
	flag.Parse()
	util.SetFlags()

	api := &httpAPI{
		client:      resty.New(),
		app:         applayer.New(storelayer.New()),
		workersPool: NewPool(),
	}

	for _, tt := range tests {
		config.SetSignMode(tt.signMode)

		t.Run(tt.name, func(t *testing.T) {
			if got := api.signBody(tt.body); got != tt.want {
				t.Errorf("httpAPI.signBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpAPI_compress(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		body    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			want:    "\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x00\x10\x00\xef\xffbody to compress\x01\x00\x00\xff\xff68=4\x10\x00\x00\x00",
			body:    "body to compress",
			wantErr: false,
		},
	}
	flag.Parse()
	util.SetFlags()

	api := &httpAPI{
		client:      resty.New(),
		app:         applayer.New(storelayer.New()),
		workersPool: NewPool(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := api.compress(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpAPI.compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("httpAPI.compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpAPI_setupClient(t *testing.T) {
	type fields struct {
		client      *resty.Client
		app         applayer.App
		workersPool *Pool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				client:      resty.New(),
				app:         applayer.New(storelayer.New()),
				workersPool: NewPool(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &httpAPI{
				client:      tt.fields.client,
				app:         tt.fields.app,
				workersPool: tt.fields.workersPool,
			}
			api.setupClient()
		})
	}
}
