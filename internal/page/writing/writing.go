package writing

import (
	"fmt"
	"strings"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/LinkinStars/words/internal/page"
	"github.com/LinkinStars/words/internal/plan"
	"github.com/LinkinStars/words/internal/voice"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var _ page.Page = (*WritingPage)(nil)

const (
	maxErrorAmount = 2

	todo answerStatus = iota
	correct
	incorrect
)

type answerStatus int

// WritingPage this page is user read words
type WritingPage struct {
	input    textinput.Model
	viewport viewport.Model

	showTip                    bool
	showAnswer                 bool
	showAddUnfamiliarWordsBook bool

	// 答题错误次数，最大 3 次
	incorrectCount int

	// 当前输入的单词还是正确的字母
	currentAnswerIsCorrect bool

	answerStatus answerStatus
}

func NewWritingPage() *WritingPage {
	input := textinput.New()
	input.Placeholder = "Enter a word..."
	input.CharLimit = 100
	input.Prompt = ""
	input.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Width = 30
	input.CursorEnd()
	input.Focus()

	vp := viewport.New(config.MaxWidth, config.ContentMaxHeight)
	vp.KeyMap = viewport.KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
		),
	}
	writingPage := &WritingPage{
		input:          input,
		viewport:       vp,
		showAnswer:     false,
		incorrectCount: 0,
		answerStatus:   todo,
	}
	return writingPage
}

func (r *WritingPage) View() (page string) {
	if !r.currentAnswerIsCorrect {
		r.input.TextStyle = lipgloss.NewStyle().Foreground(component.Red)
	} else {
		r.input.TextStyle = lipgloss.NewStyle().Foreground(component.Green)
	}
	page = r.input.View()

	if r.answerStatus == incorrect {
		page += component.WrongMark
	} else if r.answerStatus == correct {
		page += component.CheckMark
	}
	page += "\n"

	if r.showAnswer {
		page += dict.CurWord.Name
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
	}

	content := ""
	for _, tr := range dict.CurWord.Trans {
		content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranCn)
		content += fmt.Sprintf("%s %s\n", tr.Pos, tr.TranOther)
	}

	if r.showAnswer {
		for _, tr := range dict.CurWord.Sentences {
			tr.Content = wordwrap.String(tr.Content, config.MaxWidth)
			sentence := component.UnderlineWordInSentence(tr.Content, dict.CurWord.Name)
			content += fmt.Sprintf("%s\n%s\n", sentence, tr.Tran)
		}
	} else {
		if r.showTip && len(dict.CurWord.Sentences) > 0 {
			tip := dict.CurWord.Sentences[0]
			tip.Content = wordwrap.String(tip.Content, config.MaxWidth)
			sentence := component.ReplaceWordWithUnderlineInSentence(tip.Content, dict.CurWord.Name)
			content += fmt.Sprintf("%s\n%s\n", sentence, tip.Tran)
		}
	}

	r.viewport.SetContent(content)
	page += r.viewport.View()
	page += "\n"
	return page
}

func (r *WritingPage) HelpView() (page string) {
	return component.HelpStyle("回车: 确认 • 空格: 发音 " +
		fmt.Sprintf("[%s] [%s]", dict.CurrentDictionary.Brief, plan.CurrentLearningProgress()))
}

func (r *WritingPage) Update(msg tea.Msg) (cmd tea.Cmd) {
	r.viewport, cmd = r.viewport.Update(msg)

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch keyMsg.Type {
	case tea.KeyEnter:
		if r.shouldChangeWord() {
			break
		}
		r.checkAnswer()
	case tea.KeySpace:
		voice.PlayWord(dict.CurWord.Name, config.GlobalSettings.GetPronunciation())
	default:
		r.input, cmd = r.input.Update(msg)
		r.checkPrefix()
	}
	return cmd
}

func (r *WritingPage) checkAnswer() {
	if r.input.Value() == dict.CurWord.Name {
		r.showAnswer = true
		r.answerStatus = correct
		if !r.showAddUnfamiliarWordsBook {
			if err := notebook.AddWord2Vocabulary(dict.CurWord, dict.CurrentDictionary.Name, 0); err != nil {
				logger.Errorf("add word to vocabulary failed: %v", err)
			}
		}
		return
	}
	r.answerStatus = incorrect
	r.incorrectCount++
	r.showTip = true
	r.showAddUnfamiliarWordsBook = true
	if r.incorrectCount >= maxErrorAmount {
		r.showAnswer = true
		if err := notebook.AddWord2Vocabulary(dict.CurWord, dict.CurrentDictionary.Name, 3); err != nil {
			logger.Errorf("add word to vocabulary failed: %v", err)
		}
	}
}

func (r *WritingPage) checkPrefix() {
	if dict.CurWord.Name == r.input.Value() {
		r.checkAnswer()
	}
	if strings.HasPrefix(dict.CurWord.Name, r.input.Value()) {
		r.currentAnswerIsCorrect = true
	} else {
		r.currentAnswerIsCorrect = false
	}
}

func (r *WritingPage) shouldChangeWord() bool {
	if !r.showAnswer || r.input.Value() != dict.CurWord.Name {
		return false
	}

	// 如果需要显示加入生词本，则当前单词没记住
	if r.showAddUnfamiliarWordsBook {
		plan.ForgetWord(dict.CurWord.Name)
	} else {
		plan.RememberWord(dict.CurWord.Name)
	}

	r.incorrectCount = 0
	r.currentAnswerIsCorrect = true
	r.showAddUnfamiliarWordsBook = false
	r.answerStatus = todo
	r.showTip = false
	r.showAnswer = false

	r.input.Reset()

	dict.CurWord = plan.RandomWord()
	if dict.CurWord == nil {
		return true
	}
	if config.GlobalSettings.IsAutoVoice() {
		voice.PlayWord(dict.CurWord.Name, config.GlobalSettings.GetPronunciation())
	}
	r.input.CharLimit = len(dict.CurWord.Name)
	return true
}
