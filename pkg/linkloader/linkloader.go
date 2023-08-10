package linkloader

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/model"
)

func LoadEntries[T model.IEntry](ctx context.Context, opt LoadOptions) ([]T, error) {
	if opt.Workers < 1 {
		//nolint:gomnd
		opt.Workers = (len(opt.Sources) + 1) / 2
	}

	collections, err := FromSources[T](ctx, opt)
	if err != nil {
		return []T{}, err
	}

	return collections.Shuffle(), nil
}
