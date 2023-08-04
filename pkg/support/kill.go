//nolint:gomnd
package support

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WithKillSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-kill

		// logger := zerolog.Ctx(ctx)
		// logger.Warn().Msgf("OS Signal (%s)", sig.String())

		cancel()

		// Kill timeout
		<-time.After(time.Second * 20)
		// logger.Error().Msg("Stop timeout...")
		os.Exit(1)
	}()

	return ctx, cancel
}
