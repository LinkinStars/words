package model

import (
	"time"

	"github.com/LinkinStars/words/internal/dict"
)

// MistakeWord 错误单词
type MistakeWord struct {
	ID        int        `xorm:"not null pk autoincr INT(11) id"`
	CreatedAt time.Time  `xorm:"created TIMESTAMP created_at"`
	UpdatedAt time.Time  `xorm:"updated TIMESTAMP updated_at"`
	Word      string     `xorm:"not null default '' UNIQUE VARCHAR(200) word"`
	Content   *dict.Word `xorm:"not null default '' TEXT JSON content"`
	Book      string     `xorm:"not null default '' VARCHAR(200) book"`
}

// TableName mistake_word
func (MistakeWord) TableName() string {
	return "mistake_word"
}
