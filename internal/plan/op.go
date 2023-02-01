package plan

import (
	"time"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/storage"
)

var (
	TodayDate = time.Now().Format("2006-01-02")
)

func GetTodayPlan() (resp *model.Plan, exists bool, err error) {
	resp = &model.Plan{TodayDate: TodayDate}
	exists, err = storage.DB.Get(resp)
	if err != nil {
		logger.Error(err)
	}
	return
}

func AddTodayPlan() (err error) {
	cond := &model.Plan{TodayDate: TodayDate}
	_, err = storage.DB.InsertOne(cond)
	if err != nil {
		logger.Error(err)
	}
	return
}

func UpdateTodayPlan(newAmount, reviewAmount int) (err error) {
	cond := &model.Plan{NewAmount: newAmount, ReviewAmount: reviewAmount}
	_, err = storage.DB.Where("today_date = ?", TodayDate).Update(cond)
	if err != nil {
		logger.Error(err)
	}
	return
}
