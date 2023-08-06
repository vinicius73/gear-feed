package local

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

type Record interface {
	IsDeletedOrExpired() bool
	UserMeta() byte
}

func ApplyWhere(where storage.WhereOptions, item Record, rerr error) bool {
	if rerr != nil {
		if where.AllowMissed != nil && *where.AllowMissed && rerr == badger.ErrKeyNotFound {
			return true
		}

		return false
	}

	if item.IsDeletedOrExpired() {
		return false
	}

	if where.Is != nil && item.UserMeta() == where.Is.Byte() {
		return true
	}

	if where.Not != nil && item.UserMeta() != where.Not.Byte() {
		return true
	}

	return false
}
