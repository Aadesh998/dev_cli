package ui

import (
	"cli/internal/app/chat"
	"cli/internal/command"
	"cli/internal/config"
	"cli/internal/llm/claude"
	"cli/internal/llm/kimik2"
	"cli/internal/llm/openai"
	"cli/internal/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m chatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *chatModel) updateSuggestions() {
	if strings.HasPrefix(m.input.Value(), "/") {
		m.suggestionsActive = true
		query := strings.TrimPrefix(m.input.Value(), "/")
		filteredCommands := command.FilterCommands(query)
		m.suggestions = make([]Command, len(filteredCommands))
		for i, cmd := range filteredCommands {
			m.suggestions[i] = Command(cmd)
		}
		m.selectedSuggestion = 0
	} else {
		m.suggestionsActive = false
		m.suggestions = nil
	}
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !m.showChoices {
		m.input, cmd = m.input.Update(msg)
		m.updateSuggestions()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = handlekeyPress(m, msg)
	}
	return m, cmd
}

func handlekeyPress(m chatModel, msg tea.KeyMsg) (chatModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.ExitActive {
			return m, tea.Quit
		}
		m.messages = append(m.messages, "⚠ Press Ctrl+C again to exit.")
		m.ExitActive = true
		return m, nil

	case "up":
		if m.showChoices {
			if m.cursor > 0 {
				m.cursor--
			}
		} else if m.suggestionsActive {
			m.selectedSuggestion = utils.Max(0, m.selectedSuggestion-1)
		}

	case "down":
		if m.showChoices {
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		} else if m.suggestionsActive {
			m.selectedSuggestion = utils.Min(len(m.suggestions)-1, m.selectedSuggestion+1)
		}

	case "enter", "tab":
		if m.showChoices {
			selectedModel := m.choices[m.cursor]
			m.messages = append(m.messages, fmt.Sprintf("Selected model: %s", selectedModel))
			m.messages = append(m.messages, "Please enter your API key:")

			m.showChoices = false
			m.waitingForKey = true
			m.input.Placeholder = "Enter your API key"
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		}

		if m.waitingForKey {
			key := strings.TrimSpace(m.input.Value())
			if key == "" {
				m.messages = append(m.messages, "API key cannot be empty.")
				return m, nil
			}

			if err := config.SaveApiKey(m.choices[m.cursor], key); err != nil {
				m.messages = append(m.messages, "Failed to save config.")
				return m, nil
			}
			if err := config.LoadConfig(); err != nil {
				fmt.Println("Config error:", err)
				os.Exit(1)
			}

			m.messages = append(m.messages, "API key saved successfully.")
			messages := []string{
				"\n",
				"Tips for getting started:",
				"1. Ask any question to edit, generate, debug and run commands.",
				"2. Be specific for the best results.",
				"3. type /switch <MODEL_NAME> <APIKEY>",
			}

			m.messages = append(m.messages, messages...)
			m.waitingForKey = false
			m.configMissing = false
			m.input.SetValue("")
			m.input.Placeholder = "Type your message here"
			return m, nil
		}

		if m.suggestionsActive && len(m.suggestions) > 0 {
			m.input.SetValue("/" + m.suggestions[m.selectedSuggestion].Name)
			m.suggestionsActive = false
			m.suggestions = nil
		} else {
			if m.input.Value() != "" {
				userInput := m.input.Value()
				if after, ok := strings.CutPrefix(userInput, "/"); ok {
					commandStr := after
					result := command.ExecuteCommand(commandStr)
					m.messages = append(m.messages, "AI: "+result)
				} else {
					log.Printf("Model Using:- %s", config.LlmClient.ModelName)

					m.messages = append(m.messages, "You: "+m.input.Value())
					var provider chat.ChatProvider

					switch config.LlmClient.ModelName {
					case utils.ModelClaude:
						provider = claude.ClaudeProvider{}
					case utils.ModelKimik2:
						provider = kimik2.Kimik2Provider{}
					case utils.ModelOpenai:
						provider = openai.OpenaiProvider{}
					}
					reply, err := provider.ChatProcess(m.input.Value())
					if err != nil {
						log.Printf("Failed to get the response from the AI: %s", err)
					}

					m.messages = append(m.messages, "AI: "+reply.Text)
					m.input.SetValue("")
				}
			}
		}

	case "esc":
		if m.showChoices {
			m.showChoices = false
			m.messages = append(m.messages, "Model selection cancelled.")
		} else {
			m.suggestionsActive = false
			m.suggestions = nil
			m.input.SetValue("")
		}

	case " ":
		if m.showChoices {
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m chatModel) View() string {
	chatBox := strings.Join(m.messages, "\n")

	var choicesBox string
	if m.showChoices {
		var s strings.Builder
		s.WriteString("Select a model (use ↑↓ to navigate, Enter to select, Esc to cancel):\n\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice))
		}
		choicesBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(60).
			Render(s.String())
	}

	var inputBox string
	if !m.showChoices {
		inputBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(60).
			Render(m.input.View())
	}

	var suggestionBox string
	if m.suggestionsActive && len(m.suggestions) > 0 {
		var s strings.Builder
		for i, sug := range m.suggestions {
			if i == m.selectedSuggestion {
				s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("> " + sug.Name + " : " + sug.Description))
			} else {
				s.WriteString("  " + sug.Name + " : " + sug.Description)
			}
			s.WriteString("\n")
		}
		suggestionBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(60).
			Render(s.String())
	}

	var output strings.Builder
	output.WriteString(chatBox)
	output.WriteString("\n\n")

	if choicesBox != "" {
		output.WriteString(choicesBox)
		output.WriteString("\n")
	} else {
		output.WriteString("Type your message & press Enter:\n")
		output.WriteString(inputBox)

		if suggestionBox != "" {
			output.WriteString("\n")
			output.WriteString(suggestionBox)
		}
	}

	return output.String()
}

func StartChatUI() {
	p := tea.NewProgram(newChatModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
