package plan

import (
	"fmt"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/notebook"
	"github.com/dariubs/percent"
)

const (
	maxDegree = 2
)

var (
	CurrentPlan         = make(map[string]*dict.Word)
	CurrentTargetAmount = 0
	DailyNewWordAmount  = 7
	DailyReviewAmount   = 21

	TodayNewWordAmount    int
	TodayReviewWordAmount int
)

func InitCurrentPlan() {
	// read db
	todayPlan, exists, err := GetTodayPlan()
	if err != nil {
		logger.Error(err)
	}
	if !exists {
		err = AddTodayPlan()
		if err != nil {
			logger.Error(err)
		}
	} else {
		TodayNewWordAmount = todayPlan.NewAmount
		TodayReviewWordAmount = todayPlan.ReviewAmount
	}

	// 根据当前计划随机需要的单词
	newNeed, reviewNeed := DailyNewWordAmount-TodayNewWordAmount, DailyReviewAmount-TodayReviewWordAmount
	if newNeed < 0 {
		newNeed = 0
	}
	if reviewNeed < 0 {
		reviewNeed = 0
	}

	words := getWordsFromCurrentBook(newNeed)
	logger.Debugf("从单词书中获取到 %d", len(words))
	addWordsToCurrentPlan(words)

	words = getWordsFromNotebook(reviewNeed)
	logger.Debugf("从生词本中获取到 %d", len(words))
	addWordsToCurrentPlan(words)
	CurrentTargetAmount = len(CurrentPlan)
}

// ResetCurrentPlan 选择单词书后重置计划
func ResetCurrentPlan() {
	CurrentPlan = make(map[string]*dict.Word)
	dict.CurWord = nil
}

// Supplement 再来一组
func Supplement() {
	words := getWordsFromCurrentBook(DailyNewWordAmount)
	logger.Debugf("从单词书中获取到 %d", len(words))
	addWordsToCurrentPlan(words)

	words = getWordsFromNotebook(DailyReviewAmount)
	logger.Debugf("从生词本中获取到 %d", len(words))
	addWordsToCurrentPlan(words)

	dict.CurWord = RandomWord()
	CurrentTargetAmount = len(CurrentPlan)
}

func CurrentLearningProgress() string {
	target := CurrentTargetAmount
	remain := len(CurrentPlan)
	return fmt.Sprintf("%.2f%%", percent.PercentOf(target-remain, target))
}

func addWordsToCurrentPlan(words []*dict.Word) {
	for _, word := range words {
		word.Degree = 1
		CurrentPlan[word.Name] = word
	}
}

// 从当前选择的单词书中获取对应数量的单词
func getWordsFromCurrentBook(amount int) (words []*dict.Word) {
	if amount == 0 {
		return words
	}
	vocabulary, err := notebook.GetAllVocabulary(dict.CurrentDictionary.Name)
	if err != nil {
		logger.Error(err)
	}

	for _, w := range vocabulary {
		delete(dict.CurrentDictionary.DictStudyMapping, w)
	}
	for _, idx := range dict.CurrentDictionary.DictStudyMapping {
		word := dict.CurrentDictionary.Words[idx]
		word.IsNew = true
		words = append(words, word)
		if len(words) == amount {
			return words
		}
	}
	return words
}

// 从生词本中获取要复习的单词
func getWordsFromNotebook(amount int) (words []*dict.Word) {
	if amount == 0 {
		return words
	}
	vocabulary, err := notebook.RandomUnfamiliarWords(amount, dict.CurrentDictionary.Name)
	if err != nil {
		logger.Error(err)
	}
	for _, v := range vocabulary {
		// 如果当前单词书中有这个单词，那么直接使用单词书中的单词内容，以保证最新
		if idx, ok := dict.CurrentDictionary.DictMapping[v.Word]; ok {
			v.Content = dict.CurrentDictionary.Words[idx]
		}
		words = append(words, v.Content)
		if len(words) == amount {
			return words
		}
	}
	return words
}
