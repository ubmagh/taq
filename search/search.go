package search

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ubmagh/taq/types"
)

type SearchModel struct {
	input textinput.Model
	list  list.Model
	hosts []types.Host
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m SearchModel) Init() tea.Cmd { return textinput.Blink }

func toListItems(hosts []types.Host) []list.Item {
	items := []list.Item{}
	for _, h := range hosts {
		items = append(items, item(fmt.Sprintf("%s (%s)", h.Name, h.Address)))
	}
	return items
}

func (m *SearchModel) filterList() {
	query := strings.ToLower(m.input.Value())
	filtered := []types.Host{}
	for _, h := range m.hosts {
		if strings.Contains(strings.ToLower(h.Name), query) ||
			strings.Contains(strings.ToLower(h.Address), query) {
			filtered = append(filtered, h)
		}
	}
	m.list.SetItems(toListItems(filtered))
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.filterList()
	return m, cmd
}

func (m SearchModel) View() string {
	return fmt.Sprintf("Search by keywords: %20s\n\n%s", m.input.View(), m.list.View())
}

func NewSearcher(hosts []types.Host) SearchModel {
	items := toListItems(hosts)
	ti := textinput.New()
	ti.PlaceholderStyle.Blink(true).Width(1)
	ti.Placeholder = "Type to search..."
	ti.Width = 30
	ti.CharLimit = 200
	ti.Focus()

	l := list.New(items, itemDelegate{}, 0, 10)
	l.SetShowHelp(true)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SearchModel{ti, l, hosts}
}

func RunSearcher(hosts []types.Host) {
	p := tea.NewProgram(NewSearcher(hosts))
	if err := p.Start(); err != nil {
		panic(err)
	}
}
