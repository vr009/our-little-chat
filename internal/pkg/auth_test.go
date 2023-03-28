package pkg

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	tt := time.Now()
	time.Sleep(1)
	t.Log(time.Until(tt))
}
