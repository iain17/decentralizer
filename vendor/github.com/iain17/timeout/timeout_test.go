package timeout

import (
	"testing"
	"context"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestDoExpire(t *testing.T) {
	now := time.Now()
	Do(func(ctx context.Context) {
		time.Sleep(12 * time.Second)
	}, 10 * time.Second)
	assert.True(t, time.Since(now).Seconds() < 12)
}

func TestDoNotExpire(t *testing.T) {
	now := time.Now()
	Do(func(ctx context.Context) {
		time.Sleep(3 * time.Second)
	}, 10 * time.Second)
	assert.True(t, time.Since(now).Seconds() < 5)
}
