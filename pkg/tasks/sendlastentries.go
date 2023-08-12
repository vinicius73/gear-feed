package tasks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/linkloader/news"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sources"
)

var _ Task[model.IEntry] = (*SendLastEntries[model.IEntry])(nil)

const defaultSendLastEntriesLimit = 10

type SendLastEntries[T model.IEntry] struct {
	Limit   int                 `fig:"limit"   yaml:"limit"`
	Sources sources.LoadOptions `fig:"sources" yaml:"sources"`
}

func (t SendLastEntries[T]) Name() string {
	return "send_last_entries"
}

func (t SendLastEntries[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	definitions, err := sources.Load(ctx, t.Sources)
	if err != nil {
		return err
	}

	limit := t.Limit

	if limit == 0 {
		limit = defaultSendLastEntriesLimit
	} else if limit < 0 {
		zerolog.Ctx(ctx).Warn().Msg("defining limit to 0")
		limit = 0
	}

	entries, err := news.LoadEntries(ctx, news.LoadOptions[T]{
		LoadOptions: linkloader.LoadOptions{
			Sources: definitions,
			Workers: 0, // dynamic
		},
		Limit:   limit,
		Storage: opts.Storage,
	})
	if err != nil {
		return err
	}

	return opts.Sender.SendCollection(ctx, entries)
}
