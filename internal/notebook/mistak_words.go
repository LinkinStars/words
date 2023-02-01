package notebook

import (
	"encoding/json"
	"os"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/dict"
	"github.com/LinkinStars/words/internal/model"
	"github.com/LinkinStars/words/internal/storage"
)

// AddMistakeWord add mistake word
func AddMistakeWord(word *dict.Word, book string) (err error) {
	_, exists, err := GetWordFromMistakeNote(word.Name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	v := &model.MistakeWord{
		Word:    word.Name,
		Content: word,
		Book:    book,
	}
	_, err = storage.DB.InsertOne(v)
	return err
}

func GetWordFromMistakeNote(word string) (w *model.MistakeWord, exists bool, err error) {
	w = &model.MistakeWord{Word: word}
	exists, err = storage.DB.Get(w)
	return w, exists, err
}

func RemoveWordFromMistakeNote(id int) (err error) {
	w := &model.MistakeWord{}
	_, err = storage.DB.ID(id).Delete(w)
	return err
}

func ImportMistakeWords(importFilePath string) (err error) {
	content, err := os.ReadFile(importFilePath)
	if err != nil {
		return err
	}
	words := make([]*ExportMistakeWord, 0)
	_ = json.Unmarshal(content, &words)
	for _, w := range words {
		mistakeNote, exists, err := GetWordFromMistakeNote(w.Content.Name)
		if err != nil {
			return err
		}
		if !exists {
			logger.Debugf("%s word not found from mistake note", w.Content.Name)
			continue
		}
		dict.UpdateWordByBook(w.Content, w.Book)
		if err = RemoveWordFromMistakeNote(mistakeNote.ID); err != nil {
			return err
		}
	}
	return nil
}

func ExportMistakeWords(exportFilePath string) (err error) {
	words := make([]*ExportMistakeWord, 0)
	cond := make([]*model.MistakeWord, 0)
	err = storage.DB.Find(&cond)
	if err != nil {
		return
	}
	for _, word := range cond {
		words = append(words, &ExportMistakeWord{
			Content: word.Content,
			Book:    word.Book,
		})
	}
	content, _ := json.Marshal(words)
	err = os.WriteFile(exportFilePath, content, 0644)
	return
}

type ExportMistakeWord struct {
	Content *dict.Word
	Book    string
}
