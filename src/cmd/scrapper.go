package cmd

import (
	"fmt"
	"gfeed/domains/scrappers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:    "scrapper",
		Short:  "Run scrappers",
		PreRun: versionHook,
		Run: func(cmd *cobra.Command, args []string) {
			entries := scrappers.NewsEntries()

			for _, v := range entries {
				fmt.Println("--")
				fmt.Println(v)
			}
		},
	})
}
