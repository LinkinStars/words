package book

import (
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/storage"
)

const readingBookKey = "reading_book"

func GetReadingBook() {
	s := &model.Config{Key: readingBookKey}
	exist, err := storage.DB.Get(s)
	if err != nil {
		logger.Error(err)
		return
	}
	if exist {
		book := dict.FindBook(s.Value)
		if book != nil {
			dict.ChangeCurrentDictionary(book)
		}
	} else {
		AddReadingBook()
	}
}

func UpdateReadingBook() {
	cond := &model.Config{Value: dict.CurrentDictionary.Name}
	_, err := storage.DB.Where("key = ?", readingBookKey).Update(cond)
	if err != nil {
		logger.Error(err)
	}
}

func AddReadingBook() {
	_, err := storage.DB.InsertOne(&model.Config{Key: readingBookKey, Value: dict.CurrentDictionary.Name})
	if err != nil {
		logger.Error(err)
	}
}
