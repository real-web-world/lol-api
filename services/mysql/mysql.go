package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/real-web-world/bdk/logger"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

func Init(ctx context.Context, cfg *conf.MysqlConf, isDebug bool) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		db, err := initDb(ctx, cfg.Default, "default", isDebug)
		if err != nil {
			return err
		}
		global.DefaultDB = db
		return nil
	})
	return g.Wait()
}
func initDb(ctx context.Context, cfg conf.MysqlItemConf, name string, isDebug bool) (*gorm.DB, error) {
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
		return nil, errors.Wrap(err, fmt.Sprintf("init %s db otel failed", name))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIDleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTimeMinutes) * time.Minute)
	return db, sqlDB.PingContext(ctx)
}
