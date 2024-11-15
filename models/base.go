package models

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/real-web-world/lol-api/global"
	"gorm.io/gorm"
)

var (
	spanCtxKey         = global.SpanCtxKey
	gormTraceOpNameKey = global.GormTraceOpNameKey
	nopFn              = func() {}
)

var (
	errCreateFailed = errors.New("创建失败")
)

type Any interface{}
type Base struct {
	ID                 uint64          `json:"id" redis:"id" gorm:"type:bigint unsigned auto_increment;primaryKey;"`
	Ctime              *time.Time      `json:"ctime,omitempty" gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null"`
	Utime              *time.Time      `json:"utime,omitempty" gorm:"type:datetime ON UPDATE CURRENT_TIMESTAMP;default:CURRENT_TIMESTAMP;not null;"`
	RelationAffectRows int             `json:"-" gorm:"-"` // 更新时用来保存其他关联数据的更新数
	Ctx                context.Context `json:"-" gorm:"-"` // ctx
	DB                 *gorm.DB        `json:"-" gorm:"-"` // db
}
type TotalData struct {
	Total int64 `json:"total"`
}

func (m *Base) SetCtx(c context.Context) {
	m.Ctx = c
}
func (m *Base) SetOpName(opName string) func() {
	if m.Ctx != nil {
		m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "db:"+opName)
		return func() {
			m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "")
		}
	}
	return nopFn
}
func (m *Base) SetSpanCtx(ctx opentracing.SpanContext) {
	if m.Ctx != nil {
		m.Ctx = context.WithValue(m.Ctx, spanCtxKey, ctx)
	}
}
