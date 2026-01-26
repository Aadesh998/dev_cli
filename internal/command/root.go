package command

import (
	"cli/internal/config"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "Simple CLI chat app",
}

var switchModelCommand = &cobra.Command{
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

var helloCommand = &cobra.Command{
	Use:   "hello [name]",
	Short: "Prints a greeting",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Hello World %s\n", args[0])
	},
}

func init() {
	initializeLogger()
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.AddCommand(helloCommand)
	rootCmd.AddCommand(switchModelCommand)
}

func initializeLogger() {
	file, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
