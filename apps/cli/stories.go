package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gear-feed/apps/cli/actions"
	"github.com/vinicius73/gear-feed/pkg/stories"
)

func storiesCMD() *cli.Command {
	cover := &cli.Command{
		Name:        "video",
		Description: `Build cover from URL`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "URL to build cover",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       "Output dir",
				Aliases:     []string{"o"},
				Value:       fmt.Sprintf("outputs/%v-story.mp4", time.Now().Unix()),
				DefaultText: "outputs/{DATE}-story.mp4",
			},
			&cli.StringFlag{
				Name:    "footer-text",
				Usage:   "Footer text",
				Aliases: []string{"ft"},
			},
			&cli.StringFlag{
				Name:    "footer-image",
				Usage:   "Footer image",
				Aliases: []string{"fi"},
			},
		},
		Action: func(cmd *cli.Context) error {
			return actions.VideoStory(cmd.Context, actions.BuildStoryOptions{
				URL:    cmd.String("url"),
				Output: cmd.String("output"),
				Footer: stories.Footer{
					Text:  cmd.String("footer-text"),
					Image: cmd.String("footer-image"),
				},
			})
		},
	}

	return &cli.Command{
		Name:        "stories",
		Description: `Stories related commands`,
		Subcommands: []*cli.Command{cover},
	}
}
