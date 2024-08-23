package storelayer_test

import (
	"errors"
	"flag"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
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
