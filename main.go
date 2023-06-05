package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Color lipgloss
var (
	green = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#32CD32"))
	gray  = lipgloss.NewStyle().Foreground(lipgloss.Color("#696969"))
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	help     help.Model
	keys     KeyMap
}

type KeyMap struct {
	UP    key.Binding
	DOWN  key.Binding
	CHECK key.Binding
	QUIT  key.Binding
	HELP  key.Binding
}

var DefaultKeyMap = KeyMap{
	UP: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	DOWN: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	CHECK: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("ENTER/SPACE", "check"),
	),
	QUIT: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	HELP: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.HELP, k.QUIT}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UP, k.DOWN},
		{k.HELP, k.QUIT},
	}
}

func initialModel() model {
	return model{
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi", "Buy milk"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		keys:     DefaultKeyMap,
		help:     help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.QUIT):
			return m, tea.Quit
		case key.Matches(msg, m.keys.UP):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.DOWN):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.CHECK):
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case key.Matches(msg, m.keys.HELP):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		s += fmt.Sprintf(cursor + " ")
		if _, ok := m.selected[i]; ok {
			checked = "x"
			s += "[" + checked + "] "
			s += green.Render(choice)
		} else {
			s += "[" + checked + "] "
			s += gray.Render(choice)
		}
		s += "\n"
	}

	helpview := m.help.View(m.keys)
	height := 2 - strings.Count(helpview, "\n")

	s += strings.Repeat("\n", height) + helpview
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
