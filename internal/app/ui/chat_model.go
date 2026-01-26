package ui

import (
	"cli/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Command struct {
	Name        string
	Description string
}

type chatModel struct {
	messages           []string
	input              textinput.Model
	ExitActive         bool
	suggestions        []Command
	selectedSuggestion int
	suggestionsActive  bool
	configMissing      bool
	waitingForKey      bool
	choices            []string
	cursor             int
	selected           map[int]struct{}
	showChoices        bool
}

func newChatModel() chatModel {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Type your message here"
	ti.Width = 60
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	configMissing := !config.ConfigExists()
	m := chatModel{
		input:         ti,
		configMissing: configMissing,
		waitingForKey: configMissing,
	}

	if configMissing {
		m.messages = []string{"Configuration not found.", "Please enter your Claude API key:"}
		m.input.Placeholder = "Paste your Claude API key here"

		m.choices = []string{
			"claude",
			"kimik2",
			"openAI",
		}
		m.selected = make(map[int]struct{})
		m.showChoices = true
	} else {
		m.messages = []string{"1. Ask questions", "2. Be specific", "3. /switch_model <NAME>"}
	}
	return m
}
