package delivery

import (
	"github.com/google/uuid"
	"testing"
)

func BenchmarkAuthHandler_GetToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
}
