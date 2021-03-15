package bot

import (
	"fmt"
	"gfeed/domains/news"
	"gfeed/domains/news/storage"
	"gfeed/domains/scrappers"
	"math/rand"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SendNews to channel
func SendNews(c Config) {
	startTime := time.Now()

	if c.DryRun {
		logger.Warn().Msg("DryRun ON")
	}

	b, err := newClient(c)

	if err != nil {
		logger.Fatal().Err(err).Msg("Fail to create bot client")
		return
	}

	err = sendNews(b, c)

	if err != nil {
		logger.Fatal().Err(err).Msg("Fail to send news")
		return
	}

	logger.
		Info().
		Msgf("All done (%s)", time.Since(startTime).String())
}

func sendNews(b *tb.Bot, c Config) error {
	chat, err := b.ChatByID(c.Channel)

	if err != nil {
		return err
	}

	entries := scrappers.NewsEntries()

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	for _, entry := range entries {
		logger.Info().Msgf("Sending: %s", entry.Link)

		if !c.DryRun {
			storage.Put(entry)
			b.Send(chat, buildMsg(entry))
		}

		time.Sleep(time.Microsecond * 100)
	}

	return sendResume(b, c, entries)
}

func buildMsg(entry news.Entry) string {
	var builder strings.Builder

	builder.WriteString(entry.Title)
	builder.WriteString("\n")
	builder.WriteString(entry.Link)
	builder.WriteString("\n")
	builder.WriteString("#" + entry.Type)

	if len(entry.Category) > 0 {
		builder.WriteString(" #" + strings.ReplaceAll(entry.Category, " ", "_"))
	}

	return builder.String()
}

func sendResume(b *tb.Bot, c Config, entries []news.Entry) error {

	chat, err := b.ChatByID(c.User)

	if err != nil {
		return err
	}

	var builder strings.Builder

	if c.DryRun {
		builder.WriteString("🧪")
	}

	builder.WriteString("🤖")
	builder.WriteString(time.Now().String())
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Gamer Feed was executed [%v entries]", len(entries)))

	_, err = b.Send(chat, builder.String())

	return err
}
