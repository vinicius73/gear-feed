package bot

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// newClient of telegram bot
func newClient(c Config) (*tb.Bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	return b, err
}
