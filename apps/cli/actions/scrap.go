package actions

import (
	"context"

	"github.com/vinicius73/gear-feed/pkg/configurations"
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/sources"
	"github.com/vinicius73/gear-feed/pkg/storage/database"
	"github.com/vinicius73/gear-feed/pkg/tasks"
)

type LoadOptions struct {
	To           int64
	SendResumeTo []int64
	Limit        int
	Sources      sources.LoadOptions
}

func Load(ctx context.Context, opt LoadOptions) error {
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

	botSender, err := buildSender(SenderOptions{
		Chats:    []int64{opt.To},
		Storage:  store,
		Telegram: config.Telegram,
	})
	if err != nil {
		return err
	}

	return tasks.SendLastEntries[model.Entry]{
		Limit:        opt.Limit,
		Sources:      opt.Sources,
		SendResumeTo: opt.SendResumeTo,
	}.
		Run(ctx, tasks.TaskRunOptions[model.Entry]{
			Storage: store,
			Sender:  botSender,
		})
}
