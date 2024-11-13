package global

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"

	"github.com/real-web-world/lol-api/conf"
)

const (
	LocalDevKey = "LocalDev"
)

const (
	FaviconReq         = "/favicon"
	VersionApi         = "/version"
	MetricsApi         = "/metrics"
	StatusApi          = "/status"
	DebugApiPrefix     = "/debug"
	SwaggerApiPrefix   = "/swagger"
	SpanCtxKey         = "spanCtx"
	GormTraceOpNameKey = "gormTraceOpName"
	//RedisTraceOpNameKey = "redisTraceOpName"
	//HttpTraceOpNameKey  = "httpTraceOpName"
	envKey              = "Mode"
	LogWriterCleanupKey = "logWriter"
	ZapLoggerCleanupKey = "zapLogger"
	JaegerCleanupKey    = "jaeger"
)

var (
	Conf          = &conf.AppConf{}
	ValidFuncList []func()
	Logger        *zap.SugaredLogger
	// JaegerHttpRT  = otelhttp.NewTransport(http.DefaultTransport)
	Cleanups   = make(map[string]func() error)
	cleanupsMu = sync.Mutex{}
	currEnv    *string
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
	return !IsProdMode()
}
func IsProdMode() bool {
	return GetEnv() == gin.ReleaseMode
}
func IsLocalDev() bool {
	return os.Getenv(LocalDevKey) == "true"
}

func GetEnv() string {
	if currEnv == nil {
		currEnv = new(string)
		*currEnv = os.Getenv(envKey)
	}
	return *currEnv
}
func Cleanup() {
	for name, cleanup := range Cleanups {
		if name == LogWriterCleanupKey {
			continue
		}
		if err := cleanup(); err != nil {
			log.Printf("%s cleanup err:%v\n", name, err)
		}
	}
	if fn, ok := Cleanups[LogWriterCleanupKey]; ok {
		_ = fn()
	}
}
func SetCleanup(name string, fn func() error) {
	cleanupsMu.Lock()
	Cleanups[name] = fn
	cleanupsMu.Unlock()
}
