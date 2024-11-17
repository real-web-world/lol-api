package models

import (
	"context"

	"github.com/real-web-world/bdk/fastcurd"
	"gorm.io/gorm"
)

type (
	Tpl struct {
		fastcurd.Base
		Name string `json:"name" gorm:"column:name"`
	}
)

//goland:noinspection GoUnusedExportedFunction
func NewCtxTpl(ctx context.Context, db *gorm.DB) *Tpl {
	return &Tpl{
		Base: fastcurd.Base{
			Ctx: ctx,
			DB:  db,
		},
	}
}
func (m *Tpl) TableName() string {
	return "tpl"
}
