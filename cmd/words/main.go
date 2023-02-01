package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/LinkinStars/go-scaffold/contrib/log/zap"
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/page/view"
	"github.com/LinkinStars/words/internal/plan"
	"github.com/LinkinStars/words/internal/processer"
	"github.com/LinkinStars/words/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Version        string
	importFilePath string
	exportFilePath string
)

func init() {
	flag.StringVar(&importFilePath, "i", "", "import words to json file")
	flag.StringVar(&exportFilePath, "e", "", "export words to json file")

}

func main() {
	showVersion := flag.Bool("v", false, "print version")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}

	initLogger()
	dict.Init()

	if len(importFilePath+exportFilePath) > 0 {
		processer.ImportOrExport(importFilePath, exportFilePath)
		return
	}

	plan.InitCurrentPlan()

	p := tea.NewProgram(view.NewViewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initLogger() {
	level := logger.LevelInfo
	if len(os.Getenv("DEBUG")) > 0 {
		level = logger.LevelDebug
	}
	newLogger := zap.NewLogger(level, zap.WithName("words"),
		zap.WithoutStd(), zap.WithCallerFullPath(), zap.WithPath(config.LogDir))
	logger.SetLogger(newLogger)
	if err := storage.InitDB(); err != nil {
		panic(err)
	}
}
