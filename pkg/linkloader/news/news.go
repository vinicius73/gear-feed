package news

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

type LoadOptions[T model.IEntry] struct {
	LoadOptions linkloader.LoadOptions
	Storage     storage.Storage[T]
	Limit       int
}

func LoadEntries[T model.IEntry](ctx context.Context, opt LoadOptions[T]) ([]T, error) {
	logger := zerolog.Ctx(ctx)

	entries, err := linkloader.LoadEntries[T](ctx, opt.LoadOptions)
	if err != nil {
		return []T{}, err
	}

	logger.Info().Int("entries", len(entries)).Msg("loaded entries")

	if opt.Storage != nil {
		where := storage.Where(storage.WhereIs(storage.StatusNew), storage.WhereAllowMissed(true))
		entries, err = opt.Storage.Where(where, entries)
		if err != nil {
			return []T{}, err
		}

		logger.Info().Int("entries", len(entries)).Msg("filtered entries")
	}

	if opt.Limit > 0 && len(entries) > opt.Limit {
		logger.Info().Int("limit", opt.Limit).Msg("limiting entries")

		entries = entries[:opt.Limit]
	}

	return entries, nil
}
