package actions

import (
	"context"

	"github.com/vinicius73/gear-feed/pkg/configurations"
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/storage/database"
	"github.com/vinicius73/gear-feed/pkg/tasks"
)

func Cleanup(ctx context.Context) error {
	config := configurations.Ctx(ctx)

	db, err := database.Open(ctx, config.Storage)
	if err != nil {
		return err
	}

	defer db.Close()

	store, err := database.NewStorage[model.Entry](db, config.Storage)
	if err != nil {
		return err
	}

	return tasks.Cleanup[model.Entry]{
		Notify: false,
	}.Run(ctx, tasks.TaskRunOptions[model.Entry]{
		Storage: store,
		Sender:  nil,
	})
}
