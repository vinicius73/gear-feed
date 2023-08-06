package actions

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/linkloader/news"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage/local"
	"github.com/vinicius73/gamer-feed/sources"
)

type LoadOptions struct {
	Only  []string
	To    int64
	Limit int
}

func Load(ctx context.Context, opt LoadOptions) error {
	config := configurations.Ctx(ctx)

	definitions, err := sources.LoadDefinitions(ctx, sources.LoadOptions{
		Only: opt.Only,
	})
	if err != nil {
		return err
	}

	db, err := local.Open(config.Storage)
	if err != nil {
		return err
	}

	defer db.Close()

	store, err := local.NewStorage[model.Entry](db, config.Storage)
	if err != nil {
		return err
	}

	entries, err := news.LoadEntries(ctx, news.LoadOptions{
		LoadOptions: linkloader.LoadOptions{
			Sources: definitions,
			Workers: 0, // dynamic
		},
		Limit:   opt.Limit,
		Storage: store,
	})
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

	err = botSender.SendCollection(ctx, entries)

	return err
}
