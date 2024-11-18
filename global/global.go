package global

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/real-web-world/lol-api/conf"
)

const (
	LocalDevKey = "LocalDev"
)

const (
	LogWriterCleanupKey = "logWriter"
	ZapLoggerCleanupKey = "zapLogger"
	OtelCleanupKey      = "otel"
)

var (
	Conf   = &conf.AppConf{}
	Logger *zap.SugaredLogger
	// JaegerHttpRT  = otelhttp.NewTransport(http.DefaultTransport)
	Cleanups   = make(map[string]func(ctx context.Context) error)
	cleanupsMu = sync.Mutex{}
)

// Redis
var (
	DefaultCe *redis.Client
	_         = DefaultCe
)

// Mysql
var (
	DefaultDB *gorm.DB
)

func IsDevMode() bool {
	return Conf.Mode == gin.DebugMode
}
func IsProdMode() bool {
	return !IsDevMode()
}
func IsLocalDev() bool {
	return os.Getenv(LocalDevKey) == "true"
}

func Cleanup() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	for name, cleanup := range Cleanups {
		if name == LogWriterCleanupKey {
			continue
		}
		if err := cleanup(ctx); err != nil {
			log.Printf("%s cleanup err:%v\n", name, err)
		}
	}
	if fn, ok := Cleanups[LogWriterCleanupKey]; ok {
		_ = fn(ctx)
	}
}
func SetCleanup(name string, fn func(ctx context.Context) error) {
	cleanupsMu.Lock()
	Cleanups[name] = fn
	cleanupsMu.Unlock()
}
