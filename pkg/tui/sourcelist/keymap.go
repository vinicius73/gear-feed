package sourcelist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinicius73/gear-feed/pkg/scraper"
)

type delegateKeyMap struct {
	choose key.Binding
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	def := list.NewDefaultDelegate()

	def.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var item scraper.SourceDefinition

		if i, ok := m.SelectedItem().(SourceItem); ok {
			item = i.SourceDefinition
		} else {
			return nil
		}

		if msg, ok := msg.(tea.KeyMsg); ok && key.Matches(msg, keys.choose) {
			return tea.Batch(
				func() tea.Msg {
					return item
				},
			)
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	def.ShortHelpFunc = func() []key.Binding {
		return help
	}

	def.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return def
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
		},
	}
}
