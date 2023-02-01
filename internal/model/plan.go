package model

import (
	"time"
)

// Plan 计划
type Plan struct {
	ID           int       `xorm:"not null pk autoincr INT(11) id"`
	CreatedAt    time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt    time.Time `xorm:"updated TIMESTAMP updated_at"`
	TodayDate    string    `xorm:"not null default '' UNIQUE VARCHAR(200) today_date"`
	NewAmount    int       `xorm:"not null default 1 INT(11) new_amount"`
	ReviewAmount int       `xorm:"not null default 1 INT(11) review_amount"`
}

// TableName plan
func (Plan) TableName() string {
	return "plan"
}
