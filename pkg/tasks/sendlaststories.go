package tasks

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/sources"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/stories"
)

type SendLastStories[T model.IEntry] struct {
	Limit    int                 `fig:"limit"    yaml:"limit"`
	Sources  sources.LoadOptions `fig:"sources"  yaml:"sources"`
	Interval time.Duration       `fig:"interval" yaml:"interval"`
}

func (t SendLastStories[T]) Name() string {
	return "send_last_stories"
}

func (t SendLastStories[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	logger := zerolog.Ctx(ctx).With().Str("component", t.Name()).Logger()

	entries, err := t.loadEntries(ctx, opts)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		logger.Warn().Msg("no entries to send")

		return nil
	}

	stories, removeAll, err := t.loadStories(ctx, opts, entries)
	if err != nil {
		return err
	}

	defer removeAll()

	for _, story := range stories {
		if err = opts.Sender.SendStory(ctx, story); err != nil {
			return err
		}
	}

	return nil
}

func (t SendLastStories[T]) loadEntries(ctx context.Context, opts TaskRunOptions[T]) ([]T, error) {
	var entries []T

	definitions, err := sources.Load(ctx, t.Sources)
	if err != nil {
		return entries, err
	}

	names := definitions.OnlyStorieSuported().Names()

	if len(names) == 0 {
		zerolog.Ctx(ctx).Warn().Msg("no sources to load")
		return entries, nil
	}

	return opts.Storage.FindByHasStory(storage.FindByHasStoryOptions{
		SourceNames: names,
		Interval:    t.Interval,
		Limit:       t.Limit,
		Has:         false,
	})
}

func (t SendLastStories[T]) loadStories(ctx context.Context, opts TaskRunOptions[T], entries []T) ([]sender.Story[T], func(), error) {
	var records []sender.Story[T]

	urls := make([]string, len(entries))

	for i, entry := range entries {
		urls[i] = entry.Link()
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "gamer-feed-stories-*")
	if err != nil {
		return records, func() {}, err
	}

	stories, err := stories.BuildCollection(ctx, stories.BuildCollectionOptions{
		Sources:          urls,
		TargetDir:        tmpDir,
		TemplateFilename: "{{.date}}-{{.site}}-{{.hash}}--{{.filename}}",
	})
	if err != nil {
		return records, func() {}, err
	}

	removeAll := func() {
		stories.RemoveAll()
	}

	hashMap := make(map[string]T)

	for _, entry := range entries {
		hash, err := entry.Hash()
		if err != nil {
			return records, removeAll, err
		}

		hashMap[hash] = entry
	}

	for _, story := range stories {
		entry := hashMap[story.Hash]

		record := sender.Story[T]{
			Story: story,
			Entry: entry,
		}

		records = append(records, record)
	}

	return records, removeAll, nil
}
