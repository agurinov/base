package pipeline

import (
	"context"
	"testing"
)

func BenchmarkExecute(b *testing.B) {
	b.Run("success", func(b *testing.B) {
		execute(context.TODO(), &FakeLayer{})
	})
}
