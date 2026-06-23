package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/omangm/dwaar/internal/config"
	"github.com/omangm/dwaar/internal/tui/styles"
	"github.com/omangm/dwaar/internal/tui/views"
	"github.com/omangm/dwaar/internal/tunnel"
)

type Screen int

const (
	ScreenList Screen = iota
	ScreenForm
	ScreenLogs
)

type keyMap struct {
	New        key.Binding
	Edit       key.Binding
	Delete     key.Binding
	Toggle     key.Binding
	Logs       key.Binding
	GlobalLogs key.Binding
	Theme      key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.New, k.Edit, k.Delete},
		{k.Toggle, k.Logs, k.GlobalLogs, k.Theme},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	New:        key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Toggle:     key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
	Logs:       key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
	GlobalLogs: key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "global logs")),
	Theme:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

type Model struct {
	screen       Screen
	list         list.Model
	formModel    *views.RuleFormModel
	logsView     viewport.Model
	help         help.Model
	showHelp     bool
	isCatppuccin bool
	logRuleID    string
	editIdx      int
	rules        []tunnel.ForwardRule
	statuses     map[string]tunnel.TunnelStatus
	manager      *tunnel.Manager
	width        int
	height       int
	err          error
}

func NewModel(rules []tunnel.ForwardRule, mgr *tunnel.Manager) Model {
	m := Model{
		screen:   ScreenList,
		rules:    rules,
		statuses: make(map[string]tunnel.TunnelStatus),
		manager:  mgr,
		help:     help.New(),
	}

	l := list.New(m.listItems(), list.NewDefaultDelegate(), 0, 0)
	l.Title = "Dwaar - Port Forwarding"
	l.SetShowHelp(false)
	m.list = l

	m.logsView = viewport.New(0, 0)
	m.logsView.YPosition = 2

	return m
}

func (m Model) listItems() []list.Item {
	items := make([]list.Item, len(m.rules))
	for i, r := range m.rules {
		items[i] = ruleItem{rule: r, status: m.statuses[r.ID]}
	}
	return items
}

type statusMsg tunnel.StatusEvent

type tickMsg time.Time

func tickLogUpdate() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func listenForEvents(eventsCh <-chan tunnel.StatusEvent) tea.Cmd {
	return func() tea.Msg {
		return statusMsg(<-eventsCh)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenForEvents(m.manager.Events()),
		tickLogUpdate(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		f, _ := os.OpenFile("/tmp/dwaar_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if f != nil {
			f.WriteString(fmt.Sprintf("Key pressed: %s\n", keyMsg.String()))
			f.Close()
		}
	}

	switch msg := msg.(type) {
	case statusMsg:
		f, _ := os.OpenFile("/tmp/dwaar_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if f != nil {
			f.WriteString(fmt.Sprintf("Status updated: %s %v\n", msg.ID, msg.Status.State))
			f.Close()
		}
		m.statuses[msg.ID] = msg.Status
		m.list.SetItems(m.listItems())
		cmds = append(cmds, listenForEvents(m.manager.Events()))
		return m, tea.Batch(cmds...)

	case tickMsg:
		if m.screen == ScreenLogs {
			if m.logRuleID != "" {
				m.logsView.SetContent(m.manager.GetLogs(m.logRuleID))
			} else {
				m.logsView.SetContent(m.manager.GetGlobalLogs())
			}
		}
		if m.screen == ScreenList {
			for id, t := range m.statuses {
				if t.State == tunnel.StateRunning {
					m.statuses[id] = m.manager.Status(id)
				}
			}
			m.list.SetItems(m.listItems())
		}
		cmds = append(cmds, tickLogUpdate())

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.showHelp {
			m.list.SetSize(msg.Width, msg.Height-5)
		} else {
			m.list.SetSize(msg.Width, msg.Height)
		}
		m.logsView.Width = msg.Width
		m.logsView.Height = msg.Height - 4
		m.help.Width = msg.Width
	}

	switch m.screen {
	case ScreenList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if m.list.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, keys.Help):
				m.showHelp = !m.showHelp
				m.help.ShowAll = true
				if m.showHelp {
					m.list.SetSize(m.width, m.height-5)
				} else {
					m.list.SetSize(m.width, m.height)
				}
			case key.Matches(msg, keys.Theme):
				m.isCatppuccin = !m.isCatppuccin
				if m.isCatppuccin {
					styles.ActiveTheme = styles.CatppuccinTheme()
				} else {
					styles.ActiveTheme = styles.DefaultTheme()
				}
			case key.Matches(msg, keys.Toggle):
				if i, ok := m.list.SelectedItem().(ruleItem); ok {
					cmds = append(cmds, toggleTunnelCmd(m.manager, i.rule))
				}
			case key.Matches(msg, keys.Logs):
				if i, ok := m.list.SelectedItem().(ruleItem); ok {
					m.screen = ScreenLogs
					m.logRuleID = i.rule.ID
					m.logsView.SetContent(m.manager.GetLogs(i.rule.ID))
					m.logsView.GotoBottom()
				}
			case key.Matches(msg, keys.GlobalLogs):
				m.screen = ScreenLogs
				m.logRuleID = "" // Empty indicates global logs
				m.logsView.SetContent(m.manager.GetGlobalLogs())
				m.logsView.GotoBottom()
			case key.Matches(msg, keys.New):
				m.screen = ScreenForm
				m.editIdx = -1
				m.formModel = views.NewRuleForm(nil)
				cmds = append(cmds, m.formModel.Form.Init())
				return m, tea.Batch(cmds...)
			case key.Matches(msg, keys.Edit):
				if idx := m.list.Index(); idx >= 0 && len(m.rules) > 0 {
					m.screen = ScreenForm
					m.editIdx = idx
					m.formModel = views.NewRuleForm(&m.rules[idx])
					cmds = append(cmds, m.formModel.Form.Init())
					return m, tea.Batch(cmds...)
				}
			case key.Matches(msg, keys.Delete):
				if idx := m.list.Index(); idx >= 0 && len(m.rules) > 0 {
					rule := m.rules[idx]
					if m.statuses[rule.ID].State == tunnel.StateRunning {
						m.manager.Stop(rule.ID)
					}
					m.rules = append(m.rules[:idx], m.rules[idx+1:]...)
					m.list.SetItems(m.listItems())
					saveRules(m.rules)
				}
			}
		}

		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

	case ScreenLogs:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				m.screen = ScreenList
				m.logRuleID = ""
			}
		}
		m.logsView, cmd = m.logsView.Update(msg)
		cmds = append(cmds, cmd)

	case ScreenForm:
		form, cmd := m.formModel.Form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.formModel.Form = f
			cmds = append(cmds, cmd)
		}

		if m.formModel.Form.State == huh.StateCompleted {
			rule := m.formModel.ToRule()
			if m.editIdx == -1 {
				rule.ID = fmt.Sprintf("rule_%d", time.Now().UnixNano())
				m.rules = append(m.rules, rule)
			} else {
				rule.ID = m.rules[m.editIdx].ID
				rule.Via = m.rules[m.editIdx].Via
				m.rules[m.editIdx] = rule
			}
			m.list.SetItems(m.listItems())
			saveRules(m.rules)
			m.screen = ScreenList
		} else if m.formModel.Form.State == huh.StateAborted {
			m.screen = ScreenList
		}
	}

	return m, tea.Batch(cmds...)
}

func saveRules(rules []tunnel.ForwardRule) {
	path, err := config.DefaultPath()
	if err == nil {
		cfg, err := config.Load(path)
		if err == nil {
			cfg.Rules = rules
			config.Save(path, cfg)
		}
	}
}

func toggleTunnelCmd(mgr *tunnel.Manager, rule tunnel.ForwardRule) tea.Cmd {
	return func() tea.Msg {
		status := mgr.Status(rule.ID)
		if status.State == tunnel.StateRunning {
			mgr.Stop(rule.ID)
		} else {
			mgr.Start(rule)
		}
		return nil
	}
}

func (m Model) View() string {
	switch m.screen {
	case ScreenList:
		view := m.list.View()
		if m.showHelp {
			helpView := m.help.View(keys)
			view = lipgloss.JoinVertical(lipgloss.Left, view, "\n", helpView)
		}
		return view
	case ScreenForm:
		return fmt.Sprintf("  Adding/Editing Rule\n\n%s", m.formModel.Form.View())
	case ScreenLogs:
		var title string
		if m.logRuleID == "" {
			title = " Global Logs (Press ESC to return) "
		} else {
			title = fmt.Sprintf(" Logs for %s (Press ESC to return) ", m.logRuleID)
		}
		titleBar := styles.ActiveTheme.TitleBar.Render(title)
		return fmt.Sprintf("%s\n\n%s", titleBar, m.logsView.View())
	default:
		return "Unknown screen"
	}
}

type ruleItem struct {
	rule   tunnel.ForwardRule
	status tunnel.TunnelStatus
}

func (i ruleItem) Title() string {
	prefix := "[ ] "
	if i.status.State == tunnel.StateRunning {
		prefix = "[*] "
	} else if i.status.State == tunnel.StateError {
		prefix = "[!] "
	}
	title := prefix + i.rule.Name

	if i.status.State == tunnel.StateRunning {
		return styles.ActiveTheme.StatusRunning.Render(title)
	} else if i.status.State == tunnel.StateError {
		return styles.ActiveTheme.StatusError.Render(title)
	}
	return styles.ActiveTheme.StatusStopped.Render(title)
}

func (i ruleItem) Description() string {
	desc := fmt.Sprintf("Local: %d -> %s:%d", i.rule.LocalPort, i.rule.RemoteHost, i.rule.RemotePort)
	if i.status.State == tunnel.StateRunning {
		desc += fmt.Sprintf(" | ↑ %s ↓ %s", humanize.Bytes(i.status.BytesSent), humanize.Bytes(i.status.BytesRecv))
	}
	return styles.ActiveTheme.Muted.Render(desc)
}

func (i ruleItem) FilterValue() string { return i.rule.Name }
