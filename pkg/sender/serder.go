package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"gopkg.in/telebot.v3"
)

type Serder interface {
	Send(ctx context.Context, entry Sendable) error
	SendCollection(ctx context.Context, entry []Sendable) error
}

type TelegramSerder struct {
	chats   []telebot.Recipient
	storage storage.Storage[Sendable]
	bot     *telebot.Bot
}

type TelegramOptions struct {
	Chats   []int64
	Storage storage.Storage[Sendable]
}

func NewTelegramSerder(bot *telebot.Bot, opts TelegramOptions) Serder {
	ids := make([]telebot.Recipient, len(opts.Chats))
	for index, chat := range opts.Chats {
		ids[index] = telebot.ChatID(chat)
	}

	return TelegramSerder{
		chats:   ids,
		bot:     bot,
		storage: opts.Storage,
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

		err = s.storage.Store(storage.Entry[Sendable]{
			Data:   entry,
			Status: storage.StatusSent,
		})
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
	originalSize := len(entries)

	if originalSize == 0 {
		logger.Info().Msg("No entries to send")
		return nil
	}

	where := storage.Where(storage.WhereIs(storage.StatusNew), storage.WhereAllowMissed(true))
	entries, err := s.storage.Where(where, entries)
	if err != nil {
		return err
	}

	size := len(entries)

	if size == 0 {
		logger.Info().Msg("No new entries to send")
		return nil
	}

	if size != originalSize {
		logger.Info().Msgf("Found %d new entries", size)
	}

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
