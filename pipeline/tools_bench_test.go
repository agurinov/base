package pipeline

import (
	"testing"
)

func BenchmarkExecute(b *testing.B) {
	b.Run("success", func(b *testing.B) {
		execute(&FakeLayer{})
	})
}
