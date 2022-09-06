package jpterm

import "github.com/charmbracelet/lipgloss"

var (
	text            = lipgloss.AdaptiveColor{Light: "##FAFAFA", Dark: "#0A0A0A"}
	subtle          = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight       = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special         = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	expressionStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(text).
			Background(subtle).
			MarginLeft(1)
	modeBaseStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(text).
			Background(subtle).
			MarginLeft(1)
	rModeStyle = modeBaseStyle.Copy()
	eModeStyle = modeBaseStyle.Copy()
	qModeStyle = modeBaseStyle.Copy()
	IOStyle    = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Foreground(text).
			Margin(1).
			Padding(1, 2)
	// divider = lipgloss.NewStyle().
	// 	SetString("•").
	// 	Padding(0, 1).
	// 	Foreground(subtle).
	// 	String()
)
