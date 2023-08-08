//nolint:ireturn
package database

import (
	"encoding/json"
	"time"

	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

type DBEntry[T model.IEntry] struct {
	Hash       string         `db:"hash,primarykey"`
	SourceName string         `db:"source_name"`
	ImageURL   string         `db:"image_url"`
	Text       string         `db:"text"`
	Categories []byte         `db:"categories"`
	URL        string         `db:"url"`
	Status     storage.Status `db:"status"`
	CreatedAt  time.Time      `db:"created_at"`
	TTL        time.Time      `db:"ttl"`
}

func (e DBEntry[T]) ToEntry(target T) T {
	//nolint:forcetypeassert
	return target.FillFrom(model.Entry{
		Title:      e.Text,
		URL:        e.URL,
		Image:      e.ImageURL,
		Categories: []string{},
		SourceName: e.SourceName,
	}).(T)
}

func NewEntry[T model.IEntry](ttl time.Duration, entry storage.Entry[T]) (DBEntry[T], error) {
	source := entry.Data

	hash, err := source.Hash()
	if err != nil {
		return DBEntry[T]{}, err
	}

	categories, err := json.Marshal(source.Tags())
	if err != nil {
		return DBEntry[T]{}, err
	}

	return DBEntry[T]{
		Hash:       hash,
		SourceName: source.Source(),
		ImageURL:   source.ImageURL(),
		Text:       source.Text(),
		URL:        source.Link(),
		Status:     entry.Status,
		Categories: categories,
		CreatedAt:  time.Now(),
		TTL:        time.Now().Add(ttl),
	}, nil
}
