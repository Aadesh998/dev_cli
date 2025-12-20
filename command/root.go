package command

import (
	"cli/ui"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "Simple CLI chat app",

	Run: func(cmd *cobra.Command, args []string) {
		ui.Poster()
		StartChatUI()
	},
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
