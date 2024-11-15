package rds

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const (
	defaultClientName = "default"
)

func getClient(ctx context.Context, cfg conf.RedisItemConf, name string) (*redis.Client, error) {
	ce := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Pwd,
		PoolSize: cfg.Pool,
		DB:       cfg.DB,
	})
	if err := redisotel.InstrumentTracing(ce); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("初始化 %s redis InstrumentTracing 失败\n", name))
	}
	if err := redisotel.InstrumentMetrics(ce); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("初始化 %s redis InstrumentMetrics 失败\n", name))
	}
	_, err := ce.Ping(ctx).Result()
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("初始化 %s redis ping 失败\n", name))
	}
	return ce, nil
}
func Init(ctx context.Context, cfg *conf.RedisConf) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		ce, err := getClient(ctx, cfg.Default, defaultClientName)
		if err != nil {
			return err
		}
		global.DefaultCe = ce
		return nil
	})
	return g.Wait()
}
