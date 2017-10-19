package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cards",
	Short: "Manage cards and decks",
	Long: `A command line client for the cards service, designed to make
it easy to manage cards and decks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var author string

func init() {

}

func Execute() {
	RootCmd.Execute()
}
