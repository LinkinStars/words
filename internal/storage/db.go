package storage

import (
	"path/filepath"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var (
	DB *xorm.Engine
)

func InitDB() (err error) {
	conn := filepath.Join(config.DBDir, "words.db")
	DB, err = xorm.NewEngine("sqlite", conn)
	if err != nil {
		return err
	}
	DB.SetLogLevel(xormlog.LOG_WARNING)
	if err = DB.Ping(); err != nil {
		return err
	}
	DB.SetColumnMapper(names.GonicMapper{})

	err = DB.Sync(model.Tables...)
	if err != nil {
		panic(err)
	}
	logger.Debugf("db path %s", conn)
	return nil
}
