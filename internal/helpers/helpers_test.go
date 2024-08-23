package helpers

import (
	"testing"
)

func BenchmarkAbsolutePath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AbsolutePath("dfwefw", "wsefwef")
	}
}
