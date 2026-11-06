package wizard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices []string
	cursor  int
	choice  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter", "space":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := strings.Builder{}
	for i := range m.choices {
		if m.cursor == i {
			s.WriteString(">>>")
		} else {
			s.WriteString(" ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n[select] ↑/↓ 	[enter] enter,space\n\n")
	return tea.NewView(s.String())
}

// return index, error
func Ask(choices []string) (int, error) {

	if len(choices) == 0 {
		return 0, fmt.Errorf("no choices provided")
	}

	initialModel := model{
		choices: choices,
	}
	p := tea.NewProgram(initialModel)

	m, err := p.Run()
	if err != nil {
		return 0, err
	}

	r, ok := m.(model)
	if !ok {
		return 0, fmt.Errorf("tea.Model is not model")
	}

	return r.cursor, nil
}

func AskYesOrNo() bool {
	initialModel := model{
		choices: []string{"yes", "no"},
	}
	p := tea.NewProgram(initialModel)

	m, _ := p.Run()

	r, ok := m.(model)
	if !ok {
		fmt.Println("tea.Modelが model 型ではありません")
		return false
	}

	switch r.choice {
	case "yes":
		return true
	default:
		return false
	}
}
