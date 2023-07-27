package sources

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/tui"
	"github.com/vinicius73/gamer-feed/pkg/tui/loadlinks"
	"github.com/vinicius73/gamer-feed/pkg/tui/sourcelist"
)

type mode int

const (
	listMode mode = iota
	detailMode
)

type Model struct {
	list     tea.Model
	links    tea.Model
	mode     mode
	err      error
	ctx      context.Context
	quitting bool
}

func NewModel(ctx context.Context, sources []scraper.SourceDefinition) Model {
	return Model{
		ctx:  ctx,
		mode: listMode,
		list: sourcelist.New(sources),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// handle window size
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		tui.SetWindowSize(msg.Width, msg.Height)
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if m.err != nil && (k == "esc") {
			m.err = nil
			return m, nil
		}

		// if k == "q" || k == "esc" || k == "ctrl+c" {
		// 	m.quitting = true
		// 	return m, tea.Quit
		// }
	}
	switch msg := msg.(type) {
	case tui.ErrorMsg:
		m.err = msg
		return m, cmd
	case scraper.SourceDefinition:
		m.mode = detailMode
		m.links = loadlinks.New(m.ctx, msg)
		return m, m.links.Init()
	}

	switch m.mode {
	case listMode:
		m.list, cmd = m.list.Update(msg)
	case detailMode:
		m.links, cmd = m.links.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var view string

	if m.err != nil {
		return tui.ErrorStyle.Render(m.err.Error())
	}

	if m.mode == listMode {
		view = m.list.View()
	} else if m.mode == detailMode {
		view = m.links.View()
	}

	return tui.AppStyle.Render(view)
}
