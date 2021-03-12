package voxel

import (
	"gfeed/news"
	log "gfeed/utils/logger"
	"strings"

	"github.com/gocolly/colly"
)

// TYPE of vocel scrapper
const TYPE = "VOXEL"
const baseAddress = "https://www.tecmundo.com.br/voxel"

var logger log.Logger

func init() {
	logger = log.New("scrapper:voxel")
}

// Load voxel news
func Load() []news.Entry {
	entries := []news.Entry{}

	c := colly.NewCollector()

	c.OnHTML("article.tec--voxel-main__item", func(e *colly.HTMLElement) {

		image := e.ChildAttr("img.tec--voxel-main__item__thumb__image", "data-src")
		link := e.ChildAttr("figure > a", "href")
		title := e.ChildText(".tec--voxel-main__item__title a.tec--voxel-main__item__title__link")

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

	c.Visit(baseAddress)

	c.Wait()

	logger.Debug().Msg("Done.")

	if len(entries) == 0 {
		logger.Warn().Msg("Empty entries")
	}

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

	if strings.HasPrefix(e.Image, "//") {
		e.Image = "https:" + e.Image
	}

	return e
}
