package storage

import (
	"os"
	"path/filepath"

	"github.com/LinkinStars/golang-util/gu"
	"github.com/LinkinStars/words/internal/config"
)

const DataDirEnv = "WORDS_DATA"

func init() {
	dataDir()
	config.BookDir = filepath.Join(config.DataDir, "book")
	gu.CreateDirIfNotExist(config.BookDir)

	config.VoiceDir = filepath.Join(config.DataDir, "voice")
	gu.CreateDirIfNotExist(config.VoiceDir)

	config.DBDir = filepath.Join(config.DataDir, "db")
	gu.CreateDirIfNotExist(config.DBDir)

	config.LogDir = filepath.Join(config.DataDir, "log")
	gu.CreateDirIfNotExist(config.LogDir)
}

func dataDir() {
	config.DataDir = os.Getenv(DataDirEnv)
	if len(config.DataDir) == 0 {
		config.DataDir, _ = os.UserHomeDir()
	}
	if len(config.DataDir) == 0 {
		config.DataDir = os.TempDir()
	}
	config.DataDir = filepath.Join(config.DataDir, "words")
	if err := gu.CreateDirIfNotExist(config.DataDir); err != nil {
		panic(err)
	}
}
