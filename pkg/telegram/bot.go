//nolint:exhaustruct
package telegram

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const poolingTiming = 10 * time.Second

type Config struct {
	Token     string  `fig:"token"     yaml:"token"`
	Broadcast []int64 `fig:"broadcast" yaml:"broadcast"`
}

const LoggerKey = "bot:logger"

func NewBot(cfg Config) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: poolingTiming},
		OnError: func(err error, tx telebot.Context) {
			_ = tx.Reply(fmt.Sprintf("Error: %s", err.Error()))
			logger := tx.Get(LoggerKey).(zerolog.Logger) //nolint:forcetypeassert
			logger.Error().Err(err).Msg("Bot error")
		},
	}

	return telebot.NewBot(pref)
}
