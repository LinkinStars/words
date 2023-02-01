package model

import (
	"time"

	"github.com/LinkinStars/words/internal/dict"
)

// Vocabulary 词汇本
type Vocabulary struct {
	ID        int        `xorm:"not null pk autoincr INT(11) id"`
	CreatedAt time.Time  `xorm:"created TIMESTAMP created_at"`
	UpdatedAt time.Time  `xorm:"updated TIMESTAMP updated_at"`
	Word      string     `xorm:"not null default '' UNIQUE(word_book) VARCHAR(200) word"`
	Content   *dict.Word `xorm:"not null default '' TEXT JSON content"`
	Book      string     `xorm:"not null default '' UNIQUE(word_book) VARCHAR(200) book"`
	Degree    int        `xorm:"not null default 1 INT(11) degree"`
}

// TableName vocabulary
func (Vocabulary) TableName() string {
	return "vocabulary"
}
