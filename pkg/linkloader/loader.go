package linkloader

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/scraper"
	"github.com/vinicius73/gear-feed/pkg/support"
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

		for collection := range support.MergeChannels(chCollections...) {
			collections = append(collections, collection)
		}
	}()

	go func() {
		defer wg.Done()

		for err := range support.MergeChannels(chErrors...) {
			logger.Error().Err(err).Msg("Error on load worker")
		}
	}()

	for _, source := range options.Sources {
		chSources <- source
	}

	close(chSources)

	wg.Wait()

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

func loadWorker[T model.IEntry](wg *sync.WaitGroup, ctx context.Context, input <-chan scraper.SourceDefinition) (<-chan Collection[T], <-chan error) {
	out := make(chan Collection[T], 2)
	errc := make(chan error, 1)

	logger := zerolog.Ctx(ctx).With().Str("worker", "load").Logger()

	go func() {
		defer close(out)
		defer close(errc)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				logger.Warn().Msg("Context done")

				return
			case source, ok := <-input:
				if !ok {
					logger.Debug().Msg("Input closed")

					return
				}
				collection, err := FromSource[T](ctx, source)
				if err != nil {
					logger.Error().Err(err).Str("source", source.Name).Msg("Error on load worker")
					errc <- err
				} else {
					out <- collection
				}
			}
		}
	}()

	return out, errc
}
