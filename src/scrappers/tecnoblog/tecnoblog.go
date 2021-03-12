package tecnoblog

import (
	"gfeed/news"
	log "gfeed/utils/logger"
	"strings"

	"github.com/gocolly/colly"
)

// TYPE of tecnoblog scrapper
const TYPE = "TECNOBLOG"
const baseAddress = "https://tecnoblog.net"

var logger log.Logger

func init() {
	logger = log.New("scrapper:tecnoblog")
}

// Load voxel news
func Load() []news.Entry {
	entries := []news.Entry{}

	c := colly.NewCollector()

	c.OnHTML("div.posts article.bloco", func(e *colly.HTMLElement) {
		category := e.ChildText(".thumb .spread a")
		link := e.ChildAttr(".texts a", "href")

		if !isAllowed(category) {
			logger.Warn().Msgf("Skiped: %s", link)
			return
		}

		image := e.ChildAttr(".thumb img", "data-lazy-src")
		title := e.ChildText(".texts h2 a")

		entry := buildEntry(title, link, image)

		logger.Debug().Msgf("New entry: %s", entry.Link)

		entries = append(entries, entry)
	})

	c.OnError(func(r *colly.Response, e error) {
		logger.Error().Err(e).Msg("Response failure")
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Debug().Msgf("Response: %v / %v", r.StatusCode, len(r.Body))
	})

	logger.Debug().Msg("Starting...")

	c.Visit(baseAddress + "/cat/games-jogos/")

	c.Wait()

	logger.Debug().Msg("Done.")

	if len(entries) > 2 {
		return entries[0:2]
	}

	return entries
}

func buildEntry(title, link, image string) (e news.Entry) {
	e.Link = link
	e.Title = title
	e.Type = TYPE
	e.Image = image

	if strings.HasPrefix(e.Image, "/") {
		e.Image = baseAddress + e.Image
	}

	return e
}

func isAllowed(cat string) bool {
	c := strings.ToLower(cat)

	if c == "news" {
		return true
	}

	if c == "especial" {
		return true
	}

	return false
}
