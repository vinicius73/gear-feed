package actions

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/botworker"
	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage/local"
)

func BotWorker(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	config := configurations.Ctx(ctx)

	db, err := local.Open(config.Storage)
	if err != nil {
		return err
	}

	defer db.Close()

	store, err := local.NewStorage[model.Entry](db, config.Storage)
	if err != nil {
		return err
	}

	botSender, err := buildSender(SenderOptions{
		Storage:  store,
		Chats:    config.Telegram.Broadcast,
		Telegram: config.Telegram,
	})
	if err != nil {
		return err
	}

	bot := botworker.New[model.Entry](botworker.BotOptions[model.Entry]{
		Storage: store,
		Sender:  botSender,
		Config: botworker.Config[model.Entry]{
			Cron: config.Cron,
		},
	})

	logger.Info().Msg("Starting bot worker")

	if err := bot.Run(ctx); err != nil {
		return err
	}

	logger.Warn().Msg("Bot worker stopped")

	return nil
}
