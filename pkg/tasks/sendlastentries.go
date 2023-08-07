package tasks

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/linkloader/news"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/sources"
)

var _ Task[model.IEntry] = SendLastEntries[model.IEntry]{}

const defaultSendLastEntriesLimit = 10

type SendLastEntries[T model.IEntry] struct {
	Limit int      `fig:"limit" yaml:"limit"`
	Only  []string `fig:"only"  yaml:"only"`
}

func (t SendLastEntries[T]) Name() string {
	return "send_last_entries"
}

func (t SendLastEntries[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	definitions, err := sources.LoadDefinitions(ctx, sources.LoadOptions{
		Only: t.Only,
	})
	if err != nil {
		return err
	}

	entries, err := news.LoadEntries(ctx, news.LoadOptions[T]{
		LoadOptions: linkloader.LoadOptions{
			Sources: definitions,
			Workers: 0, // dynamic
		},
		Limit:   defaultSendLastEntriesLimit,
		Storage: opts.Storage,
	})
	if err != nil {
		return err
	}

	return opts.Sender.SendCollection(ctx, entries)
}
