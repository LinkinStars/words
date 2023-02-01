package voice

import (
	"fmt"
	"path/filepath"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/golang-util/gu"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/downloader"
)

const (
	// type 1 为英音 2 为美音
	downloadVoiceURL = "https://dict.youdao.com/dictvoice?audio=%s&type=%d"
)

var (
	player    Player
	playQueue = make(chan string, 1)
)

func init() {
	go func() {
		for voicePath := range playQueue {
			if player == nil {
				continue
			}
			player.Play(voicePath)
		}
	}()
}

// PlayWord voiceType 1 为英音 2 为美音
func PlayWord(word string, voiceType int) {
	voicePath := downloadWordVoice(word, voiceType)
	go playVoice(voicePath)
}

func downloadWordVoice(word string, voiceType int) (voicePath string) {
	voicePath = filepath.Join(config.VoiceDir, fmt.Sprintf("%s_%d%s", word, voiceType, ".wav"))
	if gu.CheckPathIfNotExist(voicePath) {
		return voicePath
	}

	logger.Debugf("try to save file in ", voicePath)
	err := downloader.LoadAndSaveFile(fmt.Sprintf(downloadVoiceURL, word, voiceType), voicePath)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return voicePath
}

func playVoice(voicePath string) {
	if len(voicePath) == 0 {
		return
	}
	select {
	case playQueue <- voicePath:
	default:
		return
	}
	return
}
