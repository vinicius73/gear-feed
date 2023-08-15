package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/apps/cli/actions"
)

func storiesCMD() *cli.Command {
	cover := &cli.Command{
		Name:        "cover",
		Description: `Build cover from URL`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "URL to build cover",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output dir",
				Value: "outputs",
			},
			&cli.StringFlag{
				Name:  "template-filename",
				Usage: "Template filename",
				Value: "{{.date}}-{{.site}}-{{.title}}--{{.filename}}",
			},
		},
		Action: func(cmd *cli.Context) error {
			return actions.Story(cmd.Context, actions.BuildStoryOptions{
				URL:      cmd.String("url"),
				Output:   cmd.String("output"),
				Template: cmd.String("template-filename"),
			})
		},
	}

	return &cli.Command{
		Name:        "stories",
		Description: `Stories related commands`,
		Subcommands: []*cli.Command{cover},
	}
}
