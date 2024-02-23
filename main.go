package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Styles struct {
	BorderColor    lipgloss.Color
	InputTextColor lipgloss.Color
	TextColor      lipgloss.Color
	HighlightColor lipgloss.Color
	// SelectedColor  lipgloss.Color

	InputField         lipgloss.Style
	TextField          lipgloss.Style
	HighlightTextField lipgloss.Style
	// SelectedField      lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputTextColor = lipgloss.Color("7")
	s.TextColor = lipgloss.Color("#2daecf")
	s.HighlightColor = lipgloss.Color("#00af87")
	// s.SelectedColor = lipgloss.Color("34")

	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Width(80).
		Foreground(s.InputTextColor)
	s.TextField = lipgloss.NewStyle().
		Foreground(s.TextColor)
	s.HighlightTextField = lipgloss.NewStyle().
		Foreground(s.HighlightColor).Underline(true)
	// s.SelectedField = lipgloss.NewStyle().
	// 	Background(s.SelectedColor)

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
	// cursor           int
	width            int
	height           int
}

func getFullList(items []string) []filterItem {
	newItems := make([]filterItem, len(items))
	for i, item := range items {
		newItems[i] = filterItem{value: item}
	}
	return newItems
}

func getFilteredItems(input string, items []string) []filterItem {
	matches := fuzzy.Find(input, items)

	filteredItems := make([]filterItem, len(matches))
	for i, item := range matches {
		filteredItems[i] = filterItem{
			value:   item.Str,
			matches: item.MatchedIndexes,
		}
	}
	return filteredItems
}

func (m model) styleItem(item filterItem) string {
	newItem := lipgloss.StyleRunes(
		item.value,
		item.matches,
		m.styles.HighlightTextField,
		m.styles.TextField,
	)
	// if !selected {
	// 	return m.styles.SelectedField.Render(newItem)
	// }
	return newItem
}

func (m model) StyleText() []string {
	styledItems := make([]string, len(m.filteredItems))
	for i, item := range m.filteredItems {
		// selected := i == m.cursor
		styledItems[i] = m.styleItem(item)
	}
	return styledItems
}

func (m model) joinList() string {
	var b strings.Builder
	styledList := m.StyleText()
	for _, item := range styledList {
		fmt.Fprintf(&b, "%v\n", item)
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
	m.UpdateSearch()
	return m, cmd
}

func (m *model) UpdateSearch() {
	input := m.searchField.Value()
	if m.searchFieldValue != input {
		m.searchFieldValue = input
		m.filteredItems = getFilteredItems(input, m.items)
	}
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
	// text:=getText()
	return lipgloss.JoinVertical(
		0.05,
		m.styles.InputField.Render(m.searchField.View()),
		// m.joinList(),
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.StyleText()...,
		),
	)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("fatal: %v", err)
		os.Exit(1)
	}
	folders := listFolders()
	m := New(folders)
	p := tea.NewProgram(m, tea.WithAltScreen())
	defer f.Close()
	if _, err := p.Run(); err != nil {
		log.Fatalf("There is an error: %v", err)
		os.Exit(1)
	}
}
