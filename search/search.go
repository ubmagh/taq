package search

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
	"github.com/ubmagh/taq/host"
	"github.com/ubmagh/taq/ui"
)

type phase int

const (
	phaseSearch phase = iota
	phaseUser
)

type SearchModel struct {
	phase        phase
	input        textinput.Model
	userInput    textinput.Model
	list         list.Model
	hosts        []host.Host
	selectedHost host.Host
	width        int
	height       int
}

var (
	titleStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).MarginLeft(2)
	itemNameStyle       = lipgloss.NewStyle().PaddingLeft(4)
	itemDetailStyle     = lipgloss.NewStyle().PaddingLeft(4).Faint(true)
	selectedNameStyle   = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.Color("170"))
	selectedDetailStyle = lipgloss.NewStyle().PaddingLeft(2).Faint(true).Foreground(lipgloss.Color("170"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle           = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item struct {
	host host.Host
}

func (i item) Title() string       { return i.host.Name }
func (i item) Description() string { return i.host.Address }
func (i item) FilterValue() string { return i.host.HostListDisplay() }

func itemDetail(h host.Host) string {
	parts := []string{h.Address}
	if h.User != "" {
		parts = append(parts, h.User)
	}
	if g := h.Labels["groups"]; g != "" {
		parts = append(parts, g)
	} else if g := h.Labels["groupName"]; g != "" {
		parts = append(parts, g)
	}
	return strings.Join(parts, " · ")
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	detail := itemDetail(i.host)
	if index == m.Index() {
		fmt.Fprintln(w, selectedNameStyle.Render("> "+i.host.Name))
		fmt.Fprint(w, selectedDetailStyle.Render("  "+detail))
	} else {
		fmt.Fprintln(w, itemNameStyle.Render(i.host.Name))
		fmt.Fprint(w, itemDetailStyle.Render(detail))
	}
}

func (m SearchModel) Init() tea.Cmd { return textinput.Blink }

func toListItems(hosts []host.Host) []list.Item {
	items := make([]list.Item, 0, len(hosts))
	for _, h := range hosts {
		items = append(items, item{host: h})
	}
	return items
}

func (m *SearchModel) filterList() {
	query := strings.TrimSpace(m.input.Value())
	if query == "" {
		m.list.SetItems(toListItems(m.hosts))
		return
	}

	lq := strings.ToLower(query)

	searchables := make([]string, len(m.hosts))
	for i, h := range m.hosts {
		searchables[i] = h.Searchable()
	}

	matches := fuzzy.Find(lq, searchables)

	seen := make(map[int]bool, len(matches))
	filtered := make([]host.Host, 0, len(matches))
	for _, match := range matches {
		seen[match.Index] = true
		filtered = append(filtered, m.hosts[match.Index])
	}

	// fuzzy handles names/groups well but struggles with IPs — supplement with address substring match
	for i, h := range m.hosts {
		if !seen[i] && strings.Contains(strings.ToLower(h.Address), lq) {
			filtered = append(filtered, h)
		}
	}

	m.list.SetItems(toListItems(filtered))
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-5)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.phase == phaseUser {
				m.phase = phaseSearch
				m.userInput.Blur()
				m.input.Focus()
				return m, nil
			}
			return m, tea.Quit
		case "down":
			if m.phase == phaseSearch {
				m.list.CursorDown()
			}
			return m, nil
		case "up":
			if m.phase == phaseSearch {
				m.list.CursorUp()
			}
			return m, nil
		case "enter":
			if m.phase == phaseSearch {
				if selected, ok := m.list.SelectedItem().(item); ok {
					m.selectedHost = selected.host
					m.phase = phaseUser
					m.input.Blur()
					m.userInput.SetValue("")
					m.userInput.Placeholder = m.selectedHost.User
					return m, m.userInput.Focus()
				}
			} else {
				user := strings.TrimSpace(m.userInput.Value())
				if user != "" {
					m.selectedHost.User = user
				}
				return m, tea.Sequence(tea.ClearScreen, tea.Quit)
			}
		}
	}

	if m.phase == phaseSearch {
		m.input, cmd = m.input.Update(msg)
		m.filterList()
	} else {
		m.userInput, cmd = m.userInput.Update(msg)
	}
	return m, cmd
}

func (m SearchModel) View() string {
	if m.phase == phaseUser {
		help := lipgloss.NewStyle().Faint(true).Render("`Enter` confirm • `Esc` back")
		return fmt.Sprintf("SSH username for [%s]: %s\n%s", m.selectedHost.Name, m.userInput.View(), help)
	}

	help := lipgloss.NewStyle().Faint(true).Render("`↑/↓` navigate • `Enter` connect • `Esc/Ctrl+C` exit")
	return fmt.Sprintf("Search: %s\n%s%s", m.input.View(), m.list.View(), help)
}

func NewSearcher(hosts []host.Host) SearchModel {
	items := toListItems(hosts)

	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Width = 30
	ti.CharLimit = 200
	ti.Focus()

	l := list.New(items, itemDelegate{}, 0, 20)
	l.SetShowHelp(false)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.Styles.Title = titleStyle
	l.Title = "Target instances"
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	uInput := textinput.New()
	uInput.Width = 30
	uInput.CharLimit = 100

	return SearchModel{
		phase:     phaseSearch,
		input:     ti,
		userInput: uInput,
		list:      l,
		hosts:     hosts,
	}
}

func RunSearcher(hosts []host.Host) (host.Host, bool) {
	p := tea.NewProgram(NewSearcher(hosts))
	model, err := p.Run()
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}
	if sm, ok := model.(SearchModel); ok && sm.selectedHost.Address != "" {
		return sm.selectedHost, true
	}
	return host.Host{}, false
}
