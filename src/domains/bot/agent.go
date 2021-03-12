package bot

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func Agent(c Config) {
	client, err := newClient(c)

	if err != nil {
		logger.Fatal().Err(err).Msg("Fail to create bot client")
		return
	}

	client.Handle("/me", func(m *tb.Message) {
		_, err := client.Send(m.Sender, fmt.Sprintf("Your id: %v", m.Sender.ID))

		if err != nil {
			logger.
				Warn().
				Err(err).
				Int("sender.id", m.Sender.ID).
				Msg("Fail to send message")
		}
	})

	logger.
		Info().
		Msg("Starting agent")

	client.Start()
}
