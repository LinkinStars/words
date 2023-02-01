package model

import (
	"encoding/json"
	"time"
)

// Config 配置
type Config struct {
	ID        int       `xorm:"not null pk autoincr INT(11) id"`
	CreatedAt time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt time.Time `xorm:"updated TIMESTAMP updated_at"`
	Key       string    `xorm:"not null default '' UNIQUE VARCHAR(200) key"`
	Value     string    `xorm:"not null default '' TEXT content"`
}

// TableName config
func (c *Config) TableName() string {
	return "config"
}

func (c *Config) SetJsonValue(v interface{}) {
	t, _ := json.Marshal(v)
	c.Value = string(t)
}

func (c *Config) GetJsonValue(v interface{}) {
	_ = json.Unmarshal([]byte(c.Value), v)
}
