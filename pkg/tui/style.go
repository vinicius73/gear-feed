package tui

import "github.com/charmbracelet/lipgloss"

var windowWidth = 0
var windowHeight = 0

type WindowSize struct {
	Width  int
	Height int
}

var (
	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
	AppStyle   = lipgloss.NewStyle().Margin(1, 1)
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF3B30")).
			Padding(0, 1)
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ModalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(2, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFFDF5")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func SetWindowSize(width int, height int) WindowSize {
	h, v := AppStyle.GetFrameSize()
	windowWidth = width - v
	windowHeight = height - h

	return WindowSize{
		Width:  windowWidth,
		Height: windowHeight,
	}
}

func GetWindowSize() WindowSize {
	return WindowSize{
		Width:  windowWidth,
		Height: windowHeight,
	}
}
