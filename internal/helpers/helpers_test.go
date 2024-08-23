package helpers

import (
	"testing"
)

func BenchmarkAbsolutePath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AbsolutePath("dfwefw", "wsefwef")
	}
}

func TestPrintBuildInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintBuildInfo()
		})
	}
}

func TestAbsolutePath(t *testing.T) {
	tests := []struct {
		name      string
		pathStart string
		pathEnd   string
		want      string
	}{
		{
			name:      "positive test #1",
			pathStart: "file://",
			pathEnd:   "main.go",
			want:      "file://Users/denis/metric-service/main.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsolutePath(tt.pathStart, tt.pathEnd); got != tt.want {
				t.Errorf("AbsolutePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
