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
	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
	"github.com/ubmagh/taq/ui"
)

type phase int

const (
	phaseSearch phase = iota
	phaseUser
	phasePortForward // rules collection — only reached when mode != KindSSH
)

// ResultKind describes the action the user is performing.
type ResultKind int

const (
	KindSSH ResultKind = iota
	KindLocalForward
	KindRemoteForward
)

// Result is returned by RunSearcher.
type Result struct {
	Host  host.Host
	Kind  ResultKind
	Rules []string // only set for KindLocalForward / KindRemoteForward
}

type SearchModel struct {
	phase        phase
	input        textinput.Model
	userInput    textinput.Model
	list         list.Model
	hosts        []host.Host
	searchables  []string // pre-computed parallel to hosts; never rebuilt after init
	selectedHost host.Host
	width        int
	height       int
	compact      bool
	mode         ResultKind // set from CLI flag; KindSSH = normal SSH session

	// port-forward state — only used when mode != KindSSH
	pfRuleInput textinput.Model
	pfRules     []string
}

var (
	titleStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).MarginLeft(2)
	itemNameStyle       = lipgloss.NewStyle().PaddingLeft(4)
	itemDetailStyle     = lipgloss.NewStyle().PaddingLeft(4).Faint(true)
	selectedNameStyle   = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.Color("170"))
	selectedDetailStyle = lipgloss.NewStyle().PaddingLeft(2).Faint(true).Foreground(lipgloss.Color("170"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle           = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	faintStyle          = lipgloss.NewStyle().Faint(true)
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

type itemDelegate struct {
	compact bool
}

func (d itemDelegate) Height() int {
	if d.compact {
		return 1
	}
	return 2
}
func (d itemDelegate) Spacing() int {
	if d.compact {
		return 0
	}
	return 1
}
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	if d.compact {
		if index == m.Index() {
			fmt.Fprint(w, selectedNameStyle.Render("> "+i.host.HostListDisplay()))
		} else {
			fmt.Fprint(w, itemNameStyle.Render(i.host.HostListDisplay()))
		}
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
	matches := fuzzy.Find(lq, m.searchables)

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

		case "tab":
			if m.phase == phaseSearch {
				m.compact = !m.compact
				m.list.SetDelegate(itemDelegate{compact: m.compact})
			}
			return m, nil

		case "esc":
			switch m.phase {
			case phaseUser:
				m.phase = phaseSearch
				m.userInput.Blur()
				return m, m.input.Focus()
			case phasePortForward:
				// Go all the way back to search so the user can pick a different host.
				m.phase = phaseSearch
				m.pfRules = nil
				m.pfRuleInput.SetValue("")
				m.pfRuleInput.Blur()
				return m, m.input.Focus()
			default:
				return m, tea.Quit
			}

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
			switch m.phase {
			case phaseSearch:
				if selected, ok := m.list.SelectedItem().(item); ok {
					m.selectedHost = selected.host
					m.phase = phaseUser
					m.input.Blur()
					m.userInput.SetValue("")
					m.userInput.Placeholder = m.selectedHost.User
					return m, m.userInput.Focus()
				}

			case phaseUser:
				user := strings.TrimSpace(m.userInput.Value())
				if user != "" {
					m.selectedHost.User = user
				}
				if m.mode != KindSSH {
					// Forward mode: collect rules next.
					m.phase = phasePortForward
					m.pfRules = nil
					m.pfRuleInput.SetValue("")
					m.userInput.Blur()
					return m, m.pfRuleInput.Focus()
				}
				return m, tea.Sequence(tea.ClearScreen, tea.Quit)

			case phasePortForward:
				rule := strings.TrimSpace(m.pfRuleInput.Value())
				if rule != "" {
					m.pfRules = append(m.pfRules, rule)
					m.pfRuleInput.SetValue("")
					return m, nil
				}
				if len(m.pfRules) > 0 {
					// Empty input + rules present → open tunnel.
					return m, tea.Sequence(tea.ClearScreen, tea.Quit)
				}
				// Empty input + no rules → stay, need at least one rule.
				return m, nil
			}
		}
	}

	// Route remaining messages to the active input.
	switch m.phase {
	case phaseSearch:
		m.input, cmd = m.input.Update(msg)
		m.filterList()
	case phaseUser:
		m.userInput, cmd = m.userInput.Update(msg)
	case phasePortForward:
		m.pfRuleInput, cmd = m.pfRuleInput.Update(msg)
	}
	return m, cmd
}

func (m SearchModel) View() string {
	switch m.phase {
	case phaseUser:
		action := "SSH"
		if m.mode != KindSSH {
			action = "Port Forward"
		}
		help := faintStyle.Render("`Enter` confirm  •  `Esc` back")
		return fmt.Sprintf("%s username for [%s]: %s\n%s", action, m.selectedHost.Name, m.userInput.View(), help)

	case phasePortForward:
		kindLabel := "Local (-L)"
		if m.mode == KindRemoteForward {
			kindLabel = "Remote (-R)"
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Port Forwarding [%s] — %s (%s)\n\n", kindLabel, m.selectedHost.Name, m.selectedHost.User)
		if len(m.pfRules) > 0 {
			sb.WriteString("  Rules:\n")
			for _, r := range m.pfRules {
				fmt.Fprintf(&sb, "    ✓ %s\n", r)
			}
			sb.WriteByte('\n')
		}
		rulePrompt := "  Add rule"
		if m.mode == KindLocalForward {
			rulePrompt += " (your port → server port)"
		} else {
			rulePrompt += " (server port → your port)"
		}
		fmt.Fprintf(&sb, "%s: %s\n\n", rulePrompt, m.pfRuleInput.View())
		hint := "`Enter` add rule"
		if len(m.pfRules) > 0 {
			hint += "  •  empty `Enter` to connect"
		}
		hint += "  •  `Esc` back"
		sb.WriteString(faintStyle.Render(hint))
		return sb.String()

	default: // phaseSearch
		help := faintStyle.Render("`↑/↓` navigate  •  `Enter` select  •  `Tab` toggle view  •  `Esc/Ctrl+C` exit")
		return fmt.Sprintf("Search: %s\n%s%s", m.input.View(), m.list.View(), help)
	}
}

func NewSearcher(hosts []host.Host, mode ResultKind) SearchModel {
	items := toListItems(hosts)

	// Pre-compute searchable strings once; filterList reuses this slice on every keystroke.
	searchables := make([]string, len(hosts))
	for i, h := range hosts {
		searchables[i] = h.Searchable()
	}

	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Width = 30
	ti.CharLimit = 200
	ti.Focus()

	compact := config.IsCompactMode()
	l := list.New(items, itemDelegate{compact: compact}, 0, 20)
	l.SetShowHelp(false)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	switch mode {
	case KindLocalForward:
		l.Title = "Target instances [Local Forward -L]"
	case KindRemoteForward:
		l.Title = "Target instances [Remote Forward -R]"
	default:
		l.Title = "Target instances"
	}

	uInput := textinput.New()
	uInput.Width = 30
	uInput.CharLimit = 100

	pfRuleInput := textinput.New()
	pfRuleInput.Width = 30
	pfRuleInput.CharLimit = 100
	switch mode {
	case KindLocalForward:
		pfRuleInput.Placeholder = "your_port->server_port"
	case KindRemoteForward:
		pfRuleInput.Placeholder = "server_port->your_port"
	default:
		pfRuleInput.Placeholder = "8080->3000"
	}

	return SearchModel{
		phase:       phaseSearch,
		input:       ti,
		userInput:   uInput,
		list:        l,
		hosts:       hosts,
		searchables: searchables,
		compact:     compact,
		mode:        mode,
		pfRuleInput: pfRuleInput,
	}
}

func RunSearcher(hosts []host.Host, mode ResultKind) (Result, bool) {
	p := tea.NewProgram(NewSearcher(hosts, mode))
	model, err := p.Run()
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}
	if sm, ok := model.(SearchModel); ok && sm.selectedHost.Address != "" {
		switch sm.phase {
		case phaseUser:
			return Result{Host: sm.selectedHost, Kind: KindSSH}, true
		case phasePortForward:
			return Result{Host: sm.selectedHost, Kind: sm.mode, Rules: sm.pfRules}, true
		}
	}
	return Result{}, false
}
