package mysql

import (
	"context"
	"fmt"
	gormLogger "gorm.io/gorm/logger"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/pkg/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
// spanCtxKey         = global.SpanCtxKey
// gormTraceOpNameKey = global.GormTraceOpNameKey
)

func Init(ctx context.Context, cfg *conf.MysqlConf, isDebug bool) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		db, err := initDb(cfg.Default, "default", isDebug)
		if err != nil {
			return err
		}
		global.DefaultDB = db
		return nil
	})
	return g.Wait()
}
func initDb(cfg conf.MysqlItemConf, name string, isDebug bool) (*gorm.DB, error) {
	var l = logger.GormLogger
	if isDebug {
		l = l.LogMode(gormLogger.Info)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.UserName, cfg.Pwd, cfg.Host, cfg.Port, cfg.Database, cfg.Charset),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   l,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("init %s db connect failed", name))
	}
	err = db.Use(tracing.NewPlugin())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("init %s db trace failed", name))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIDleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Minute)
	return db, nil
}
