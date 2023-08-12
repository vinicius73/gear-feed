package news

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

type LoadOptions[T model.IEntry] struct {
	LoadOptions linkloader.LoadOptions
	Storage     storage.Storage[T]
	Limit       int
}

type SourceResultEntries[T model.IEntry] struct {
	SourceResult
	Entries []T
}

type SourceResult struct {
	Source   string
	Total    int
	Filtered int
}

type SourceResultEntriesList[T model.IEntry] []SourceResultEntries[T]

type Result[T model.IEntry] struct {
	Results  []SourceResult
	Loaded   int
	Filtered int
	Entries  []T
}

func LoadEntries[T model.IEntry](ctx context.Context, opt LoadOptions[T]) (Result[T], error) {
	logger := zerolog.Ctx(ctx)

	entries, err := linkloader.LoadEntries[T](ctx, opt.LoadOptions)
	if err != nil {
		return Result[T]{}, err
	}

	logger.Info().Int("entries", len(entries)).Msg("loaded entries")

	results, err := buildResult[T](ctx, opt, entries)
	if err != nil {
		return Result[T]{}, err
	}

	return results, nil
}

func buildResult[T model.IEntry](ctx context.Context, opt LoadOptions[T], loadedEntries []T) (Result[T], error) {
	logger := zerolog.Ctx(ctx)

	grouped := map[string][]T{}

	where := storage.WhereNotSent()

	for _, entry := range loadedEntries {
		grouped[entry.Source()] = append(grouped[entry.Source()], entry)
	}

	results := SourceResultEntriesList[T]{}

	for source, entries := range grouped {
		total := len(entries)
		entries, err := opt.Storage.Where(where, entries)
		if err != nil {
			return Result[T]{}, err
		}

		result := SourceResultEntries[T]{
			Entries: entries,
			SourceResult: SourceResult{
				Total:    total,
				Source:   source,
				Filtered: len(entries),
			},
		}

		logger.Info().
			Str("source", source).
			Int("total", total).
			Int("filtered", len(entries)).
			Msg("filtered entries")

		results = append(results, result)
	}

	entries := results.Limit(ctx, opt.Limit)

	return Result[T]{
		Entries:  entries,
		Loaded:   len(loadedEntries),
		Filtered: len(entries),
		Results:  results.SourceResults(),
	}, nil
}

func (r SourceResultEntriesList[T]) Limit(ctx context.Context, limit int) []T {
	entries := []T{}

	for _, result := range r {
		entries = append(entries, result.Entries...)
	}

	entries = support.Shuffle(entries)

	if limit > 0 && len(entries) > limit {
		zerolog.Ctx(ctx).Info().Int("limit", limit).Msg("limiting entries")
		entries = entries[:limit]
	}

	return entries
}

func (r SourceResultEntriesList[T]) SourceResults() []SourceResult {
	results := []SourceResult{}

	for _, result := range r {
		results = append(results, result.SourceResult)
	}

	return results
}

func (r Result[T]) Resume() sender.Resume {
	sources := make([]sender.ResumeSource, len(r.Results))

	for index, result := range r.Results {
		sources[index] = sender.ResumeSource{
			Source:   result.Source,
			Loaded:   result.Total,
			Filtered: result.Filtered,
		}
	}

	return sender.Resume{
		Loaded:   r.Loaded,
		Filtered: r.Filtered,
		Sources:  sources,
	}
}
