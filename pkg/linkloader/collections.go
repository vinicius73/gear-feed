package linkloader

import (
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/support"
)

type Collection[T model.IEntry] struct {
	SourceName string
	Entries    []T
}

type Collections[T model.IEntry] []Collection[T]

func (c Collections[T]) Entries() []T {
	entries := []T{}

	for _, collection := range c {
		entries = append(entries, collection.Entries...)
	}

	return entries
}

func (c Collections[T]) Shuffle() []T {
	return support.Shuffle(c.Entries())
}
