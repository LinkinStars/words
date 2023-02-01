package notebook

import (
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/storage"
)

// AddWord2Vocabulary add word to vocabulary
func AddWord2Vocabulary(word *dict.Word, book string, degree int) (err error) {
	old, exists, err := GetWordFromVocabulary(word.Name, book)
	if err != nil {
		return err
	}
	if exists {
		_, err = storage.DB.ID(old.ID).Incr("degree").Update(&model.Vocabulary{})
		return err
	}
	v := &model.Vocabulary{
		Word:    word.Name,
		Content: word,
		Book:    book,
		Degree:  degree,
	}
	_, err = storage.DB.InsertOne(v)
	return err
}

func GetWordFromVocabulary(word, book string) (w *model.Vocabulary, exists bool, err error) {
	w = &model.Vocabulary{Word: word, Book: book}
	exists, err = storage.DB.Get(w)
	return w, exists, err
}

func GetVocabularyPage(page, pageSize int) (words []*model.Vocabulary, count int64, err error) {
	words = make([]*model.Vocabulary, 0)
	startNum := (page - 1) * pageSize
	count, err = storage.DB.Where("degree > 0").Limit(pageSize, startNum).FindAndCount(&words)
	return
}

func GetAllVocabulary(book string) (words []string, err error) {
	words = make([]string, 0)
	session := storage.DB.Table("vocabulary")
	session.Where("book = ?", book)
	err = session.Select("word").Find(&words)
	if err != nil {
		logger.Error(err)
	}
	return words, err
}

func CountKnownVocabulary(book string) (count int64, err error) {
	return storage.DB.Where("degree = 0").Count(&model.Vocabulary{Book: book})
}

func ClearVocabulary(book string) (err error) {
	_, err = storage.DB.Delete(&model.Vocabulary{Book: book})
	return err
}

func RandomUnfamiliarWords(amount int, book string) (words []*model.Vocabulary, err error) {
	ids := make([]int, 0)
	session := storage.DB.Table("vocabulary")
	session.Where("book = ?", book)
	session.Where("degree > 0")
	err = session.Select("id").Find(&ids)
	if err != nil {
		return nil, err
	}

	words = make([]*model.Vocabulary, 0)
	if len(ids) == 0 {
		return words, nil
	}

	choseIDs := make([]int, 0)
	m := make(map[int]bool, 0)
	for _, id := range ids {
		m[id] = true
	}
	for id := range m {
		choseIDs = append(choseIDs, id)
		if len(choseIDs) == amount {
			break
		}
	}

	err = storage.DB.In("id", choseIDs).Find(&words)
	if err != nil {
		return nil, err
	}
	return words, nil
}
