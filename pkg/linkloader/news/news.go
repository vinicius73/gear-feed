package news

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

type LoadOptions struct {
	linkloader.LoadOptions
	Storage storage.Storage[model.Entry]
	Limit   int
}

func LoadEntries(ctx context.Context, opt LoadOptions) ([]model.Entry, error) {
	logger := zerolog.Ctx(ctx)

	entries, err := linkloader.LoadEntries(ctx, opt.LoadOptions)
	if err != nil {
		return []model.Entry{}, err
	}

	logger.Info().Int("entries", len(entries)).Msg("loaded entries")

	if opt.Storage != nil {
		entries, err := opt.Storage.Where(storage.WhereOptions{}, entries)
		if err != nil {
			return []model.Entry{}, err
		}

		logger.Info().Int("entries", len(entries)).Msg("filtered entries")
	}

	if opt.Limit > 0 && len(entries) > opt.Limit {
		logger.Info().Int("limit", opt.Limit).Msg("limiting entries")

		entries = entries[:opt.Limit]
	}

	return entries, nil
}
