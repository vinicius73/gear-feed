package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg"
)

func main() {
	app := &cli.App{
		Name:  "gamerfeed",
		Usage: "Gamer Feed Bot CLI",
		Commands: []*cli.Command{
			listCMD(),
		},
		Action: func(*cli.Context) error {
			fmt.Println(pkg.VersionVerbose())

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
