package loadlinks

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinicius73/gamer-feed/pkg/browser"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/tui"
)

const loadTimeout = time.Second * 10

type state int

const (
	loading state = iota
	ready
)

type readyMsg struct {
	entries []list.Item
}

//nolint:containedctx
type Model struct {
	ctx     context.Context
	entry   scraper.SourceDefinition
	spinner spinner.Model
	state   state
	list    list.Model `exhaustruct:"optional"`
}

func buildList(entry scraper.SourceDefinition, entries []list.Item) list.Model {
	windowSize := tui.GetWindowSize()

	delegate := newItemDelegate(newDelegateKeyMap())

	l := list.New(entries, delegate, windowSize.Width, windowSize.Height)
	l.Title = entry.Name
	l.Styles.Title = tui.TitleStyle

	return l
}

func New(ctx context.Context, entry scraper.SourceDefinition) Model {
	spn := spinner.New()
	spn.Spinner = spinner.Dot
	spn.Style = tui.SpinnerStyle

	return Model{
		ctx:     ctx,
		spinner: spn,
		state:   loading,
		entry:   entry,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadLinks(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	// handle window size message
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h, v := tui.AppStyle.GetFrameSize()

		m.list.SetSize(msg.Width-h, msg.Height-v)

		return m, nil
	}

	// handle ready message
	if readyMsg, ok := msg.(readyMsg); ok {
		m.state = ready
		m.list = buildList(m.entry, readyMsg.entries)

		cmds = append(
			cmds,
			tea.ClearScreen,
			tea.EnterAltScreen,
		)

		return m, tea.Sequence(cmds...)
	}

	// handle esc key
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyEscape {
			return m, tui.Back
		}
	}

	// handle link selection
	if link, ok := msg.(Link); ok {
		err := browser.OpenURL(link.Entry.URL)
		if err != nil {
			return m, func() tea.Msg {
				return tui.Error(err)
			}
		}

		return m, nil
	}

	switch msg := msg.(type) {
	case readyMsg:
	case tea.WindowSizeMsg:
		h, v := tui.AppStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case spinner.TickMsg:
		newSpinner, cmd := m.spinner.Update(msg)
		m.spinner = newSpinner
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	if m.state != ready {
		return m, tea.Batch(cmds...)
	}

	newList, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.list = newList

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.state == loading {
		return tui.ModalStyle.SetString(
			fmt.Sprintf("\n\n   %s Loading %s links... press q to stop\n\n", m.spinner.View(), m.entry.BaseURL),
		).String()
	}

	l := m.list.View()

	return l
}

func (m Model) loadLinks() tea.Cmd {
	ctx, cancel := context.WithTimeout(m.ctx, loadTimeout)

	return func() tea.Msg {
		defer cancel()

		entries, err := scraper.FindEntries[model.Entry](ctx, m.entry)
		if err != nil {
			return tui.Error(err)
		}

		list := make([]list.Item, len(entries))

		for i, entry := range entries {
			list[i] = Link{Entry: entry}
		}

		return readyMsg{list}
	}
}
