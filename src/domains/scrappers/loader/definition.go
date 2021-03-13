package loader

import (
	"gfeed/domains/news"
	"strings"

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

	attributes := options.Attributes

	skip := false
	count := 0

	c.OnHTML(attributes.Wrapper, func(e *colly.HTMLElement) {
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
		category := attributes.Category.findAttribute(e)

		if len(attributes.Category.Path) > 0 && len(category) > 0 {
			if !attributes.Category.isAllowed(category) {
				logger.Warn().
					Str("category", category).
					Msgf("Skiped, that category is not allowed: %s", link)
				return
			}
		}

		image := attributes.Image.findAttribute(e)
		title := attributes.Title.findAttribute(e)

		entry := options.buildEntry(title, link, image, category)

		logger.Debug().Msgf("New entry: %s", entry.Link)

		entries = append(entries, entry)
	})

	c.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))

		// fmt.Println(string(r.Body))
	})

	logger.Debug().Msg("Starting...")

	c.Visit(options.visitURL())

	c.Wait()

	logger.Info().Msgf("Done with %v results and %v entries.", count, len(entries))

	if len(entries) == 0 {
		logger.Warn().
			Msg("Empty result")
	}

	return entries
}

func (d Definitions) visitURL() string {
	return d.BaseURL + d.Path
}

func (d Definitions) buildEntry(title, link, image, category string) news.Entry {
	return news.Entry{
		Type:     d.Name,
		Title:    title,
		Category: category,
		Link:     d.absouteURL(link),
		Image:    d.absouteURL(image),
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
