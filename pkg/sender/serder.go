package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Serder interface {
	Send(ctx context.Context, entry Sendable) error
	SendCollection(ctx context.Context, entry []Sendable) error
}

type TelegramSerder struct {
	chats []telebot.Recipient

	bot *telebot.Bot
}

func NewTelegramSerder(bot *telebot.Bot, chats []int64) Serder {
	ids := make([]telebot.Recipient, len(chats))
	for index, chat := range chats {
		ids[index] = telebot.ChatID(chat)
	}

	return TelegramSerder{
		chats: ids,
		bot:   bot,
	}
}

func (s TelegramSerder) Send(ctx context.Context, entry Sendable) error {
	logger := zerolog.Ctx(ctx)

	msg := BuildMessage(entry)

	for _, chat := range s.chats {
		_, err := s.bot.Send(chat, msg)
		if err != nil {
			return err
		}

		logger.Info().
			Str("recipient", chat.Recipient()).
			Strs("tags", entry.Tags()).
			Msgf("Message sent %s", entry.Link())
	}

	return nil
}

func (s TelegramSerder) SendCollection(ctx context.Context, entries []Sendable) error {
	logger := zerolog.Ctx(ctx)

	size := len(entries)

	sendInterval := CalculeSendInterval(size)

	logger.Info().Msgf("Sending %d entries with %s interval, speding %s", size, sendInterval, sendInterval*time.Duration(size))

	startedAt := time.Now()

	for _, item := range entries {
		err := s.Send(ctx, item)
		if err != nil {
			logger.Error().Err(err).Msg("error sending message")

			return fmt.Errorf("error sending message: %w", err)
		}

		time.Sleep(sendInterval)
	}

	dur := time.Since(startedAt)

	logger.Info().Dur("spend", dur).Msgf("Finished sending %d entries in %s", size, dur)

	return nil
}
