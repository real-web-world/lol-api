package bootstrap

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/grafana/pyroscope-go"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/real-web-world/bdk"
	bdkgin "github.com/real-web-world/bdk/gin"
	bdkmid "github.com/real-web-world/bdk/gin/middleware"
	"github.com/real-web-world/bdk/valid"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/middleware"
	"github.com/real-web-world/lol-api/routes"
	"github.com/real-web-world/lol-api/services/logger"
	"github.com/real-web-world/lol-api/services/mysql"
	"github.com/real-web-world/lol-api/services/otel"
	"github.com/real-web-world/lol-api/services/rds"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

const (
	envFilename      = ".env"
	envLocalFilename = ".env.local"
	defaultConfPath  = "./config"
	configFile       = "config.json"
)

func initConf() error {
	_ = godotenv.Load(envFilename)
	if bdk.IsFile(envLocalFilename) {
		_ = godotenv.Overload(envLocalFilename)
	}
	confPath := defaultConfPath + "/" + configFile
	return configor.Load(global.Conf, confPath)
}
func initLog() error {
	cfg := global.Conf.Log
	ws := zapcore.AddSync(log.Writer())
	logLevel := zapcore.DebugLevel
	if global.IsProdMode() {
		bufWriter := bufio.NewWriter(log.Writer())
		logWriter := bdk.NewConcurrentWriter(bufWriter)
		log.SetOutput(logWriter)
		global.SetCleanup(global.LogWriterCleanupKey, logWriter.Close)
		ws = zapcore.AddSync(logWriter)
		logLevel, _ = zapcore.ParseLevel(cfg.Level)
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.StringDurationEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		ws,
		zap.NewAtomicLevelAt(logLevel),
	)
	global.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)).Sugar()
	if global.IsProdMode() {
		global.SetCleanup(global.ZapLoggerCleanupKey, func(_ context.Context) error {
			return global.Logger.Sync()
		})
	}
	return nil
}
func initCache(ctx context.Context) error {
	return rds.Init(ctx, &global.Conf.Redis)
}

func initEngine() *gin.Engine {
	gin.SetMode(global.Conf.Mode)
	engine := bdkgin.NewGin()
	binding.Validator = &valid.DefaultValidator{}
	return engine
}
func initMiddleware(e *gin.Engine) {
	Conf := global.Conf
	e.Use(bdkmid.RecoveryWithLogFn(logger.Error))
	//e.Use(gin.LoggerWithConfig(gin.LoggerConfig{
	//	Formatter: bdkgin.LogFormatter,
	//	Output:    log.Writer(),
	//}))
	e.Use(middleware.Prometheus)
	if Conf.Otel.Enabled {
		e.Use(otelgin.Middleware(Conf.ProjectName, otelgin.WithFilter(func(request *http.Request) bool {
			return !bdk.IsSkipLogReq(request, http.StatusOK)
		})))
	}
	e.Use(middleware.Cors(Conf))
	e.Use(bdkmid.NewHttpTraceWithDefaultNotTrace(logger.Info))
}
func InitApp() (*gin.Engine, error) {
	var err error
	if err = initConf(); err != nil {
		return nil, err
	}
	if err = initLog(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return initLib(ctx)
	})
	g.Go(func() error {
		return initSdk(ctx)
	})
	g.Go(func() error {
		return initDB(ctx)
	})
	g.Go(func() error {
		return initCache(ctx)
	})
	e := initEngine()
	initMiddleware(e)
	routes.RouterSetup(e)
	if global.Conf.PProf.Enabled {
		pprof.Register(e.Group("", bdkmid.DevAccess))
	}
	return e, g.Wait()
}
func initDB(ctx context.Context) error {
	return mysql.Init(ctx, &global.Conf.Mysql, global.IsDevMode())
}
func initLib(_ context.Context) error {
	now.WeekStartDay = time.Monday

	return nil
}

func initSdk(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if global.Conf.Otel.Enabled {
			return nil

		}
		shutdown, err := otel.InitOtel(ctx)
		if err != nil {
			err = errors.WithMessage(err, "初始化otel失败")
		} else {
			global.SetCleanup("otel", shutdown)
		}
		return err
	})
	g.Go(func() error {
		if global.Conf.Pyroscope.Enabled {
			return nil
		}
		if err := initPyroscope(global.Conf.Pyroscope); err != nil {
			return errors.WithMessage(err, "初始化pyroscope失败")
		}
		return nil
	})
	return g.Wait()
}

func initPyroscope(cfg conf.PyroscopeConf) error {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	mode := global.Conf.Mode
	isLocal := global.IsLocalDev()
	isLocalStr := "false"
	if isLocal {
		isLocalStr = "true"
	}
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: mode + "_" + cfg.AppName,
		ServerAddress:   cfg.ServerAddress,
		Logger:          nil,
		Tags: map[string]string{
			"buff.hostname": os.Getenv("HOSTNAME"),
			"buff.env":      mode,
			"buff.isLocal":  isLocalStr,
		},
		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	return err
}
