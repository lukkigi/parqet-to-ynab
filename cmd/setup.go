package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lukkigi/parqet-to-ynab/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup your parqet-to-ynab config",
	Long:  `Setup your parqet-to-ynab config`,
	Run: func(cmd *cobra.Command, args []string) {
		result := createModel()

		p := tea.NewProgram(result)

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}

		for i, val := range result.inputs {
			if len(val.Value()) == 0 {
				continue
			}

			if i == 0 {
				viper.Set(config.ParqetPortfolioId, val.Value())
			} else if i == 1 {
				viper.Set(config.YnabApiKey, val.Value())
			} else if i == 2 {
				viper.Set(config.YnabBudgetId, val.Value())
			} else if i == 3 {
				viper.Set(config.YnabInvestingAccountId, val.Value())
			}
		}

		saveNewConfig()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Save ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Save"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func createModel() model {
	newModel := model{
		inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range newModel.inputs {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Parqet Portfolio ID"
			t.CharLimit = 24
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "YNAB API Key"
			t.CharLimit = 32
		case 2:
			t.Placeholder = "YNAB Budget ID"
			t.CharLimit = 32
		case 3:
			t.Placeholder = "YNAB Account ID"
			t.CharLimit = 32
		}

		newModel.inputs[i] = t
	}

	return newModel
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func saveNewConfig() {
	home, _ := os.UserHomeDir()
	var configPath = home + "/.parqet-to-ynab.yaml"

	_, statErr := os.Stat(configPath)
	if !os.IsExist(statErr) {
		if _, createErr := os.Create(configPath); createErr != nil {
			log.Fatal(createErr)
		}
	}

	configErr := viper.WriteConfig()

	if configErr != nil {
		log.Fatal(configErr)
	}
}
