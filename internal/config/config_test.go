package config

import (
	"reflect"
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	tests := []struct {
		name string
		want *VersionInfo
	}{
		{
			name: "positive test #1",
			want: &VersionInfo{
				BuildVersion: "N/A",
				BuildCommit:  "N/A",
				BuildDate:    "N/A",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVersionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
