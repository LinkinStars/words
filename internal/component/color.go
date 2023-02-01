package component

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	Underline = lipgloss.NewStyle().Underline(true)
)
