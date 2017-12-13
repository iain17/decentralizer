package timeout

import (
	"time"
	"context"
)

func Do(op func(ctx context.Context), expire time.Duration) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(expire))
	go func() {
		op(ctx)
		cancel()
	}()
	<-ctx.Done()
}