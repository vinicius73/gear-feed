package linkloader

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/model"
)

func LoadEntries(ctx context.Context, opt LoadOptions) ([]model.Entry, error) {
	if opt.Workers < 1 {
		opt.Workers = (len(opt.Sources) + 1) / 2
	}

	collections, err := FromSources(ctx, opt)
	if err != nil {
		return []model.Entry{}, err
	}

	return collections.Shuffle(), nil
}
