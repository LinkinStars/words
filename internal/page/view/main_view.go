package view

import (
	"strconv"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/LinkinStars/words/internal/page"
	"github.com/LinkinStars/words/internal/page/book"
	"github.com/LinkinStars/words/internal/page/dashboard"
	"github.com/LinkinStars/words/internal/page/note"
	"github.com/LinkinStars/words/internal/page/reading"
	"github.com/LinkinStars/words/internal/page/search"
	"github.com/LinkinStars/words/internal/page/setting"
	"github.com/LinkinStars/words/internal/page/writing"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dariubs/percent"
)

type PageMode int

const (
	NormalPageMode PageMode = iota + 1
	ConfigPageMode
	NotePageMode
	BookPageMode
	SearchPageMode
	MaxPageMode
)

type ViewModel struct {
	Dashboard   *dashboard.Page
	SettingPage *setting.SettingPage
	ReadingPage *reading.ReadingPage
	WritingPage *writing.WritingPage
	SearchPage  *search.SearchPage
	NotePage    *note.NotePage
	BookPage    *book.BookPage
	mode        PageMode
	showHelp    bool
}

func NewViewModel() *ViewModel {
	return &ViewModel{
		Dashboard:   dashboard.NewPage(),
		SettingPage: setting.NewSettingPage(),
		ReadingPage: reading.NewReadingPage(),
		WritingPage: writing.NewWritingPage(),
		SearchPage:  search.NewSearchPage(),
		NotePage:    note.NewNotePage(),
		BookPage:    book.NewBookPage(),
		showHelp:    true,
		mode:        1,
	}
}

func (m *ViewModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m *ViewModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if m.switchMode(key) {
			return m, nil
		}
		switch key {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "`":
			m.showHelp = !m.showHelp
		case "?":
			m.TryToAddMistakeWord()
		default:
			m.UpdateByMode(message)
		}
	}
	return m, nil
}

func (m *ViewModel) switchMode(key string) bool {
	mode, _ := strconv.Atoi(key)
	if 0 < mode && mode < int(MaxPageMode) {
		m.mode = PageMode(mode)
		return true
	}
	return false
}

func (m *ViewModel) UpdateByMode(msg tea.Msg) (t tea.Model, cmd tea.Cmd) {
	p := m.GetCurrentPage()
	cmd = p.Update(msg)
	return m, cmd
}

func (m *ViewModel) TryToAddMistakeWord() {
	if m.mode != NormalPageMode {
		return
	}
	err := notebook.AddMistakeWord(dict.CurWord, dict.CurrentDictionary.Name)
	if err != nil {
		logger.Error(err)
	}
	dict.CurWord.Mistake = true
}

func (m *ViewModel) View() (page string) {
	p := m.GetCurrentPage()
	page = p.View()
	if m.showHelp {
		return page + p.HelpView()
	}
	return page
}

func (m *ViewModel) GetCurrentPage() (page page.Page) {
	switch m.mode {
	case NormalPageMode:
		if dict.CurWord == nil {
			page = m.Dashboard
		} else if config.GlobalSettings.IsReadMode() {
			page = m.ReadingPage
		} else if config.GlobalSettings.IsSpellMode() {
			page = m.WritingPage
		}
	case ConfigPageMode:
		page = m.SettingPage
	case SearchPageMode:
		page = m.SearchPage
	case NotePageMode:
		page = m.NotePage
	case BookPageMode:
		for _, dictionary := range dict.Dictionaries {
			c, err := notebook.CountKnownVocabulary(dictionary.Name)
			if err != nil {
				logger.Error(err)
			}
			dictionary.Progress = percent.PercentOf(int(c), len(dictionary.Words))
		}
		page = m.BookPage
	default:
		page = m.ReadingPage
	}
	return page
}
