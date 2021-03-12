package cmd

import (
	"gfeed/domains/bot"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "bot",
		Short: "Run bot client",
		Run: func(cmd *cobra.Command, args []string) {
			bot.Agent(getBotConfig())
		},
	})
}
