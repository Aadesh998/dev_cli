package main

import (
	"cli/internal/app/ui"
	"cli/internal/command"
	"cli/internal/config"
)

func main() {
	config.LoadConfig()
	ui.Poster()
	command.Execute()
	ui.StartChatUI()
}
