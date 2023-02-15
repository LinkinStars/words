package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LinkinStars/go-scaffold/contrib/log/zap"
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/page/view"
	"github.com/LinkinStars/words/internal/plan"
	"github.com/LinkinStars/words/internal/processer"
	"github.com/LinkinStars/words/internal/storage"
	"github.com/blang/semver"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var (
	Version        = "v0.0.0"
	importFilePath string
	exportFilePath string
	upgrade        bool
)

func init() {
	flag.StringVar(&importFilePath, "i", "", "import words to json file")
	flag.StringVar(&exportFilePath, "e", "", "export words to json file")
	flag.BoolVar(&upgrade, "u", false, "upgrade words")
}

func main() {
	showVersion := flag.Bool("v", false, "print version")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}

	if upgrade {
		doSelfUpdate()
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

func doSelfUpdate() {
	fmt.Println("当前版本为：", Version, "，正在检查更新...若网络访问不稳定，请耐心等待。")
	v := semver.MustParse(strings.TrimPrefix(Version, "v"))
	latest, err := selfupdate.UpdateSelf(v, "LinkinStars/words")
	if err != nil {
		fmt.Println("更新失败：", err)
		return
	}
	if latest.Version.Equals(v) {
		fmt.Println("当前版本已是最新版本", Version)
	} else {
		fmt.Println("更新成功，当前版本为：", latest.Version)
		fmt.Println("更新内容：\n", latest.ReleaseNotes)
	}
}
