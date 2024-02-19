package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor       lipgloss.Color
	TextColor         lipgloss.Color
	TextHiglightColor lipgloss.Color

	InputField lipgloss.Style
	TextField  lipgloss.Style
}

func (s Styles) Highlight(text string) string {
	s.TextField.Foreground(s.TextHiglightColor)
	highLightText := s.TextField.Render(text)
	s.TextField.Foreground(s.TextColor)
	return highLightText
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.TextColor = lipgloss.Color("32")
	s.TextHiglightColor = lipgloss.Color("36")

	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Width(80).
		Foreground(s.TextColor)
	s.TextField = lipgloss.NewStyle().
		Foreground(s.TextColor)

	return s
}

type model struct {
	items       []string
	searchField textinput.Model
	styles      *Styles
	width       int
	height      int
}

func New(items []string) *model {
	styles := DefaultStyles()
	inputField := textinput.New()
	inputField.Placeholder = "Search Project"
	inputField.Focus()
	return &model{
		items:       items,
		searchField: inputField,
		styles:      styles,
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
	m.items = append(m.items, m.styles.Highlight("testhghmmm"))
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			0.05,
			m.styles.InputField.Render(m.searchField.View()),
			m.styles.TextField.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					m.items...,
				),
			),
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
