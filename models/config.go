package models

import (
	"context"

	"github.com/real-web-world/bdk/fastcurd"
	"gorm.io/gorm"
)

const (
	ClientConfDbID        = iota + 1 // 客户端默认配置
	CurrVersionDbID                  // 当前版本信息
	DownloadUrlPrefixDbID            // 下载url前缀
)

type (
	Config struct {
		fastcurd.Base
		Key string `json:"key" gorm:"column:k"`
		Val string `json:"Val" gorm:"column:v"`
	}
)

func NewCtxConfig(ctx context.Context, db *gorm.DB) *Config {
	return &Config{
		Base: fastcurd.Base{
			Ctx: ctx,
			DB:  db,
		},
	}
}
func (m *Config) TableName() string {
	return "config"
}

func (m *Config) GetClientConf() (*Config, error) {
	return fastcurd.GetDetailByID(m, ClientConfDbID)
}
func (m *Config) GetCurrVersion() (*Config, error) {
	return fastcurd.GetDetailByID(m, CurrVersionDbID)
}
func (m *Config) GetDownloadUrlPrefix() (*Config, error) {
	return fastcurd.GetDetailByID(m, DownloadUrlPrefixDbID)
}
