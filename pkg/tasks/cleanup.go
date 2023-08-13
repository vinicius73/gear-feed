package tasks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
)

var _ Task[model.IEntry] = (*Cleanup[model.IEntry])(nil)

type Cleanup[T model.IEntry] struct {
	Notify bool `fig:"notify" yaml:"notify"`
}

func (t Cleanup[T]) Name() string {
	return "cleanup"
}

func (t Cleanup[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	logger := zerolog.Ctx(ctx)

	count, err := opts.Storage.Cleanup()
	if err != nil {
		return err
	}

	logger.Info().Int64("count", count).Msg("cleanup done")

	if opts.Sender == nil && !t.Notify {
		return nil
	}

	return t.notify(ctx, count, opts)
}

func (t Cleanup[T]) notify(ctx context.Context, count int64, opts TaskRunOptions[T]) error {
	return opts.Sender.SendCleanupNotify(ctx, sender.SendCleanupNotifyOptions{
		Count: count,
	})
}
