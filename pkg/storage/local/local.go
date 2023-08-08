package local

import (
	"errors"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

var _ storage.Storage[storage.Hashable] = Storage[storage.Hashable]{} // Ensure interface implementation

type Options struct {
	storage.Options `fig:",squash" yaml:",inline"`
	Path            string `fig:"path"    yaml:"path"`
}

type Storage[T storage.Hashable] struct {
	ttl time.Duration
	db  *badger.DB
}

func Open(opt Options) (*badger.DB, error) {
	return badger.Open(badger.DefaultOptions(opt.Path).
		WithMaxLevels(3).             //nolint:gomnd
		WithValueLogMaxEntries(50).   //nolint:gomnd
		WithIndexCacheSize(20 << 20), //nolint:gomnd // 20mb
	)
}

func NewStorage[T storage.Hashable](db *badger.DB, opt Options) (storage.Storage[T], error) {
	return Storage[T]{
		ttl: opt.TTL,
		db:  db,
	}, nil
}

func (s Storage[T]) Has(hash string) (bool, error) {
	var has bool
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(hash))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}

		has = err == nil

		return nil
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (s Storage[T]) Store(value storage.Entry[T]) error {
	data, err := value.Marshal()
	if err != nil {
		return storage.ErrFailToMarshalData.Wrap(err)
	}

	hash, err := value.Hash()
	if err != nil {
		return err
	}

	entry := badger.NewEntry(hash, data).WithMeta(value.Status.Byte())

	if value.Status.Is(storage.StatusNew) && s.ttl > 0 {
		entry = entry.WithTTL(s.ttl)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (s Storage[T]) Where(opts storage.WhereOptions, list []T) ([]T, error) {
	var result []T

	err := s.db.View(func(txn *badger.Txn) error {
		for _, value := range list {
			hash, err := value.Hash()
			if err != nil {
				return err
			}

			entry, err := txn.Get([]byte(hash))

			if ApplyWhere(opts, entry, err) {
				result = append(result, value)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
