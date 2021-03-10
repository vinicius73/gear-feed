package theenemy

import (
	"gfeed/news"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly"
)

// TYPE of the enemy scrapper
const TYPE = "THE_ENEMY"
const baseAddress = "https://www.theenemy.com.br"

var reURL *regexp.Regexp
var logger *log.Entry

func init() {
	reURL = regexp.MustCompile(`\((.*?)\)`)
	logger = log.New().WithField("scrapper", "theenemy")
}

// Load the enemy news
func Load() []news.Entry {
	entries := []news.Entry{}

	c := colly.NewCollector()

	c.OnHTML("div[data-content-type='featured'] .card__wrapper", func(e *colly.HTMLElement) {

		style := e.ChildAttr("a.card__image__anchor", "style")
		link := e.ChildAttr(".card__content > a", "href")
		title := e.ChildText(".card__content > a")

		log.Info("New Entry: " + link)

		entries = append(entries, buildEntry(title, link, style))
	})

	c.Visit(baseAddress + "/news")

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
