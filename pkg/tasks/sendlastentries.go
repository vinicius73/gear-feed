package tasks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/linkloader/news"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/sources"
)

var _ Task[model.IEntry] = (*SendLastEntries[model.IEntry])(nil)

const defaultSendLastEntriesLimit = 10

type SendLastEntries[T model.IEntry] struct {
	Limit        int                 `fig:"limit"          yaml:"limit"`
	SendResumeTo []int64             `fig:"send_resume_to" yaml:"send_resume_to"`
	Sources      sources.LoadOptions `fig:"sources"        yaml:"sources"`
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

	if len(entries.Entries) == 0 {
		zerolog.Ctx(ctx).Info().Msg("no entries to send")
	} else if err = opts.Sender.SendCollection(ctx, entries.Entries); err != nil {
		return err
	}

	if len(t.SendResumeTo) > 0 {
		if err = t.sendResume(ctx, entries, opts); err != nil {
			return err
		}
	}

	return nil
}

func (t SendLastEntries[T]) sendResume(ctx context.Context, entries news.Result[T], opts TaskRunOptions[T]) error {
	sources := make([]sender.ResumeSource, len(entries.Results))

	for index, result := range entries.Results {
		sources[index] = sender.ResumeSource{
			Source:   result.Source,
			Loaded:   result.Total,
			Filtered: result.Filtered,
		}
	}

	return opts.Sender.SendResume(ctx, sender.SendResumeOptions{
		Chats: t.SendResumeTo,
		Resume: sender.Resume{
			Loaded:   entries.Loaded,
			Filtered: entries.Filtered,
			Sources:  sources,
		},
	})
}
