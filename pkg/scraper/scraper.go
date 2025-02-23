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
	ErrCategoryNotAllowed  = errors.New("category not allowed")
	ErrFailToCrateRequest  = errors.New("fail to create request")
	ErrCloudflareChallenge = errors.New("cloudflare challenge detected")
)

const (
	requestTimeout = time.Second * 15
	titleLimit     = 150
	maxRetries     = 3
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:90.0) Gecko/20100101 Firefox/90.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1 Mobile/15E148 Safari/604.1",
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

// getBrowserHeaders retorna cabeçalhos que imitam um navegador real.
func getBrowserHeaders() http.Header {
	return http.Header{
		"User-Agent":                {getRandomUserAgent()},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Connection":                {"keep-alive"},
		"Charset":                   {"UTF-8"},
		"Upgrade-Insecure-Requests": {"1"},
		"DNT":                       {"1"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Cache-Control":             {"max-age=0"},
	}
}

func newCollector(ctx context.Context) *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(getRandomUserAgent()),
		colly.MaxDepth(1),
		colly.Async(false), // Síncrono para melhor controle de re-tentativas
		colly.IgnoreRobotsTxt(),
		// colly.CacheDir(_tmpDir),
		colly.AllowURLRevisit(),
	)
	c.SetRequestTimeout(requestTimeout)
	// logger := zerolog.Ctx(ctx)
	// logger.Debug().Str("tmpDir", _tmpDir).Msg("colly.NewCollector: Using temporary directory")

	return c
}

// isCloudflareChallenge verifica se a resposta é um desafio da Cloudflare.
func isCloudflareChallenge(body []byte) bool {
	strBody := string(body)

	return strings.Contains(strBody, "Attention Required! | Cloudflare") ||
		strings.Contains(strBody, "id=\"challenge-form\"")
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

	randGen.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	doRequest := func(url string) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return ErrFailToCrateRequest
		}
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

		if isCloudflareChallenge(body) {
			return ErrCloudflareChallenge
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
		for attempt := 0; attempt <= maxRetries; attempt++ {
			err := doRequest(url)
			if err == nil {
				break
			}
			if errors.Is(err, ErrCloudflareChallenge) {
				logger.Warn().Msgf("Cloudflare challenge detected, retrying (%d/%d)", attempt+1, maxRetries)
				time.Sleep(time.Duration(randGen.Intn(5)+5) * time.Second) // Atraso de 5-10 segundos
			} else {
				logger.Error().Err(err).Msgf("Fail to visit %s", url)

				return nil, err
			}
		}
		time.Sleep(time.Duration(randGen.Intn(5)) * time.Second) // Atraso após cada requisição
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
	collector.Async = false // Síncrono para facilitar re-tentativas
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

	var isChallenge bool
	collector.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))
		if isCloudflareChallenge(r.Body) {
			isChallenge = true
		} else {
			isChallenge = false
		}
	})

	collector.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
		isChallenge = false
	})

	urls := source.urls()
	randGen.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	for _, url := range urls {
		logger.Info().Msgf("Visiting %s", url)
		for attempt := 0; attempt <= maxRetries; attempt++ {
			isChallenge = false
			headers := getBrowserHeaders()
			err := collector.Request("GET", url, nil, nil, headers)
			if err != nil {
				logger.Error().Err(err).Str("url", url).Msg("Fail to visit")

				return fmt.Errorf("error on visit (%s): %w", url, err)
			}
			collector.Wait()
			if !isChallenge {
				break
			}
			logger.Warn().Msgf("Cloudflare challenge detected, retrying (%d/%d)", attempt+1, maxRetries)
			time.Sleep(time.Duration(randGen.Intn(5)+5) * time.Second)
		}
		if isChallenge {
			logger.Error().Msgf("Failed to bypass Cloudflare after %d retries", maxRetries)

			return fmt.Errorf("failed to bypass Cloudflare for %s", url)
		}
		time.Sleep(time.Duration(randGen.Intn(5)) * time.Second) // Atraso após cada requisição
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
