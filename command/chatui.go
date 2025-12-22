package command

import (
	"cli/config"
	"cli/llm/claude"
	"cli/utils"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	messages           []string
	input              textinput.Model
	ExitActive         bool
	suggestions        []Command
	selectedSuggestion int
	suggestionsActive  bool
	Command            tea.Cmd
	configMissing      bool
	waitingForKey      bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Type your message here"
	ti.Width = 60
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	configMissing := !config.ConfigExists()

	messages := []string{
		"Tips for getting started:",
		"1. Ask any question to edit, generate, debug and run commands.",
		"2. Be specific for the best results.",
		"3. type /switch_model <MODEL_NAME>",
	}

	if configMissing {
		messages = []string{
			"⚠ Configuration not found.",
			"Please enter your Claude API key to continue:",
		}
		ti.Placeholder = "Paste your Claude API key here"
	}

	return model{
		messages:      messages,
		input:         ti,
		configMissing: configMissing,
		waitingForKey: configMissing,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) updateSuggestions() {
	if strings.HasPrefix(m.input.Value(), "/") {
		m.suggestionsActive = true
		query := strings.TrimPrefix(m.input.Value(), "/")
		m.suggestions = FilterCommands(query)
		m.selectedSuggestion = 0
	} else {
		m.suggestionsActive = false
		m.suggestions = nil
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = handleKeyMsg(m, msg)
	}
	return m, cmd
}

func handleKeyMsg(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.ExitActive {
			return m, tea.Quit
		}
		m.messages = append(m.messages, "⚠ Press Ctrl+C again to exit.")
		m.ExitActive = true
		return m, nil
	case "up":
		if m.suggestionsActive {
			m.selectedSuggestion = utils.Max(0, m.selectedSuggestion-1)
		}
	case "down":
		if m.suggestionsActive {
			m.selectedSuggestion = utils.Min(len(m.suggestions)-1, m.selectedSuggestion+1)
		}
	case "enter", "tab":
		if m.waitingForKey {
			key := strings.TrimSpace(m.input.Value())
			if key == "" {
				m.messages = append(m.messages, "API key cannot be empty.")
				return m, nil
			}

			if err := config.SaveClaudeKey(key); err != nil {
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
				"3. type /switch_model <MODEL_NAME>",
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
				if strings.HasPrefix(userInput, "/") {
					commandStr := strings.TrimPrefix(userInput, "/")
					result := ExecuteCommand(commandStr)
					m.messages = append(m.messages, "AI: "+result)
				} else {
					m.messages = append(m.messages, "You: "+m.input.Value())
					reply := claude.ClaudeModelReply(m.input.Value())
					m.messages = append(m.messages, "AI: "+reply)
					m.input.SetValue("")
				}
			}
		}
	case "esc":
		m.suggestionsActive = false
		m.suggestions = nil
		m.input.SetValue("")
	}
	return m, nil
}

func (m model) View() string {
	chatBox := strings.Join(m.messages, "\n")
	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(60).
		Render(m.input.View())

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

	return fmt.Sprintf(
		"%s\n\nType your message & press Enter:\n%s\n%s",
		chatBox,
		inputBox,
		suggestionBox,
	)
}

func StartChatUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
