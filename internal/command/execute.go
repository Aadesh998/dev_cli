package command

import (
	"bytes"
	"fmt"
	"strings"
)

func ExecuteCommand(commandStr string) string {
	args := strings.Fields(commandStr)
	if len(args) == 0 {
		return "No command provided"
	}

	cmd, _, err := rootCmd.Find(args)
	if err != nil {
		return fmt.Sprintf("Command not found: %s", args[0])
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	rootCmd.SetArgs(args)

	err = cmd.Execute()
	if err != nil {
		return fmt.Sprintf("Error executing command: %v", err)
	}

	output := buf.String()
	if output == "" {
		output = "Command executed successfully (no output)"
	}

	return output
}
