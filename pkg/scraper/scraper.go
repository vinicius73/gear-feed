package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var ErrCategoryNotAllowed = errors.New("category not allowed")
var ErrFailToCrateRequest = errors.New("fail to create request")

const (
	requestTimeout = time.Second * 15
	titleLimit     = 150
	userAgent      = "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"
)

func newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(userAgent),
	)
	c.SetRequestTimeout(requestTimeout)

	return c
}

func FindEntries[T model.IEntry](ctx context.Context, source SourceDefinition) ([]T, error) {
	switch strings.ToUpper(source.Parser) {
	case "JSON":
		return FindEntriesJSON[T](ctx, source)
	default:
		return FindEntriesXHTML[T](ctx, source)
	}
}

func FindEntriesJSON[T model.IEntry](ctx context.Context, source SourceDefinition) ([]T, error) {
	logger := zerolog.Ctx(ctx).With().Str("source", source.Name).Logger()

	//nolint:exhaustivestruct
	httpClient := http.Client{Timeout: requestTimeout}

	urls := source.urls()
	limit := source.Limit

	entries := []T{}

	doRequest := func(url string) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {

			return ErrFailToCrateRequest
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := httpClient.Do(req)

		if err != nil {
			return fmt.Errorf("error on request: %w", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error on read body: %w", err)
		}

		resultEntries := gjson.GetBytes(body, source.Attributes.EntrySelector).Array()

		for _, row := range resultEntries {
			if limit == 0 {
				break
			}

			limit--

			title := row.Get(source.Attributes.Title.Path).String()
			link := row.Get(source.Attributes.Link.Path).String()
			image := row.Get(source.Attributes.Image.Path).String()

			// TODO: support categories

			entry := source.buildEntry(title, link, image, []string{}).(T)

			entries = append(entries, entry)

			if limit == 0 {
				logger.Warn().Msgf("Limit reached (%v)", source.Limit)
			}
		}

		return nil
	}

	for _, url := range urls {
		logger.Info().Msgf("Visiting %s", url)

		if err := doRequest(url); err != nil {
			logger.
				Error().
				Err(err).
				Msgf("Fail to visit %s", url)

			return nil, err
		}
	}

	return entries, nil
}

func FindEntriesXHTML[T model.IEntry](ctx context.Context, source SourceDefinition) ([]T, error) {
	logger := zerolog.Ctx(ctx).With().Str("source", source.Name).Logger()

	ctx = logger.WithContext(ctx)

	entries := []T{}

	limit := source.Limit

	if limit == 0 {
		limit = math.MaxInt
	}

	callback := func(e Element) {
		if limit == 0 {
			return
		}

		entry, err := onEntry[T](ctx, source, e)
		if err != nil {
			if !errors.Is(err, ErrCategoryNotAllowed) {
				logger.Error().Err(err).Msg("Error on entry")
			}

			return
		}

		logger.
			Debug().
			Msgf("New entry: %s", entry.Link())

		entries = append(entries, entry)

		limit--

		if limit == 0 {
			logger.Warn().Msgf("Limit reached (%v)", source.Limit)
		}
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

	urls := source.urls()

	for _, url := range urls {
		logger.Info().Msgf("Visiting %s", url)

		if err := collector.Visit(url); err != nil {
			logger.
				Error().
				Err(err).
				Msgf("Fail to visit %s", url)

			return fmt.Errorf("error on visit: %w", err)
		}
	}

	collector.Wait()

	return nil
}

func onEntry[T model.IEntry](ctx context.Context, source SourceDefinition, el Element) (T, error) { //nolint:ireturn
	var result T
	attributes := source.Attributes

	title := attributes.Title.findAttribute(el)

	categories := support.ToLower(attributes.Category.findCategories(el))

	if !attributes.Category.isAllowed(categories) {
		zerolog.Ctx(ctx).Debug().
			Strs("categories", categories).
			Str("title", title).
			Msg(ErrCategoryNotAllowed.Error())

		return result, ErrCategoryNotAllowed
	}

	link := attributes.Link.findAttribute(el)
	image := attributes.Image.findAttribute(el)

	if len(title) > titleLimit {
		title = title[:titleLimit]
	}

	//nolint:forcetypeassert
	result = source.buildEntry(title, link, image, categories).(T)

	return result, nil
}
