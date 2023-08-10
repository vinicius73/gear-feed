package tui

import tea "github.com/charmbracelet/bubbletea"

type BackMsg struct{}

func Back() tea.Msg {
	return BackMsg{}
}
