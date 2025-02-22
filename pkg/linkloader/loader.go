package linkloader

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

type LoadOptions struct {
	Workers int
	Sources []scraper.SourceDefinition
}

func FromSources[T model.IEntry](ctx context.Context, options LoadOptions) (Collections[T], error) {
	collections := Collections[T]{}

	chCollections := []<-chan Collection[T]{}
	chErrors := []<-chan error{}
	chSources := make(chan scraper.SourceDefinition, options.Workers)

	logger := zerolog.Ctx(ctx)

	var wg sync.WaitGroup

	for range options.Workers {
		wg.Add(1)
		out, errc := loadWorker[T](&wg, ctx, chSources)

		chCollections = append(chCollections, out)
		chErrors = append(chErrors, errc)
	}

	//nolint:gomnd
	wg.Add(2) // collections and errors

	go func() {
		defer wg.Done()

		for collection := range support.MergeChanners(chCollections...) {
			collections = append(collections, collection)
		}
	}()

	var err error

	go func() {
		defer wg.Done()

		for err := range support.MergeChanners(chErrors...) {
			logger.Error().Err(err).Msg("Error on load worker")
		}
	}()

	for _, source := range options.Sources {
		chSources <- source
	}

	close(chSources)

	wg.Wait()

	if err != nil {
		logger.Warn().Msg("There are errors on load workers")
	}

	return collections, nil
}

func FromSource[T model.IEntry](ctx context.Context, source scraper.SourceDefinition) (Collection[T], error) {
	entries, err := scraper.FindEntries[T](ctx, source)
	if err != nil {
		return Collection[T]{}, err
	}

	return Collection[T]{
		SourceName: source.Name,
		Entries:    entries,
	}, nil
}

//nolint:lll,revive
func loadWorker[T model.IEntry](wg *sync.WaitGroup, ctx context.Context, input <-chan scraper.SourceDefinition) (<-chan Collection[T], <-chan error) {
	//nolint:gomnd
	out := make(chan Collection[T], 2)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		defer wg.Done()

		for source := range input {
			collection, err := FromSource[T](ctx, source)
			if err != nil {
				errc <- err
			}

			out <- collection
		}
	}()

	return out, errc
}
