package bot

import (
	"gfeed/news"
	"gfeed/news/data"
	"gfeed/scrappers"
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
	b, err := newClient(c)

	if err != nil {
		logger.Fatal(err)
		return
	}

	err = sendNews(b, c)

	if err != nil {
		logger.Fatal(err)
		return
	}
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
		logger.Info("Sending: " + entry.Link)

		data.Put(entry)

		b.Send(chat, buildMsg(entry))

		time.Sleep(time.Second * 1)
	}

	return nil
}

func buildMsg(entry news.Entry) string {
	var builder strings.Builder

	builder.WriteString(entry.Title)
	builder.WriteString("\n")
	builder.WriteString(entry.Link)
	builder.WriteString("\n")
	builder.WriteString("#" + entry.Type)

	return builder.String()
}
