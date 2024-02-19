package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchFieldStyle struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type highlightStyle struct {
	TextColor lipgloss.Color
	TextField lipgloss.Style
}

func DefaultStyles() (*searchFieldStyle, *highlightStyle) {
	s := new(searchFieldStyle)
	t := new(highlightStyle)
	s.BorderColor = lipgloss.Color("36")
	t.TextColor = lipgloss.Color("32")

	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		// Padding(1).
		Width(80).
		Foreground(t.TextColor)
	t.TextField = lipgloss.NewStyle().
		Foreground(s.BorderColor)

	return s, t
}

type model struct {
	items            []string
	searchField      textinput.Model
	searchFieldStyle *searchFieldStyle
	highlightStyle   *highlightStyle
	width            int
	height           int
}

func New(items []string) *model {
	searchFieldStyle, highlightStyle := DefaultStyles()
	inputField := textinput.New()
	inputField.Placeholder = "Search Project"
	inputField.Focus()
	return &model{
		items:            items,
		searchField:      inputField,
		searchFieldStyle: searchFieldStyle,
		highlightStyle:   highlightStyle,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.searchField, cmd = m.searchField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading.."
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			0.05,
			m.searchFieldStyle.InputField.Render(m.searchField.View()),
			m.highlightStyle.TextField.Render(lipgloss.JoinVertical(lipgloss.Left, m.items...)),
		),
	)
}

func main() {
	fuzzyTest()
	items := []string{"test1", "test2loool"}
	m := New(items)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There is an error: %v", err)
		os.Exit(1)
	}
}
