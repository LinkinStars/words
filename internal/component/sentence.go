package component

import (
	"strings"
)

// UnderlineWordInSentence underline the words in the example sentence
func UnderlineWordInSentence(sentence string, word string) string {
	// TODO 当前还无法替换例句中不同时态的当前单词
	return strings.ReplaceAll(sentence, word, Underline.Render(word))
}

// ReplaceWordWithUnderlineInSentence replace the words in the example sentence using underline
func ReplaceWordWithUnderlineInSentence(sentence string, word string) string {
	return strings.ReplaceAll(sentence, word, Underline.Render("?"))
}
