package config

const (
	StudyMode = iota
	Pronunciation
	AutoPlayVoice
)

var (
	GlobalSettings Settings
)

type Setting struct {
	Label          string   `json:"label"`
	Name           string   `json:"name"`
	OptionNames    []string `json:"option_names"`
	CurrentPointer int      `json:"current_pointer"`
}

type Settings []*Setting

func (sets Settings) IsReadMode() bool {
	return sets[StudyMode].CurrentPointer == 0
}

func (sets Settings) IsSpellMode() bool {
	return sets[StudyMode].CurrentPointer == 1
}

func (sets Settings) IsBritish() bool {
	return sets[Pronunciation].CurrentPointer == 0
}

func (sets Settings) IsAmerican() bool {
	return sets[Pronunciation].CurrentPointer == 1
}

// GetPronunciation 1 为英音 2 为美音
func (sets Settings) GetPronunciation() int {
	return sets[Pronunciation].CurrentPointer + 1
}

func (sets Settings) IsAutoVoice() bool {
	return sets[AutoPlayVoice].CurrentPointer == 1
}
