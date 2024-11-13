package models

import (
	"context"
	"encoding/json"
	"github.com/real-web-world/lol-api/pkg/fastcurd"
	"golang.org/x/sync/errgroup"

	"gorm.io/gorm"
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
		Base
		Key string `json:"key" gorm:"column:k"`
		Val string `json:"Val" gorm:"column:v"`
	}
)

const (
	ClientConfDbID = 1
)

func NewDefaultSceneConfig(m *Config) map[string]any {
	return map[string]any{
		"id":  m.ID,
		"key": m.Key,
		"val": m.Val,
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
func (m *Config) SetOpName(opName string) func() {
	if m.Ctx != nil {
		m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "db:"+opName)
		return func() {
			m.Ctx = context.WithValue(m.Ctx, gormTraceOpNameKey, "")
		}
	}
	return nopFn
}

//goland:noinspection GoUnusedExportedFunction
func NewCtxConfig(ctx context.Context, db *gorm.DB) *Config {
	return &Config{
		Base: Base{
			Ctx: ctx,
			DB:  db,
		},
	}
}
func (m *Config) GetGormQuery() *gorm.DB {
	db := m.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(m)
}
func (m *Config) GetTxGormQuery(tx *gorm.DB) *gorm.DB {
	db := tx
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(m)
}

func (m *Config) GetFmtList(arr []Config, sceneParam ...string) any {
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
func (m *Config) GetDetailByID(id uint64) (*Config, error) {
	record := &Config{}
	err := m.GetGormQuery().Where("id = ?", id).First(record).Error
	if err != nil {
		record = nil
	}
	return record, err
}
func (m *Config) ListByIDArr(idArr []uint64) ([]Config, error) {
	if len(idArr) == 0 {
		return nil, nil
	}
	list := make([]Config, 0, len(idArr))
	err := m.GetGormQuery().Where("id in ?", idArr).Find(&list).Error
	return list, err
}
func (m *Config) dbEditByID(db *gorm.DB, id uint64, values map[string]any) (int64, error) {
	res := db.Where("id = ?", id).Updates(values)
	return res.RowsAffected, res.Error
}
func (m *Config) EditByID(id uint64, values map[string]any) (int64, error) {
	return m.dbEditByID(m.GetGormQuery(), id, values)
}
func (m *Config) TxEditByID(tx *gorm.DB, id uint64, values map[string]any) (int64, error) {
	return m.dbEditByID(m.GetTxGormQuery(tx), id, values)
}
func (m *Config) dbEditByIDArr(db *gorm.DB, idArr []uint64, values map[string]any) (int64, error) {
	res := db.Where("id in ?", idArr).Updates(values)
	return res.RowsAffected, res.Error
}
func (m *Config) EditByIDArr(idArr []uint64, values map[string]any) (int64, error) {
	return m.dbEditByIDArr(m.GetGormQuery(), idArr, values)
}
func (m *Config) TxEditByIDArr(tx *gorm.DB, idArr []uint64, values map[string]any) (int64, error) {
	return m.dbEditByIDArr(m.GetTxGormQuery(tx), idArr, values)
}
func (m *Config) dbDelByIDArr(db *gorm.DB, idArr []uint64) (int64, error) {
	if len(idArr) == 0 {
		return 0, nil
	}
	res := db.Where("id in ?", idArr).Delete(m)
	return res.RowsAffected, res.Error
}
func (m *Config) DelByIDArr(idArr []uint64) (int64, error) {
	return m.dbDelByIDArr(m.GetGormQuery(), idArr)
}
func (m *Config) TxDelByIDArr(tx *gorm.DB, idArr []uint64) (int64, error) {
	return m.dbDelByIDArr(m.GetTxGormQuery(tx), idArr)
}
func (m *Config) dbCreateRecord(db *gorm.DB, record *Config) (*Config, error) {
	res := db.Create(record)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected != 1 {
		return nil, errCreateFailed
	}
	return record, nil
}
func (m *Config) CreateRecord(record *Config) (*Config, error) {
	return m.dbCreateRecord(m.GetGormQuery(), record)
}
func (m *Config) TxCreateRecord(tx *gorm.DB, record *Config) (*Config, error) {
	return m.dbCreateRecord(m.GetTxGormQuery(tx), record)
}
func (m *Config) dbCreateList(db *gorm.DB, list []Config) ([]Config, error) {
	res := db.Create(&list)
	if res.Error != nil {
		return nil, res.Error
	}
	return list, nil
}
func (m *Config) CreateList(list []Config) ([]Config, error) {
	return m.dbCreateList(m.GetGormQuery(), list)
}
func (m *Config) TxCreateList(tx *gorm.DB, list []Config) ([]Config, error) {
	return m.dbCreateList(m.GetTxGormQuery(tx), list)
}
func (m *Config) ListRecord(page int, limit int, filter fastcurd.Filter, order map[string]string) ([]Config, int64, error) {
	var count int64
	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	list := make([]Config, 0, limit)
	db := m.GetGormQuery()
	query, err := fastcurd.BuildFilterCond(ConfigFilterKeyMapDbField, db, filter)
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
		dataQuery = fastcurd.BuildOrderCond(ConfigOrderKeyMapDbField, dataQuery, order)
		return dataQuery.Offset(offset).Limit(limit).Find(&list).Error
	})
	return list, count, g.Wait()
}

func (m *Config) GetClientConf() (*Config, error) {
	return m.GetDetailByID(ClientConfDbID)
}
