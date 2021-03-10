package theenemy

import (
	"gfeed/news"
	"gfeed/utils"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// TYPE of the enemy scrapper
const TYPE = "THE_ENEMY"
const baseAddress = "https://www.theenemy.com.br"

var reURL *regexp.Regexp
var logger utils.Logger

func init() {
	reURL = regexp.MustCompile(`\((.*?)\)`)
	logger = utils.NewLogger("scrapper:theenemy")
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

		logger.Debug("New Entry: " + link)

		entries = append(entries, entry)
	})

	logger.Debug("Starting...")

	c.Visit(baseAddress + "/news")

	logger.Debug("Done...")

	c.Wait()

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
