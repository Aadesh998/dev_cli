package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help command used for CLI-related available commands and help",
	Run: func(cmd *cobra.Command, args []string) {
		helpList := getHelp()
		for _, h := range helpList {
			fmt.Println(h)
		}
	},
}

func getHelp() []string {
	help := []string{
		"1. Help",
		"2. Dev CLI",
	}
	return help
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
