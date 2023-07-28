package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

func main() {
	support.SetupLogger("DEBUG", "text", nil)

	app := &cli.App{
		Name:  "gamerfeed",
		Usage: "Gamer Feed Bot CLI",
		Commands: []*cli.Command{
			sourcesCMD(),
			scrapCMD(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
