package linkloader

import "github.com/vinicius73/gamer-feed/pkg/scraper"

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
