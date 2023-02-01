package page

import tea "github.com/charmbracelet/bubbletea"

type Page interface {
	Update(msg tea.Msg) tea.Cmd
	View() (page string)
	HelpView() (page string)
}
