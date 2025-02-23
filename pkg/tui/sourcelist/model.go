package sourcelist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinicius73/gear-feed/pkg/scraper"
	"github.com/vinicius73/gear-feed/pkg/tui"
)

type Model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

func New(sources []scraper.SourceDefinition) Model {
	itens := make([]list.Item, len(sources))

	for i, source := range sources {
		itens[i] = SourceItem{source}
	}

	delegateKeys := newDelegateKeyMap()

	delegate := newItemDelegate(delegateKeys)

	list := list.New(itens, delegate, 0, 0)
	list.Title = "Sources"
	list.Styles.Title = tui.TitleStyle

	return Model{
		list:         list,
		delegateKeys: delegateKeys,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := tui.AppStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
