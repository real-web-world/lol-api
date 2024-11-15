package models

import (
	"context"

	"github.com/real-web-world/bdk/fastcurd"
	"github.com/real-web-world/bdk/json"
	"gorm.io/gorm"
)

var (
	TplFilterKeyMapDbField = map[string]string{
		fastcurd.PrimaryField:    fastcurd.PrimaryField,
		fastcurd.CreateTimeField: fastcurd.CreateTimeField,
		fastcurd.UpdateTimeField: fastcurd.UpdateTimeField,
	}
	TplOrderKeyMapDbField = map[string]string{
		fastcurd.PrimaryField:    fastcurd.PrimaryField,
		fastcurd.CreateTimeField: fastcurd.CreateTimeField,
		fastcurd.UpdateTimeField: fastcurd.UpdateTimeField,
	}
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
func (m *Tpl) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
func (m *Tpl) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *Tpl) TableName() string {
	return "Tpl"
}

func (m *Tpl) GetFmtDetail(scenes ...string) any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model any
	switch scene {
	default:
		model = NewDefaultSceneTpl(m)
	}
	return model
}
func (m *Tpl) GetFilterKeyMapDBField() map[string]string {
	return TplFilterKeyMapDbField
}
func (m *Tpl) GetOrderKeyMapDBField() map[string]string {
	return TplOrderKeyMapDbField
}
func NewDefaultSceneTpl(m *Tpl) map[string]any {
	return map[string]any{
		"id":   m.ID,
		"name": m.Name,
	}
}
