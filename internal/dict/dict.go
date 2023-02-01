package dict

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/internal/config"
)

var (
	Dictionaries      []*Dictionary
	CurrentDictionary *Dictionary
	CurWord           *Word
)

type Dictionary struct {
	Name             string         `json:"name"`
	Brief            string         `json:"brief"`
	Words            []*Word        `json:"words"`
	DictMapping      map[string]int `json:"-"`
	DictStudyMapping map[string]int `json:"-"`
	FilePath         string         `json:"-"`
	Progress         float64        `json:"-"`
}

func (d *Dictionary) InitIndex() {
	if d.DictMapping == nil {
		d.DictMapping = make(map[string]int)
		d.DictStudyMapping = make(map[string]int)
		for i, word := range d.Words {
			d.DictMapping[word.Name] = i
			d.DictStudyMapping[word.Name] = i
		}
	}
}

func Init() {
	initBooks()
	files, err := os.ReadDir(config.BookDir)
	if err != nil {
		panic(err)
	}
	Dictionaries = make([]*Dictionary, 0)
	CurrentDictionary = nil
	CurWord = nil

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".json" {
			continue
		}
		bookFilePath := filepath.Join(config.BookDir, file.Name())
		fileContent, err := os.ReadFile(bookFilePath)
		if err != nil {
			logger.Error(err)
			continue
		}
		var dictionary *Dictionary
		err = json.Unmarshal(fileContent, &dictionary)
		if err != nil {
			logger.Error(err)
			continue
		}
		dictionary.InitIndex()
		dictionary.FilePath = bookFilePath
		Dictionaries = append(Dictionaries, dictionary)
	}

	if len(Dictionaries) == 0 {
		log.Fatalf("dictionary is empty, please download it to %s", config.BookDir)
		return
	}
	CurrentDictionary = Dictionaries[0]
}

func initBooks() {
	files, err := os.ReadDir(config.BookDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".json" {
			return
		}
	}
	err = extractBooks()
	if err != nil {
		panic(err)
	}
}

func SearchWord(word string) *Word {
	if idx, ok := CurrentDictionary.DictMapping[word]; ok {
		return CurrentDictionary.Words[idx]
	}
	return nil
}

func UpdateWordByBook(word *Word, bookName string) {
	book := FindBook(bookName)
	if book == nil {
		return
	}
	if idx, ok := book.DictMapping[word.Name]; ok {
		book.Words[idx] = word
		content, _ := json.Marshal(book)
		err := os.WriteFile(book.FilePath, content, 0644)
		if err != nil {
			logger.Error(err)
		}
	}
}

func FindBook(bookName string) *Dictionary {
	for _, d := range Dictionaries {
		if d.Name == bookName {
			return d
		}
	}
	return nil
}

func ChangeCurrentDictionary(d *Dictionary) {
	if d.DictMapping == nil {
		d.DictMapping = make(map[string]int)
		for i, word := range d.Words {
			d.DictMapping[word.Name] = i
		}
	}
	CurrentDictionary = d
}
