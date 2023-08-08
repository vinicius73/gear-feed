package database

import (
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

func Where[T model.IEntry](where storage.WhereOptions, has bool, entry DBEntry[T]) bool {
	if !has {
		return where.AllowMissed != nil && *where.AllowMissed
	}

	if condition := where.Is; condition != nil {
		return condition.Is(entry.Status)
	}

	if condition := where.Not; condition != nil {
		return !condition.Is(entry.Status)
	}

	return true
}
