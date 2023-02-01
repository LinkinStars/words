package plan

import (
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/dict"
)

func RandomWord() *dict.Word {
	if len(CurrentPlan) == 0 {
		return nil
	}
	for _, word := range CurrentPlan {
		// 如果当前随机的单词是当前正在学习的单词，并且还有别的单词可选，则尽可能不要连续同一个单词
		if word == dict.CurWord && len(CurrentPlan) > 1 {
			continue
		}
		return word
	}
	return nil
}

func RememberWord(word string) {
	if CurrentPlan[word] == nil {
		return
	}
	CurrentPlan[word].Degree--
	if CurrentPlan[word].Degree <= 0 {
		if CurrentPlan[word].IsNew {
			TodayNewWordAmount++
		} else {
			TodayReviewWordAmount++
		}
		delete(CurrentPlan, word)
		err := UpdateTodayPlan(TodayNewWordAmount, TodayReviewWordAmount)
		if err != nil {
			logger.Error(err)
		}
	}
}

func ForgetWord(word string) {
	if CurrentPlan[word] == nil {
		return
	}
	CurrentPlan[word].Degree++
	if CurrentPlan[word].Degree > maxDegree {
		CurrentPlan[word].Degree = maxDegree
	}
}
