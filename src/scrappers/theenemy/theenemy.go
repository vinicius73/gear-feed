package theenemy

import (
	"gfeed/news"
	log "gfeed/utils/logger"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// TYPE of the enemy scrapper
const TYPE = "THE_ENEMY"
const baseAddress = "https://www.theenemy.com.br"

var reURL *regexp.Regexp
var logger log.Logger

func init() {
	reURL = regexp.MustCompile(`\((.*?)\)`)
	logger = log.New("scrapper:theenemy")
}

// Load the enemy news
func Load() []news.Entry {
	entries := []news.Entry{}

	c := colly.NewCollector()

	c.OnHTML("div[data-content-type='featured'] .card__wrapper", func(e *colly.HTMLElement) {

		style := e.ChildAttr("a.card__image__anchor", "style")
		link := e.ChildAttr(".card__content > a", "href")
		title := e.ChildText(".card__content > a")

		entry := buildEntry(title, link, style)

		logger.Debug().Msgf("New Entry: %s", link)

		entries = append(entries, entry)
	})

	c.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))
	})

	logger.Debug().Msg("Starting...")

	c.Visit(baseAddress + "/news")

	logger.Debug().Msg("Done.")

	c.Wait()

	if len(entries) == 0 {
		logger.Warn().Msg("Empty entries")
	}

	return entries
}

func buildEntry(title, link, style string) (e news.Entry) {
	e.Link = baseAddress + link
	e.Title = title
	e.Type = TYPE
	e.Image = parseStyle(style)

	return e
}

func parseStyle(style string) string {
	result := reURL.FindString(style)

	result = strings.TrimLeft(result, "(")
	result = strings.TrimRight(result, ")")

	if strings.HasPrefix(result, "//") {
		return "https:" + result
	}

	return result

}
