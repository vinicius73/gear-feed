package main

import (
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/apps/cli/actions"
	"github.com/vinicius73/gamer-feed/pkg/sources"
)

//nolint:funlen
func scrapCMD() *cli.Command {
	load := &cli.Command{
		Name:        "load",
		Description: `Load scrap data from sources.`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Usage:   "Limit the number of entries to send",
				Aliases: []string{"l"},
				Value:   10, //nolint:gomnd // default value
				Action: func(c *cli.Context, val int) error {
					if val < 1 {
						return cli.Exit("Limit must be greater than 1", 1)
					}

					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:    "only",
				Aliases: []string{"o"},
				Usage:   "Only show the specified sources",
			},
			&cli.StringSliceFlag{
				Name:     "sources",
				Aliases:  []string{"s"},
				Usage:    "Load sources from the specified paths",
				Required: true,
			},
			&cli.Int64Flag{
				Name:     "to",
				Usage:    "Send the loaded data to the specified channel",
				Required: true,
			},
			&cli.Int64SliceFlag{
				Name:    "send-resume-to",
				Usage:   "Send the loaded data to the specified channel",
				Aliases: []string{"r"},
			},
		},
		Action: func(cmd *cli.Context) error {
			return actions.Load(cmd.Context, actions.LoadOptions{
				To:           cmd.Int64("to"),
				Limit:        cmd.Int("limit"),
				SendResumeTo: cmd.Int64Slice("send-resume-to"),
				Sources: sources.LoadOptions{
					Only:  cmd.StringSlice("only"),
					Paths: cmd.StringSlice("sources"),
				},
			})
		},
	}

	stories := &cli.Command{
		Name:        "stories",
		Description: `Send entries and send stories to telegram.`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Usage:   "Limit the number of entries to send",
				Aliases: []string{"l"},
				Value:   10, //nolint:gomnd // default value
				Action: func(c *cli.Context, val int) error {
					if val < 1 {
						return cli.Exit("Limit must be greater than 1", 1)
					}

					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:    "only",
				Aliases: []string{"o"},
				Usage:   "Only show the specified sources",
			},
			&cli.StringSliceFlag{
				Name:     "sources",
				Aliases:  []string{"s"},
				Usage:    "Load sources from the specified paths",
				Required: true,
			},
			&cli.DurationFlag{
				Name:     "period",
				Aliases:  []string{"p"},
				Usage:    "Period to load entries and build stories",
				Value:    time.Hour * 48, //nolint:gomnd // default value
				Required: false,
			},
			&cli.Int64Flag{
				Name:     "to",
				Usage:    "Send the loaded data to the specified channel",
				Required: true,
			},
		},
		Action: func(cmd *cli.Context) error {
			return actions.SendStories(cmd.Context, actions.SendStoriesOptions{
				To:     cmd.Int64("to"),
				Limit:  cmd.Int("limit"),
				Period: cmd.Duration("period"),
				// SendResumeTo: cmd.Int64Slice("send-resume-to"),
				Sources: sources.LoadOptions{
					Only:  cmd.StringSlice("only"),
					Paths: cmd.StringSlice("sources"),
				},
			})
		},
	}

	return &cli.Command{
		Name:        "scrap",
		Subcommands: []*cli.Command{load, stories},
	}
}
