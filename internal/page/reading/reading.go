package reading

import (
	"fmt"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/LinkinStars/words/internal/page"
	"github.com/LinkinStars/words/internal/plan"
	"github.com/LinkinStars/words/internal/voice"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var _ page.Page = (*ReadingPage)(nil)

// ReadingPage this page is user read words
type ReadingPage struct {
	viewport                   viewport.Model
	showDetail                 bool
	showAddUnfamiliarWordsBook bool
}

func NewReadingPage() *ReadingPage {
	vp := viewport.New(config.MaxWidth, config.ContentMaxHeight)
	readingPage := &ReadingPage{
		viewport: vp,
	}
	return readingPage
}

func (r *ReadingPage) View() (page string) {
	page = dict.CurWord.Name
	if dict.CurWord.Mistake {
		page += component.HelpStyle("[已加入错误本]")
	}
	if config.GlobalSettings.IsBritish() {
		page += fmt.Sprintf(" [%s]", dict.CurWord.UKPhone)
	}
	if config.GlobalSettings.IsAmerican() {
		page += fmt.Sprintf(" [%s]", dict.CurWord.USPhone)
	}
	if r.showAddUnfamiliarWordsBook {
		page += component.HelpStyle(" [已加入生词本]")
	}
	page += "\n"

	content := ""
	if r.showDetail {
		for _, tr := range dict.CurWord.Trans {
			content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranCn)
			content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranOther)
		}
		content += "\n"
		for _, tr := range dict.CurWord.Sentences {
			content += fmt.Sprintf("%s\n%s\n", tr.Content, tr.Tran)
		}
	}
	r.viewport.SetContent(content)
	page += r.viewport.View()
	page += "\n"
	return page
}

func (r *ReadingPage) HelpView() (page string) {
	return component.HelpStyle("←: 不认识 • →: 认识 • 空格: 发音 " +
		fmt.Sprintf("[%s] [%s]", dict.CurrentDictionary.Brief, plan.CurrentLearningProgress()))
}

func (r *ReadingPage) Update(msg tea.Msg) (cmd tea.Cmd) {
	r.viewport, cmd = r.viewport.Update(msg)
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch key.Type {
	case tea.KeyLeft:
		r.showDetail = true
		r.showAddUnfamiliarWordsBook = true
		if err := notebook.AddWord2Vocabulary(dict.CurWord, dict.CurrentDictionary.Name, 3); err != nil {
			logger.Debugf("add word to vocabulary failed: %v", err)
		}
		plan.ForgetWord(dict.CurWord.Name)
	case tea.KeyRight:
		if r.showDetail {
			// 如果已经展示详情，并且没有加入生词本
			if !r.showAddUnfamiliarWordsBook {
				plan.RememberWord(dict.CurWord.Name)
				if err := notebook.AddWord2Vocabulary(dict.CurWord, dict.CurrentDictionary.Name, 0); err != nil {
					logger.Debugf("add word to vocabulary failed: %v", err)
				}
			}
			r.showDetail = false
			dict.CurWord = plan.RandomWord()
			if dict.CurWord != nil && config.GlobalSettings.IsAutoVoice() {
				voice.PlayWord(dict.CurWord.Name, config.GlobalSettings.GetPronunciation())
			}
		} else {
			r.showDetail = true
		}
		r.showAddUnfamiliarWordsBook = false
	case tea.KeySpace:
		voice.PlayWord(dict.CurWord.Name, config.GlobalSettings.GetPronunciation())
	}
	return nil
}
