package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Styles struct {
	BorderColor        lipgloss.Color
	TextColor          lipgloss.Color
	HighlightTextColor lipgloss.Color

	InputField         lipgloss.Style
	TextField          lipgloss.Style
	HighlightTextField lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.TextColor = lipgloss.Color("32")
	s.HighlightTextColor = lipgloss.Color("36")

	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Width(80).
		Foreground(s.TextColor)
	s.TextField = lipgloss.NewStyle().
		Foreground(s.TextColor)
	s.HighlightTextField = lipgloss.NewStyle().
		Foreground(s.HighlightTextColor)

	return s
}

type filterItem struct {
	value   string
	matches []int
}

type model struct {
	items            []string
	filteredItems    []filterItem
	searchField      textinput.Model
	searchFieldValue string
	styles           *Styles
	width            int
	height           int
}

func fullList(items []string) []filterItem {
	newItems := make([]filterItem, len(items))
	for i, item := range items {
		newItems[i] = filterItem{value: item}
	}
	return newItems
}

func (m *model) FuzzyFind(searchFieldValue string) {
	matches := fuzzy.Find(searchFieldValue, m.items)

	if len(matches) == 0 {
		m.filteredItems = fullList(m.items)
		return
	}

	filteredItems := make([]filterItem, len(matches))
	for i, item := range matches {
		filteredItems[i] = filterItem{
			value:   item.Str,
			matches: item.MatchedIndexes,
		}
	}
	m.filteredItems = filteredItems
}

func (m model) StyleText() []string {
	styledItems := make([]string, len(m.filteredItems))
	for i, item := range m.filteredItems {
		styledItems[i] = lipgloss.StyleRunes(
			item.value,
			item.matches,
			m.styles.HighlightTextField,
			m.styles.TextField,
		)
	}
	return styledItems
}

func styleText(w io.Writer, item filterItem, m model) {
	newItem := lipgloss.StyleRunes(
		item.value,
		item.matches,
		m.styles.HighlightTextField,
		m.styles.TextField,
	)
	fmt.Fprintf(w, "%s\n", newItem)
}

func (m model) RenderText() string {
	var b strings.Builder
	for _, item := range m.filteredItems {
		styleText(&b, item, m)
	}
	return b.String()
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
	searchFieldValue := m.searchField.Value()
	if m.searchFieldValue != searchFieldValue {
		m.searchFieldValue = searchFieldValue
		m.FuzzyFind(searchFieldValue)
	}
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading.."
	}
	// lipgloss.Place(
	// 		m.width,
	// 		m.height,
	// 		lipgloss.Center,
	// 		lipgloss.Top,
	//
	// 			m.styles.InputField.Render(m.searchField.View())
	//       )
	// (m.RenderText())
	return lipgloss.JoinVertical(
		0.05,
		m.styles.InputField.Render(m.searchField.View()),
		m.RenderText(),
		// lipgloss.JoinVertical(
		// 	lipgloss.Left,
		// 	m.StyleText()...,
		// ),
	)
}

func main() {
	folders := listFolders()
	m := New(folders)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There is an error: %v", err)
		os.Exit(1)
	}
}
