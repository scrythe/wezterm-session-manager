package main

import (
	"log"
	"os"

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
	SelectedColor  lipgloss.Color

	InputField         lipgloss.Style
	TextField          lipgloss.Style
	HighlightTextField lipgloss.Style
	SelectedField      lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputTextColor = lipgloss.Color("7")
	s.TextColor = lipgloss.Color("#2daecf")
	s.HighlightColor = lipgloss.Color("#00af87")
	s.SelectedColor = lipgloss.Color("#FFFFFF")

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
	s.SelectedField = lipgloss.NewStyle().
		Background(s.SelectedColor)

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
	cursor           int
	startItem        int
	listLength       int
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

func (m model) getDisplayItems(items []filterItem, start, length int) []filterItem {
	displayItems := items[start : start+length]
	return displayItems
}

func (m model) styleItem(item filterItem, selected bool) string {
	var textStyle lipgloss.Style
	var highlightedTextStyle lipgloss.Style
	if selected {
		textStyle = m.styles.SelectedField.Copy().Inherit(m.styles.TextField)
		highlightedTextStyle = m.styles.SelectedField.Copy().Inherit(m.styles.HighlightTextField)
	} else {
		textStyle = m.styles.TextField
		highlightedTextStyle = m.styles.HighlightTextField
	}
	newItem := lipgloss.StyleRunes(
		item.value,
		item.matches,
		highlightedTextStyle,
		textStyle,
	)
	return newItem
}

func (m model) StyleText(items []filterItem) []string {
	styledItems := make([]string, len(items))
	for i, item := range items {
		selected := i == m.cursor-m.startItem
		styledItems[i] = m.styleItem(item, selected)
	}
	return styledItems
}

func (m model) StyleList() string {
	listEnd := min(len(m.filteredItems), m.startItem+m.listLength)
	displayItems := m.filteredItems[m.startItem:listEnd]
	// log.Print(displayItems)
	styleItems := m.StyleText(displayItems)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styleItems...,
	)
}

func New(items []string) *model {
	styles := DefaultStyles()
	inputField := textinput.New()
	inputField.Placeholder = "Search Project"
	inputField.Focus()
	return &model{
		items:       items,
		listLength:  42,
		searchField: inputField,
		styles:      styles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) CursorDown() {
	m.cursor++
	itemsLength := len(m.filteredItems)
	if m.cursor >= itemsLength {
		m.cursor = itemsLength - 1
	}
	if m.cursor >= m.startItem+m.listLength {
		m.startItem++
	}
}

func (m *model) CursorUp() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.startItem {
		m.startItem--
	}
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
		case "ctrl+j", "ctrl+n":
			m.CursorDown()
		case "ctrl+k", "ctrl+p":
			m.CursorUp()
		}
	}
	m.searchField, cmd = m.searchField.Update(msg)
	m.UpdateSearch()
	return m, cmd
}

func (m *model) CheckCursor() {
	itemsLength := len(m.filteredItems)
	if m.cursor >= itemsLength {
		m.cursor = itemsLength - 1
	} else if m.cursor < 0 {
		// list length 0 -> cursor -1
		m.cursor = 0
	}
}

func (m *model) CheckStartItem() {
	itemsLength := len(m.filteredItems)
	if m.startItem+m.listLength > itemsLength {
		m.startItem = 0
	}
}

func (m *model) UpdateSearch() {
	input := m.searchField.Value()
	if input == "" {
		m.filteredItems = getFullList(m.items)
	} else if m.searchFieldValue != input {
		m.filteredItems = getFilteredItems(input, m.items)
	}
	m.CheckCursor()
  m.CheckStartItem()
	m.searchFieldValue = input
}

func (m model) View() string {
	if m.width == 0 {
		return "loading.."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			0.05,
			m.styles.InputField.Render(m.searchField.View()),
			m.StyleList(),
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
