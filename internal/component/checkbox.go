package component

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	Green     = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Red       = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"}
	CheckMark = lipgloss.NewStyle().SetString("✓").Foreground(Green).String()
	WrongMark = lipgloss.NewStyle().SetString("✗").Foreground(Red).String()
)

func Checkbox(label string, checked bool) string {
	if checked {
		return fmt.Sprintf("[%s] %s", CheckMark, label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
