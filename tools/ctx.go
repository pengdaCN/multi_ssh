package tools

import "context"

func WithCancel(ctx context.Context, fn func()) {
	f := make(chan struct{})
	go func() {
		fn()
		f <- struct{}{}
	}()
	select {
	case <-ctx.Done():
	case <-f:
	}
}
