package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gear-feed/pkg/sources"
	layout "github.com/vinicius73/gear-feed/pkg/tui/sources"
)

func sourcesCMD() *cli.Command {
	list := &cli.Command{
		Name: "list",
		Flags: []cli.Flag{
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
		},
		Action: func(cmd *cli.Context) error {
			list, err := sources.Load(cmd.Context, sources.LoadOptions{
				Only:  cmd.StringSlice("only"),
				Paths: cmd.StringSlice("sources"),
			})
			if err != nil {
				return err
			}

			logger := zerolog.Ctx(cmd.Context)

			ctx := logger.WithContext(cmd.Context)

			p := tea.NewProgram(
				layout.NewModel(ctx, list),
				tea.WithAltScreen(),
				tea.WithContext(ctx),
			)
			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return &cli.Command{
		Name:        "sources",
		Description: "Interact with sources",
		Subcommands: []*cli.Command{list},
	}
}
