package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(ctx context.Context) (context.Context, context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
	case <-ctx.Done():
	}

	signal.Stop(quit)
	close(quit)

	return context.WithTimeout(ctx, 5*time.Second)
}
