package search

import (
	"fmt"
	"strings"

	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/voice"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SearchPage this page use to search for word
type SearchPage struct {
	textarea textarea.Model
	viewport viewport.Model
	word     *dict.Word
}

func NewSearchPage() *SearchPage {
	ta := textarea.New()
	ta.Placeholder = "Enter a word..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 100

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(config.MaxWidth, config.ContentMaxHeight)
	vp.KeyMap = viewport.KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
		),
	}
	settingPage := &SearchPage{
		textarea: ta,
		viewport: vp,
	}
	return settingPage
}

func (s *SearchPage) View() (page string) {
	page = s.textarea.View()
	if s.word == nil {
		return page
	}

	page += "\n"
	page += s.word.Name
	if config.GlobalSettings.IsBritish() {
		page += fmt.Sprintf(" [%s]", s.word.UKPhone)
	}
	if config.GlobalSettings.IsAmerican() {
		page += fmt.Sprintf(" [%s]", s.word.USPhone)
	}
	page += "\n"

	content := ""
	for _, tr := range s.word.Trans {
		content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranCn)
		content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranOther)
	}
	content += "\n"
	for _, tr := range s.word.Sentences {
		content += fmt.Sprintf("%s\n%s\n", tr.Content, tr.Tran)
	}
	s.viewport.SetContent(content)
	page += s.viewport.View()

	return page
}

func (s *SearchPage) HelpView() (page string) {
	return component.HelpStyle("\n输入单词后回车查询")
}

func (s *SearchPage) Update(msg tea.Msg) (cmd tea.Cmd) {
	s.viewport, cmd = s.viewport.Update(msg)

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.Type {
	case tea.KeyEnter:
		word := strings.TrimSpace(s.textarea.Value())
		s.textarea.Reset()
		s.word = dict.SearchWord(word)
	case tea.KeySpace:
		if s.word != nil {
			voice.PlayWord(s.word.Name, config.GlobalSettings.GetPronunciation())
		}
	default:
		s.textarea, cmd = s.textarea.Update(msg)
	}
	return
}
