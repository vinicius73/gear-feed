package sender

import (
	"context"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	"gopkg.in/telebot.v3"
)

var _ Serder[model.IEntry] = (*TelegramSerder[model.IEntry])(nil) // Ensure interface implementation

var (
	ErrNoChats    = apperrors.Business("no chats to send message", "SENDER:NO_CHATS")
	ErrFailToSend = apperrors.System(nil, "fail to send message", "SENDER:FAIL_TO_SEND")
)

type SendFileOptions struct {
	FilePath string
	Caption  string
}

type SendResumeOptions struct {
	Resume Resume
	Chats  []int64
}

type SendCleanupNotifyOptions struct {
	Count int64
}

type Serder[T model.IEntry] interface {
	Send(ctx context.Context, entry T) error
	SendCollection(ctx context.Context, entry []T) error
	SendResume(ctx context.Context, opt SendResumeOptions) error
	SendCleanupNotify(ctx context.Context, opt SendCleanupNotifyOptions) error
	SendFile(ctx context.Context, opt SendFileOptions) error
	WithChats(ids []int64) Serder[T]
}

type TelegramSerder[T model.IEntry] struct {
	chats   []telebot.Recipient
	storage storage.Storage[T]
	bot     *telebot.Bot
}

type TelegramOptions[T model.IEntry] struct {
	Chats   []int64
	Storage storage.Storage[T]
}

func NewTelegramSerder[T model.IEntry](bot *telebot.Bot, opts TelegramOptions[T]) TelegramSerder[T] {
	ids := make([]telebot.Recipient, len(opts.Chats))
	for index, chat := range opts.Chats {
		ids[index] = telebot.ChatID(chat)
	}

	return TelegramSerder[T]{
		chats:   ids,
		bot:     bot,
		storage: opts.Storage,
	}
}

func (s TelegramSerder[T]) SendFile(ctx context.Context, opt SendFileOptions) error {
	logger := zerolog.Ctx(ctx).With().Str("file", opt.FilePath).Logger()

	file := &telebot.Document{
		File:                 telebot.FromDisk(opt.FilePath),
		Caption:              opt.Caption,
		FileName:             filepath.Base(opt.FilePath),
		DisableTypeDetection: true,
	}

	for _, chat := range s.chats {
		_, err := s.bot.Send(chat, file, telebot.ModeHTML)
		if err != nil {
			return ErrFailToSend.Wrap(err)
		}

		logger.Info().
			Str("recipient", chat.Recipient()).
			Msgf("File sent")
	}

	return nil
}

func (s TelegramSerder[T]) Send(ctx context.Context, entry T) error {
	logger := zerolog.Ctx(ctx)

	if len(s.chats) == 0 {
		return ErrNoChats
	}

	msg := BuildMessage(entry)

	for _, chat := range s.chats {
		_, err := s.bot.Send(chat, msg)
		if err != nil {
			return ErrFailToSend.Wrap(err)
		}

		logger.Info().
			Str("recipient", chat.Recipient()).
			Strs("tags", entry.Tags()).
			Msgf("Message sent %s", entry.Link())
	}

	err := s.storage.Store(storage.Entry[T]{
		Data:   entry,
		Status: storage.StatusSent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s TelegramSerder[T]) SendCollection(ctx context.Context, entries []T) error {
	logger := zerolog.Ctx(ctx)

	size := len(entries)

	if size == 0 {
		logger.Warn().Msg("No entries to send")

		return nil
	}

	sendInterval := CalculeSendInterval(size)

	logger.Info().
		Msgf("Sending %d entries with %s interval, speding %s", size, sendInterval, sendInterval*time.Duration(size))

	startedAt := time.Now()

	for _, item := range entries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sendInterval):
			err := s.Send(ctx, item)
			if err != nil {
				logger.Error().Err(err).Msg("Error sending message")

				return err
			}
		}
	}

	dur := time.Since(startedAt)

	logger.Info().Dur("spend", dur).Msgf("Finished sending %d entries in %s", size, dur)

	return nil
}

func (s TelegramSerder[T]) SendResume(ctx context.Context, opt SendResumeOptions) error {
	logger := zerolog.Ctx(ctx)

	chats := s.chats

	if len(opt.Chats) > 0 {
		chats = make([]telebot.Recipient, len(opt.Chats))

		for index, chat := range opt.Chats {
			chats[index] = telebot.ChatID(chat)
		}
	}

	if len(chats) == 0 {
		return ErrNoChats
	}

	msg := opt.Resume.HTML()

	for _, chat := range chats {
		_, err := s.bot.Send(chat, msg, telebot.ModeHTML)
		if err != nil {
			return ErrFailToSend.Wrap(err)
		}

		logger.Info().
			Str("recipient", chat.Recipient()).
			Msgf("Resume sent")
	}

	return nil
}

func (s TelegramSerder[T]) SendCleanupNotify(ctx context.Context, opt SendCleanupNotifyOptions) error {
	logger := zerolog.Ctx(ctx)

	if len(s.chats) == 0 {
		return ErrNoChats
	}

	msg := BuildCleanupMessage(opt.Count)

	for _, chat := range s.chats {
		_, err := s.bot.Send(chat, msg, telebot.ModeHTML)
		if err != nil {
			return ErrFailToSend.Wrap(err)
		}

		logger.Info().
			Str("recipient", chat.Recipient()).
			Msgf("Cleanup notify sent")
	}

	return nil
}

// WithChats add chats to send messages.
func (s TelegramSerder[T]) WithChats(ids []int64) Serder[T] {
	chats := make([]telebot.Recipient, len(ids))

	for index, chat := range ids {
		chats[index] = telebot.ChatID(chat)
	}

	s.chats = append(s.chats, chats...)

	return s
}
