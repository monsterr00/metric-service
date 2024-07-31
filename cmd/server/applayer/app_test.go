package applayer

import (
	"testing"

	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
)

func BenchmarkSaveMetricsFile(b *testing.B) {
	b.StopTimer()
	storeLayer := storelayer.New()
	appLayer := New(storeLayer)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		appLayer.SaveMetricsFile()
	}
}
