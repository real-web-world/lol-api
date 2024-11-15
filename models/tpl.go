package models

import (
	"context"
	"encoding/json"

	"github.com/real-web-world/bdk/fastcurd"
	"golang.org/x/sync/errgroup"
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
		Base
		Name string `json:"name"`
	}
)

func NewDefaultSceneTpl(m *Tpl) map[string]any {
	return map[string]any{
		"id":   m.ID,
		"name": m.Name,
		//"ctime": m.Ctime,
		//"utime": m.Utime,
	}
}

func (m *Tpl) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
func (m *Tpl) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *Tpl) TableName() string {
	return "tpl"
}
func (m *Tpl) SetOpName(opName string) func() {
	if m.Ctx != nil {
		m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "db:"+opName)
		return func() {
			m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "")
		}
	}
	return nopFn
}

//goland:noinspection GoUnusedExportedFunction
func NewCtxTpl(ctx context.Context, db *gorm.DB) *Tpl {
	return &Tpl{
		Base: Base{
			Ctx: ctx,
			DB:  db,
		},
	}
}
func (m *Tpl) GetGormQuery() *gorm.DB {
	db := m.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(m)
}
func (m *Tpl) GetTxGormQuery(tx *gorm.DB) *gorm.DB {
	db := tx
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(m)
}

func (m *Tpl) GetFmtList(arr []Tpl, sceneParam ...string) any {
	scene := ""
	if len(sceneParam) > 0 {
		scene = sceneParam[0]
	}
	fmtList := make([]any, 0, len(arr))
	actList := arr
	for _, item := range actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
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
func (m *Tpl) GetDetailByID(id uint64) (*Tpl, error) {
	record := &Tpl{}
	err := m.GetGormQuery().Where("id = ?", id).First(record).Error
	if err != nil {
		record = nil
	}
	return record, err
}
func (m *Tpl) ListByIDArr(idArr []uint64) ([]Tpl, error) {
	if len(idArr) == 0 {
		return nil, nil
	}
	list := make([]Tpl, 0, len(idArr))
	err := m.GetGormQuery().Where("id in ?", idArr).Find(&list).Error
	return list, err
}
func (m *Tpl) dbEditByID(db *gorm.DB, id uint64, values map[string]any) (int64, error) {
	res := db.Where("id = ?", id).Updates(values)
	return res.RowsAffected, res.Error
}
func (m *Tpl) EditByID(id uint64, values map[string]any) (int64, error) {
	return m.dbEditByID(m.GetGormQuery(), id, values)
}
func (m *Tpl) TxEditByID(tx *gorm.DB, id uint64, values map[string]any) (int64, error) {
	return m.dbEditByID(m.GetTxGormQuery(tx), id, values)
}
func (m *Tpl) dbEditByIDArr(db *gorm.DB, idArr []uint64, values map[string]any) (int64, error) {
	res := db.Where("id in ?", idArr).Updates(values)
	return res.RowsAffected, res.Error
}
func (m *Tpl) EditByIDArr(idArr []uint64, values map[string]any) (int64, error) {
	return m.dbEditByIDArr(m.GetGormQuery(), idArr, values)
}
func (m *Tpl) TxEditByIDArr(tx *gorm.DB, idArr []uint64, values map[string]any) (int64, error) {
	return m.dbEditByIDArr(m.GetTxGormQuery(tx), idArr, values)
}
func (m *Tpl) dbDelByIDArr(db *gorm.DB, idArr []uint64) (int64, error) {
	if len(idArr) == 0 {
		return 0, nil
	}
	res := db.Where("id in ?", idArr).Delete(m)
	return res.RowsAffected, res.Error
}
func (m *Tpl) DelByIDArr(idArr []uint64) (int64, error) {
	return m.dbDelByIDArr(m.GetGormQuery(), idArr)
}
func (m *Tpl) TxDelByIDArr(tx *gorm.DB, idArr []uint64) (int64, error) {
	return m.dbDelByIDArr(m.GetTxGormQuery(tx), idArr)
}
func (m *Tpl) dbCreateRecord(db *gorm.DB, record *Tpl) (*Tpl, error) {
	res := db.Create(record)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected != 1 {
		return nil, errCreateFailed
	}
	return record, nil
}
func (m *Tpl) CreateRecord(record *Tpl) (*Tpl, error) {
	return m.dbCreateRecord(m.GetGormQuery(), record)
}
func (m *Tpl) TxCreateRecord(tx *gorm.DB, record *Tpl) (*Tpl, error) {
	return m.dbCreateRecord(m.GetTxGormQuery(tx), record)
}
func (m *Tpl) dbCreateList(db *gorm.DB, list []Tpl) ([]Tpl, error) {
	res := db.Create(&list)
	if res.Error != nil {
		return nil, res.Error
	}
	return list, nil
}
func (m *Tpl) CreateList(list []Tpl) ([]Tpl, error) {
	return m.dbCreateList(m.GetGormQuery(), list)
}
func (m *Tpl) TxCreateList(tx *gorm.DB, list []Tpl) ([]Tpl, error) {
	return m.dbCreateList(m.GetTxGormQuery(tx), list)
}
func (m *Tpl) ListRecord(page int, limit int, filter fastcurd.Filter, order map[string]string) ([]Tpl, int64, error) {
	var count int64
	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	list := make([]Tpl, 0, limit)
	db := m.GetGormQuery()
	query, err := fastcurd.BuildFilterCond(TplFilterKeyMapDbField, db, filter)
	if err != nil {
		return nil, count, err
	}
	dataQuery := query.WithContext(m.Ctx)
	g := errgroup.Group{}
	g.Go(func() error {
		query.Count(&count)
		return nil
	})
	g.Go(func() error {
		dataQuery = fastcurd.BuildOrderCond(TplOrderKeyMapDbField, dataQuery, order)
		return dataQuery.Offset(offset).Limit(limit).Find(&list).Error
	})
	return list, count, g.Wait()
}
