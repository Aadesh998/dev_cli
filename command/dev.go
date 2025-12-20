package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev command for CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dev command")
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
