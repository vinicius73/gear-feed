package tasks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
)

var _ Task[model.IEntry] = (*Cleanup[model.IEntry])(nil)

type Cleanup[T model.IEntry] struct{}

func (t Cleanup[T]) Name() string {
	return "cleanup"
}

func (t Cleanup[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	logger := zerolog.Ctx(ctx)

	num, err := opts.Storage.Cleanup()
	if err != nil {
		return err
	}

	logger.Info().Int64("num", num).Msg("cleanup done")

	return nil
}
