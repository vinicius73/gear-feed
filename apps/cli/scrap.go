package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/telegram"
	"github.com/vinicius73/gamer-feed/sources"
	"gopkg.in/telebot.v3"
)

func loadEntries(ctx context.Context, only []string) (linkloader.Collections, error) {
	list, err := sources.LoadDefinitions(ctx, sources.LoadOptions{
		Only: only,
	})
	if err != nil {
		return nil, err
	}

	collections, err := linkloader.FromSources(ctx, linkloader.LoadOptions{
		Workers: (len(list) + 1) / 2,
		Sources: list,
	})
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func scrapCMD() *cli.Command {
	load := &cli.Command{
		Name:        "load",
		Description: `Load scrap data from sources.`,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "only",
				Usage:   "Load only the specified loaders",
				Aliases: []string{"o"},
			},
			&cli.Int64Flag{
				Name:     "to",
				Usage:    "Send the loaded data to the specified channel",
				Required: true,
			},
		},
		Action: func(cmd *cli.Context) error {
			logger := zerolog.Ctx(cmd.Context)
			config := configurations.Ctx(cmd.Context)

			only := cmd.StringSlice("only")

			collections, err := loadEntries(cmd.Context, only)
			if err != nil {
				return err
			}

			entries := collections.Entries()

			logger.Info().Msgf("Found %d entries", len(entries))

			bot, err := telegram.NewBot(config.Telegram)

			if err != nil {
				return err
			}

			to := telebot.ChatID(cmd.Int64("to"))

			count := 0

			for _, entry := range entries {
				count++

				if count > 4 {
					break
				}

				logger.Info().Msgf("Sending %s to %v", entry.Title, to)

				msg := fmt.Sprintf("%s\n%s\n\n#%s", entry.Title, entry.Link, entry.Type)

				_, err := bot.Send(to, msg)

				time.Sleep(100 * time.Millisecond)

				if err != nil {
					logger.Error().Err(err).Msgf("Failed to send %s to %v", entry.Title, to)
				}
			}

			return nil
		},
	}

	return &cli.Command{
		Name:        "scrap",
		Subcommands: []*cli.Command{load},
	}
}
