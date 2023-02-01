package book

import (
	"fmt"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/LinkinStars/words/internal/page"
	"github.com/LinkinStars/words/internal/plan"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ page.Page = (*BookPage)(nil)
)

// BookPage this page contains all settings
type BookPage struct {
	// current user pointer position
	CurrentPointer int
}

func NewBookPage() *BookPage {
	settingPage := &BookPage{}
	GetReadingBook()
	for i, dictionary := range dict.Dictionaries {
		if dict.CurrentDictionary.Name == dictionary.Name {
			settingPage.CurrentPointer = i
		}
	}
	return settingPage
}

func (s *BookPage) View() (page string) {
	for _, dictionary := range dict.Dictionaries {
		choose := dict.CurrentDictionary.Name == dictionary.Name
		if choose {
			page += "➜ "
		} else {
			page += "  "
		}
		page += component.Checkbox(dictionary.Brief, choose)
		page += fmt.Sprintf(" %.2f%%", dictionary.Progress)
		page += "\n"
	}
	return page
}

func (s *BookPage) HelpView() (page string) {
	return component.HelpStyle("←/↑/↓/→: 选择后自动保存 ctrl+p: 清空单词书进度")
}

func (s *BookPage) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.Type {
	case tea.KeyUp:
		s.CurrentPointer = sub(s.CurrentPointer)
	case tea.KeyDown:
		s.CurrentPointer = add(s.CurrentPointer, len(dict.Dictionaries)-1)
	case tea.KeyCtrlP:
		book := dict.Dictionaries[s.CurrentPointer]
		err := notebook.ClearVocabulary(book.Name)
		if err != nil {
			logger.Error(err)
		}
		book.Progress = 0
	}
	dict.ChangeCurrentDictionary(dict.Dictionaries[s.CurrentPointer])
	UpdateReadingBook()
	plan.ResetCurrentPlan()
	return nil
}

func sub(cur int) int {
	if cur <= 0 {
		return 0
	}
	return cur - 1
}

func add(cur, max int) int {
	if cur == max {
		return max
	}
	return cur + 1
}
