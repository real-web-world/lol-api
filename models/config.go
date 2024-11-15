package models

import (
	"context"

	"github.com/real-web-world/bdk/fastcurd"
	"github.com/real-web-world/bdk/json"
	"gorm.io/gorm"
)

const (
	ClientConfDbID        = iota + 1 // 客户端默认配置
	CurrVersionDbID                  // 当前版本信息
	DownloadUrlPrefixDbID            // 下载url前缀
)

var (
	ConfigFilterKeyMapDbField = map[string]string{
		fastcurd.PrimaryField:    fastcurd.PrimaryField,
		fastcurd.CreateTimeField: fastcurd.CreateTimeField,
		fastcurd.UpdateTimeField: fastcurd.UpdateTimeField,
	}
	ConfigOrderKeyMapDbField = map[string]string{
		fastcurd.PrimaryField:    fastcurd.PrimaryField,
		fastcurd.CreateTimeField: fastcurd.CreateTimeField,
		fastcurd.UpdateTimeField: fastcurd.UpdateTimeField,
	}
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
func (m *Config) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
func (m *Config) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *Config) TableName() string {
	return "config"
}

func (m *Config) GetFmtDetail(scenes ...string) any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model any
	switch scene {
	default:
		model = NewDefaultSceneConfig(m)
	}
	return model
}
func (m *Config) GetFilterKeyMapDBField() map[string]string {
	return ConfigFilterKeyMapDbField
}
func (m *Config) GetOrderKeyMapDBField() map[string]string {
	return ConfigOrderKeyMapDbField
}
func NewDefaultSceneConfig(m *Config) map[string]any {
	return map[string]any{
		"id":  m.ID,
		"key": m.Key,
		"val": m.Val,
	}
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
