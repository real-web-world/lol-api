package pgsql

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"time"

	"github.com/pkg/errors"
	"github.com/real-web-world/bdk/logger"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

func Init(ctx context.Context, cfg *conf.DbConf, isDebug bool) error {
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
func initDb(ctx context.Context, cfg conf.DbItemConf, name string, isDebug bool) (*gorm.DB, error) {
	var l = logger.GormLogger
	if isDebug {
		l = l.LogMode(gormLogger.Info)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			cfg.Host, cfg.UserName, cfg.Pwd, cfg.Database, cfg.Port, cfg.Tz),
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
	sqlDB.SetConnMaxIdleTime(60 * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTimeMinutes) * time.Minute)
	return db, nil
}
