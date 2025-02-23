package actions

import (
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/sender"
	"github.com/vinicius73/gear-feed/pkg/storage"
	"github.com/vinicius73/gear-feed/pkg/telegram"
)

type SenderOptions struct {
	Chats    []int64
	Storage  storage.Storage[model.Entry]
	Telegram telegram.Config
}

func buildSender(opt SenderOptions) (sender.Serder[model.Entry], error) {
	bot, err := telegram.NewBot(opt.Telegram)
	if err != nil {
		return nil, err
	}

	return sender.NewTelegramSerder(bot, sender.TelegramOptions[model.Entry]{
		Storage: opt.Storage,
		Chats:   opt.Chats,
	}), nil
}
