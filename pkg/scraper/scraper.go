package scraper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var ErrCategoryNotAllowed = errors.New("category not allowed")

const (
	requestTimeout = time.Second * 15
	titleLimit     = 150
)

func newCollector() *colly.Collector {
	c := colly.NewCollector()
	c.SetRequestTimeout(requestTimeout)

	return c
}

func FindEntries(ctx context.Context, source SourceDefinition) ([]Entry, error) {
	logger := zerolog.Ctx(ctx).With().Str("source", source.Name).Logger()

	ctx = logger.WithContext(ctx)

	entries := []Entry{}

	callback := func(e Element) {
		entry, err := onEntry(ctx, source, e)
		if err != nil {
			if !errors.Is(err, ErrCategoryNotAllowed) {
				logger.Error().Err(err).Msg("Error on entry")
			}

			return
		}

		logger.
			Debug().
			Msgf("New entry: %s", entry.Link)

		entries = append(entries, entry)
	}

	startTime := time.Now()

	err := visit(ctx, source, callback)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)

	logger.
		Info().
		Dur("duration", duration).
		Msgf("Done with %v entries (%s).", len(entries), duration.String())

	return entries, nil
}

func visit(ctx context.Context, source SourceDefinition, callback func(e Element)) error {
	logger := zerolog.Ctx(ctx)

	collector := newCollector()
	entrySelector := source.Attributes.EntrySelector

	parser := strings.ToUpper(source.Parser)

	if parser == "XML" {
		collector.OnXML(entrySelector, func(e *colly.XMLElement) {
			callback(e)
		})
	} else {
		collector.OnHTML(entrySelector, func(e *colly.HTMLElement) {
			callback(e)
		})
	}

	collector.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
	})

	collector.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))
	})

	err := collector.Visit(source.visitURL())
	if err != nil {
		logger.
			Error().
			Err(err).
			Msgf("Fail to visit %s", source.visitURL())

		return fmt.Errorf("error on visit: %w", err)
	}

	logger.Info().Msg("Fetching source")

	collector.Wait()

	return nil
}

func onEntry(ctx context.Context, source SourceDefinition, el Element) (Entry, error) {
	attributes := source.Attributes

	title := attributes.Title.findAttribute(el)

	categories := support.ToLower(attributes.Category.findCategories(el))

	if !attributes.Category.isAllowed(categories) {
		zerolog.Ctx(ctx).Debug().
			Strs("categories", categories).
			Str("title", title).
			Msg(ErrCategoryNotAllowed.Error())

		return Entry{}, ErrCategoryNotAllowed
	}

	link := attributes.Link.findAttribute(el)
	image := attributes.Image.findAttribute(el)

	if len(title) > titleLimit {
		title = title[:titleLimit]
	}

	return source.buildEntry(title, link, image, categories), nil
}
