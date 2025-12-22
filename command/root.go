package command

import (
	"cli/config"
	"cli/ui"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "Simple CLI chat app",

	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		ui.Poster()
		StartChatUI()
	},
}

var helloCmd = &cobra.Command{
	Use:   "Hello",
	Short: "Hello Command",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Hello World")
	},
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)

	rootCmd.AddCommand(helloCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
