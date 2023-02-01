package setting

import (
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/storage"
)

const settingsKey = "settings"

func ReadSetting() {
	s := &model.Config{Key: settingsKey}
	exist, err := storage.DB.Get(s)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if exist {
		s.GetJsonValue(&config.GlobalSettings)
	}
}

func UpdateSetting() {
	cond := &model.Config{}
	cond.SetJsonValue(config.GlobalSettings)
	_, err := storage.DB.Where("key = ?", settingsKey).Update(cond)
	if err != nil {
		logger.Error(err.Error())
	}
}

func AddSettings() {
	bean := &model.Config{Key: settingsKey}
	bean.SetJsonValue(config.GlobalSettings)
	_, err := storage.DB.InsertOne(bean)
	if err != nil {
		logger.Error(err.Error())
	}
}
