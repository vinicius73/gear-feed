package linkloader

import (
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

type Collection struct {
	SourceName string
	Entries    []scraper.Entry
}

type Collections []Collection

func (c Collections) Entries() []scraper.Entry {
	entries := []scraper.Entry{}

	for _, collection := range c {
		entries = append(entries, collection.Entries...)
	}

	return entries
}

func (c Collections) Shuffle() []scraper.Entry {
	return support.Shuffle(c.Entries())
}
