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

var switchModelCmd = &cobra.Command{
	Use:   "switch [model_name] [apikey]",
	Short: "Switch model and API key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		apikey := args[1]
		config.SaveApiKey(modelName, apikey)
		config.LoadConfig()
		fmt.Fprintf(cmd.OutOrStdout(), "Model Switch Successfully")
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello [name]",
	Short: "Prints a greeting",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Hello World %s\n", args[0])
	},
}

func init() {
	initLogger()
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)

	rootCmd.AddCommand(helloCmd)
	rootCmd.AddCommand(switchModelCmd)
}

func initLogger() {
	file, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
