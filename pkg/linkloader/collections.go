package linkloader

import (
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

type Collection struct {
	SourceName string
	Entries    []model.Entry
}

type Collections []Collection

func (c Collections) Entries() []model.Entry {
	entries := []model.Entry{}

	for _, collection := range c {
		entries = append(entries, collection.Entries...)
	}

	return entries
}

func (c Collections) Shuffle() []model.Entry {
	return support.Shuffle(c.Entries())
}
