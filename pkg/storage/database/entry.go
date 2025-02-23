package database

import (
	"encoding/json"
	"time"

	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/storage"
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
	HasStory   bool           `db:"has_story"`
	TTL        time.Time      `db:"ttl"`
}

type DBEntryToUpdate[T model.IEntry] struct {
	Hash     string         `db:"hash,primarykey"`
	Status   storage.Status `db:"status"`
	HasStory bool           `db:"has_story"`
}

func (e DBEntry[T]) ToEntry(target T) T {
	//nolint:forcetypeassert
	return target.FillFrom(model.Entry{
		Title:      e.Text,
		URL:        e.URL,
		Image:      e.ImageURL,
		HaveStory:  e.HasStory,
		SourceName: e.SourceName,
		Categories: []string{},
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
		HasStory:   source.HasStory(),
		Status:     entry.Status,
		Categories: categories,
		CreatedAt:  time.Now(),
		TTL:        time.Now().Add(ttl),
	}, nil
}

func EntryToUpdate[T model.IEntry](entry storage.Entry[T]) (DBEntryToUpdate[T], error) {
	source := entry.Data

	hash, err := source.Hash()
	if err != nil {
		return DBEntryToUpdate[T]{}, err
	}

	return DBEntryToUpdate[T]{
		Hash:     hash,
		HasStory: source.HasStory(),
		Status:   entry.Status,
	}, nil
}
