package dashboard

import (
	"fmt"

	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/page"
	"github.com/LinkinStars/words/internal/plan"
	"github.com/LinkinStars/words/internal/voice"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	_         page.Page = (*Page)(nil)
	descStyle           = lipgloss.NewStyle().MarginTop(1)
	subtle              = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	infoStyle           = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#1E222A")).
				Background(lipgloss.Color("#80A1C1")).
				MarginRight(2).
				Underline(true)
)

// Page 首页
type Page struct {
}

func NewPage() *Page {
	p := &Page{}
	return p
}

func (s *Page) View() (page string) {
	todayAmount := plan.TodayNewWordAmount + plan.TodayReviewWordAmount
	todayAmountStyle :=
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Render(fmt.Sprintf("%d", todayAmount))
	desc := lipgloss.JoinVertical(lipgloss.Center,
		descStyle.Render("黑发不知勤学早，头秃方悔读书迟~"),
		infoStyle.Render(
			fmt.Sprintf("今日在 words 已学 %s 个单词 加油💪", todayAmountStyle)),
	)

	okButton := activeButtonStyle.Render("按空格键 开始学习")
	if todayAmount > 0 {
		okButton = activeButtonStyle.Render("按空格键 再来一组")
	}
	row := lipgloss.JoinVertical(lipgloss.Center, plan.TodayDate, desc, okButton)
	page += row + "\n"
	return page
}

func (s *Page) HelpView() (page string) {
	return ""
}

func (s *Page) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.Type {
	case tea.KeySpace:
		plan.Supplement()
		if dict.CurWord != nil && config.GlobalSettings.IsAutoVoice() {
			voice.PlayWord(dict.CurWord.Name, config.GlobalSettings.GetPronunciation())
		}
	}
	return nil
}
