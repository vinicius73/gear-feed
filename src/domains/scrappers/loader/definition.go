package loader

import (
	"gfeed/domains/news"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Definitions struct {
	Name       string           `yaml:"name"`
	Enabled    bool             `yaml:"enabled"`
	BaseURL    string           `yaml:"base_url"`
	Path       string           `yaml:"path"`
	Limit      int8             `yaml:"limit"`
	Attributes AttributesFinder `yaml:"attributes"`
}

// FindEnties from definition
func (options Definitions) FindEnties() []news.Entry {
	entries := []news.Entry{}

	logger := baseLogger.With().Str("scrapper", options.Name).Logger()

	c := colly.NewCollector()
	startTime := time.Now()

	attributes := options.Attributes

	skip := false
	count := 0

	c.SetRequestTimeout(time.Second * 15)

	callback := func(e Element, parser string) {
		count++
		if skip {
			return
		}

		if len(entries) >= int(options.Limit) {
			skip = true

			logger.
				Warn().
				Msgf("There is more than %v entries, skiping all new entries", options.Limit)

			return
		}

		link := attributes.Link.findAttribute(e)
		categories := attributes.Category.findCategories(e)

		if len(attributes.Category.Path) > 0 && len(categories) > 0 {
			if !attributes.Category.isAllowed(categories) {
				logger.Warn().
					Strs("categories", categories).
					Msgf("Skiped, that category is not allowed: %s", link)
				return
			}
		}

		image := attributes.Image.findAttribute(e)
		title := attributes.Title.findAttribute(e)

		entry := options.buildEntry(title, link, image, categories)

		logger.
			Debug().
			Str("parser", parser).
			Msgf("New entry: %s", entry.Link)

		entries = append(entries, entry)
	}

	c.OnXML(attributes.Wrapper, func(e *colly.XMLElement) {
		callback(e, "XML")
	})

	c.OnHTML(attributes.Wrapper, func(e *colly.HTMLElement) {
		callback(e, "HTML")
	})

	c.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))
	})

	logger.
		Debug().
		Msg("Starting...")

	c.Visit(options.visitURL())

	c.Wait()

	duration := time.Since(startTime)

	logger.
		Info().
		Dur("duration", duration).
		Msgf("Done with %v results and %v entries (%s).", count, len(entries), duration.String())

	if len(entries) == 0 {
		logger.Warn().
			Msg("Empty result")
	}

	return entries
}

func (d Definitions) visitURL() string {
	return d.BaseURL + d.Path
}

func (d Definitions) buildEntry(title, link, image string, categories []string) news.Entry {
	return news.Entry{
		Type:       d.Name,
		Title:      title,
		Categories: categories,
		Link:       d.absouteURL(link),
		Image:      d.absouteURL(image),
	}
}

func (d Definitions) absouteURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}

	if strings.HasPrefix(path, "//") {
		return "https:" + path
	}

	return d.BaseURL + path
}
