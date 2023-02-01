package note

import (
	"fmt"
	"strings"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/LinkinStars/words/internal/page"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ page.Page = (*NotePage)(nil)

// NotePage this page is user read words
type NotePage struct {
	words []*model.Vocabulary

	page      int
	pageSize  int
	paginator paginator.Model
}

func NewNotePage() *NotePage {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	notePage := &NotePage{
		paginator: p,
	}
	return notePage
}

func (r *NotePage) View() (page string) {
	var b strings.Builder
	words, count, err := notebook.GetVocabularyPage(r.paginator.Page+1, r.paginator.PerPage)
	if err != nil {
		logger.Error(err)
		return
	}
	r.paginator.SetTotalPages(int(count))
	r.words = words
	for _, item := range r.words {
		var tr string
		if len(item.Content.Trans) > 0 {
			tr = item.Content.Trans[0].TranCn
		}
		b.WriteString(fmt.Sprintf("[%d] %s %s\n", item.Degree, item.Word, tr))
	}
	b.WriteString(r.paginator.View() + "\n")
	return b.String()
}

func (r *NotePage) HelpView() (page string) {
	return component.HelpStyle("←/→: 翻页")
}

func (r *NotePage) Update(msg tea.Msg) (cmd tea.Cmd) {
	r.paginator, cmd = r.paginator.Update(msg)
	return cmd
}
