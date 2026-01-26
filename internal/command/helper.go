package command

import (
	"strings"

	"github.com/spf13/cobra"
)

type Command struct {
	Name        string
	Description string
}

func allCommands() []Command {
	cmds := GetCommands()
	var commands []Command
	for _, c := range cmds {
		commands = append(commands, Command{
			Name:        c.Use,
			Description: c.Short,
		})
	}
	return commands
}

func FilterCommands(input string) []Command {
	var filtered []Command
	for _, cmd := range allCommands() {
		if strings.HasPrefix(strings.ToLower(cmd.Name), strings.ToLower(input)) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

func GetCommands() []*cobra.Command {
	return rootCmd.Commands()
}
