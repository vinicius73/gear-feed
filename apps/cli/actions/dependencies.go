package actions

import (
	"context"
	"database/sql"

	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/storage/database"
)

func buildDB[T model.IEntry](
	ctx context.Context,
	config *configurations.AppConfig,
) (storage.Storage[T], *sql.DB, error) {
	db, err := database.Open(ctx, config.Storage)
	if err != nil {
		return nil, db, err
	}

	store, err := database.NewStorage[T](db, config.Storage)

	return store, db, err
}
