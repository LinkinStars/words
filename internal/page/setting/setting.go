package setting

import (
	"fmt"

	"github.com/LinkinStars/words/internal/component"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/page"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ page.Page = (*SettingPage)(nil)
)

// SettingPage this page contains all settings
type SettingPage struct {
	// current user pointer position
	CurrentPointer int
}

func NewSettingPage() *SettingPage {
	settingPage := &SettingPage{}
	initSetting()
	return settingPage
}

func initSetting() {
	ReadSetting()
	if len(config.GlobalSettings) > 0 {
		return
	}
	config.GlobalSettings = make([]*config.Setting, 3)
	config.GlobalSettings[config.StudyMode] = &config.Setting{
		Label:       "学习模式",
		Name:        "StudyMode",
		OptionNames: []string{"认识", "拼写"},
	}
	config.GlobalSettings[config.Pronunciation] = &config.Setting{
		Label:       "单词发音",
		Name:        "Pronunciation",
		OptionNames: []string{"英音", "美音"},
	}
	config.GlobalSettings[config.AutoPlayVoice] = &config.Setting{
		Label:       "自动发音",
		Name:        "AutoPlayVoice",
		OptionNames: []string{"否", "是"},
	}
	AddSettings()
}

func (s *SettingPage) View() (page string) {
	for i, curSetting := range config.GlobalSettings {
		if s.CurrentPointer == i {
			page += "➜ "
		} else {
			page += "  "
		}
		page += fmt.Sprintf("%s\t", curSetting.Label)
		for j, option := range curSetting.OptionNames {
			page += component.Checkbox(option, j == curSetting.CurrentPointer)
			page += "  "
		}
		page += "\n"
	}
	return page
}

func (s *SettingPage) HelpView() (page string) {
	return component.HelpStyle("←/↑/↓/→: 选择后自动保存")
}

func (s *SettingPage) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.Type {
	case tea.KeyUp:
		s.CurrentPointer = sub(s.CurrentPointer)
	case tea.KeyDown:
		s.CurrentPointer = add(s.CurrentPointer, len(config.GlobalSettings)-1)
	case tea.KeyLeft:
		set := config.GlobalSettings[s.CurrentPointer]
		set.CurrentPointer = sub(set.CurrentPointer)
	case tea.KeyRight:
		set := config.GlobalSettings[s.CurrentPointer]
		set.CurrentPointer = add(set.CurrentPointer, len(set.OptionNames)-1)
	}
	UpdateSetting()
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
