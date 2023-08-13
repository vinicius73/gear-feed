package database

import (
	"database/sql"
	"time"

	"github.com/go-gorp/gorp/v3"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var ErrFailedToCreateEntry = apperrors.System(nil, "failed to create entry", "DB:FailedToCreateEntry")

type Storage[T model.IEntry] struct {
	ttl time.Duration
	db  *gorp.DbMap
}

func NewStorage[T model.IEntry](db *sql.DB, opt Options) (storage.Storage[T], error) {
	dbmap := &gorp.DbMap{
		Db: db, Dialect: gorp.SqliteDialect{},
		ExpandSliceArgs: true,
		TypeConverter:   nil,
	}

	//nolint:exhaustruct
	dbmap.AddTableWithName(DBEntry[T]{}, "entries")

	return Storage[T]{
		ttl: opt.TTL,
		db:  dbmap,
	}, nil
}

func (s Storage[T]) Has(hash string) (bool, error) {
	res, err := s.db.SelectNullStr("SELECT hash FROM entries WHERE hash = ?", hash)

	return res.Valid, err
}

func (s Storage[T]) Store(entry storage.Entry[T]) error {
	record, err := NewEntry[T](s.ttl, entry)
	if err != nil {
		return err
	}

	err = s.db.Insert(&record)

	if err != nil {
		return ErrFailedToCreateEntry.Wrap(err)
	}

	return nil
}

func (s Storage[T]) Where(where storage.WhereOptions, list []T) ([]T, error) {
	hashMap, hashs, err := GroupByHash(list)
	if err != nil {
		return nil, err
	}

	found := []DBEntry[T]{}

	// Select only fields that are needed (hash and status)
	_, err = s.db.Select(&found, "SELECT hash, status FROM entries WHERE hash IN (:hashs)", map[string]interface{}{
		"hashs": hashs,
	})

	if err != nil {
		return nil, err
	}

	foundMap := map[string]DBEntry[T]{}

	for _, entry := range found {
		foundMap[entry.Hash] = entry
	}

	result := []T{}

	for hash, entry := range hashMap {
		dbEntry, has := foundMap[hash]

		if Where(where, has, dbEntry) {
			result = append(result, entry)
		}
	}

	return result, nil
}

func (s Storage[T]) Cleanup() (int64, error) {
	res, err := s.db.Exec("DELETE FROM entries WHERE ttl < ?", time.Now())
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func GroupByHash[T model.IEntry](entries []T) (map[string]T, []string, error) {
	hashMap := map[string]T{}
	hashs := make([]string, len(entries))

	for index, entry := range entries {
		hash, err := entry.Hash()
		if err != nil {
			return nil, nil, err
		}

		hashMap[hash] = entry
		hashs[index] = hash
	}

	return hashMap, hashs, nil
}
