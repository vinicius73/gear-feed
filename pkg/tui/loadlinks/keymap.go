package loadlinks

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type delegateKeyMap struct {
	open key.Binding
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	def := list.NewDefaultDelegate()

	def.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// var title string
		var item Link

		if i, ok := m.SelectedItem().(Link); ok {
			item = i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.open):
				return tea.Batch(
					func() tea.Msg {
						return item
					},
				)
			}
		}

		return nil
	}

	help := []key.Binding{keys.open}

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
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open link"),
		),
	}
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.open,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.open,
		},
	}
}
