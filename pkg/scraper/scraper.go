package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var (
	ErrCategoryNotAllowed = errors.New("category not allowed")
	ErrFailToCrateRequest = errors.New("fail to create request")
)

const (
	requestTimeout = time.Second * 15
	titleLimit     = 150
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var userAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:135.0) Gecko/20100101 Firefox/135.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	// Add more user agents as needed
}

var _tmpDir string

func init() {
	tmpDir, err := os.UserCacheDir()
	if err != nil {
		tmpDir = os.TempDir()
	}

	_tmpDir = filepath.Join(tmpDir, "gamer-feed/colly")
}

func getRandomUserAgent() string {
	return userAgents[randGen.Intn(len(userAgents))]
}

// getBrowserHeaders returns a set of headers that mimic a typical browser request.
func getBrowserHeaders() http.Header {
	return http.Header{
		"User-Agent":                {getRandomUserAgent()},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Connection":                {"keep-alive"},
		"Upgrade-Insecure-Requests": {"1"},
	}
}

func newCollector(ctx context.Context) *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(getRandomUserAgent()),
		colly.MaxDepth(1),
		colly.Async(true),
		colly.IgnoreRobotsTxt(),
		colly.CacheDir(_tmpDir),
		colly.AllowURLRevisit(),
	)

	c.SetRequestTimeout(requestTimeout)

	logger := zerolog.Ctx(ctx)

	logger.Debug().Str("tmpDir", _tmpDir).Msg("colly.NewCollector: Using temporary directory")

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

	httpClient := http.Client{Timeout: requestTimeout}

	urls := source.urls()
	limit := source.Limit

	entries := []T{}

	// Shuffle URLs to avoid scraping in a predictable order
	randGen.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	doRequest := func(url string) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return ErrFailToCrateRequest
		}

		// Set browser-like headers
		headers := getBrowserHeaders()
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

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
			logger.Error().Err(err).Msgf("Fail to visit %s", url)

			return nil, err
		}

		// Add a random delay after each request (0-5 seconds)
		time.Sleep(time.Duration(randGen.Intn(5)) * time.Second)
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

		logger.Debug().Msgf("New entry: %s", entry.Link())
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

	logger.Info().Dur("duration", duration).Msgf("Done with %v entries (%s).", len(entries), duration.String())

	return entries, nil
}

func visit(ctx context.Context, source SourceDefinition, callback func(e Element)) error {
	logger := zerolog.Ctx(ctx)

	collector := newCollector(ctx)
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

	// Shuffle URLs to avoid predictable scraping patterns
	randGen.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	for _, url := range urls {
		logger.Info().Msgf("Visiting %s", url)

		// Use browser-like headers for each request
		headers := getBrowserHeaders()

		if err := collector.Request("GET", url, nil, nil, headers); err != nil {
			logger.Error().Err(err).Str("url", url).Msgf("Fail to visit %s", url)

			return fmt.Errorf("error on visit (%s): %w", url, err)
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

	result = source.buildEntry(title, link, image, categories).(T)

	return result, nil
}
