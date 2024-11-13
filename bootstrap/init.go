package bootstrap

import (
	"bufio"
	"context"
	"fmt"
	"github.com/grafana/pyroscope-go"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/routes"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	apiProj "github.com/real-web-world/lol-api"
	"github.com/real-web-world/lol-api/global"
	mid "github.com/real-web-world/lol-api/middleware"
	"github.com/real-web-world/lol-api/pkg/bdk"
	ginApp "github.com/real-web-world/lol-api/pkg/gin"
	"github.com/real-web-world/lol-api/pkg/logger"
	"github.com/real-web-world/lol-api/pkg/valid"
	"github.com/real-web-world/lol-api/services/mysql"
	"github.com/real-web-world/lol-api/services/rds"
	"github.com/real-web-world/lol-api/services/trace"
)

const (
	defaultTZ        = "Asia/Shanghai"
	envFilename      = ".env"
	envLocalFilename = ".env.local"
	tzEnv            = "TZ"
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
	ws := zapcore.AddSync(os.Stdout)
	var logLevel logger.LogLevelStr
	if global.IsDevMode() {
		logLevel = logger.LevelDebugStr
	} else {
		bufWriter := bufio.NewWriter(os.Stdout)
		logWriter := bdk.NewConcurrentWriter(bufWriter)
		logWriter.Consume()
		log.SetOutput(logWriter)
		global.SetCleanup(global.LogWriterCleanupKey, logWriter.Flush)
		ws = zapcore.AddSync(logWriter)
		logLevel = cfg.Level
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.StringDurationEncoder
	level, err := logger.Str2ZapLevel(logLevel)
	if err != nil {
		return err
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		ws,
		zap.NewAtomicLevelAt(level),
	)
	global.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	if global.IsProdMode() {
		global.SetCleanup(global.ZapLoggerCleanupKey, global.Logger.Sync)
	}
	return nil
}
func initCache(ctx context.Context) error {
	return rds.Init(ctx, &global.Conf.Redis)
}

func logFormatter(p gin.LogFormatterParams) string {
	if bdk.IsSkipLogReq(p.Request, p.StatusCode) {
		return ""
	}
	reqTime := p.TimeStamp.Format(time.DateTime)
	path := p.Request.URL.Path
	method := p.Request.Method
	code := p.StatusCode
	clientIp := p.ClientIP
	errMsg := p.ErrorMessage
	processTime := p.Latency
	return fmt.Sprintf("API: %s %d %s %s %s %v %s\n", reqTime, code, clientIp, path, method, processTime,
		errMsg)
}
func initRegValidates() {
	for _, fn := range global.ValidFuncList {
		fn()
	}
}
func initEngine() *gin.Engine {
	gin.SetMode(global.Conf.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logFormatter))
	defaultValid := &valid.DefaultValidator{}
	binding.Validator = defaultValid
	initRegValidates()
	return engine
}
func initMiddleware(e *gin.Engine) {
	Conf := global.Conf
	// sentry
	e.Use(mid.GenerateSentryMiddleware(Conf))
	e.Use(ginApp.PrepareProc)
	e.Use(mid.Prometheus)
	// trace
	if Conf.Trace.Enabled {
		e.Use(otelgin.Middleware(Conf.ProjectName, otelgin.WithFilter(func(request *http.Request) bool {
			return !bdk.IsSkipLogReq(request, http.StatusOK)
		})))
	}
	// httpTrace
	e.Use(mid.HTTPTrace)
	e.Use(mid.InjectRemoteIP)
	// csrf
	e.Use(mid.Cors(Conf))
	e.Use(mid.Logger)
	e.Use(mid.ErrHandler)
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
	return e, g.Wait()
}
func initDB(ctx context.Context) error {
	return mysql.Init(ctx, &global.Conf.Mysql, global.IsDevMode())
}
func initLib(_ context.Context) error {
	var err error
	_ = os.Setenv(tzEnv, defaultTZ)
	now.WeekStartDay = time.Monday
	if global.Conf.Trace.Enabled {
		if err = trace.InitJeager(&global.Conf.Trace); err != nil {
			return errors.WithMessage(err, "初始化jaeger失败")
		}
	}
	if global.Conf.Sentry.Enabled {
		if err = initSentry(global.Conf.Sentry.Dsn); err != nil {
			return errors.WithMessage(err, "Sentry initialization failed")
		}
	}
	if global.Conf.PProf.Enabled {
		initPprof(global.Conf.PProf.Addr)
	}
	if global.Conf.Pyroscope.Enabled {
		if err = initPyroscope(global.Conf.Pyroscope); err != nil {
			return errors.WithMessage(err, "初始化pyroscope失败")
		}
	}
	rand.Seed(time.Now().UnixNano())
	return nil
}
func initPprof(addr string) {
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println("pprof start failed:", err)
		}
	}()
}

func initSentry(dsn string) error {
	isDebugMode := global.IsDevMode()
	sampleRate := 1.0
	if !isDebugMode {
		sampleRate = 1.0
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Debug:       isDebugMode,     // 是否是测试环境
		SampleRate:  sampleRate,      // 固定 1.0
		Release:     apiProj.Commit,  // 测试环境传dev  正式环境传 git commit id
		Environment: global.GetEnv(), // 环境变量 [beta,debug] [prod,release]
	})
	if err == nil {
		global.SetCleanup("sentryFlush", func() error {
			sentry.Flush(2 * time.Second)
			return nil
		})
	}
	return err
}

func initSdk(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return nil
	})
	return g.Wait()
}

func initPyroscope(cfg conf.PyroscopeConf) error {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	envInfo := global.GetEnv()
	isLocal := global.IsLocalDev()
	isLocalStr := "false"
	if isLocal {
		isLocalStr = "true"
	}
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: envInfo + "_" + cfg.AppName,
		// replace this with the address of pyroscope server
		ServerAddress: cfg.ServerAddress,
		// you can disable logging by setting this to nil
		//Logger: pyroscope.StandardLogger,
		Logger: nil,
		// you can provide static tags via a map:
		Tags: map[string]string{
			"hostname": os.Getenv("HOSTNAME"),
			"env":      global.GetEnv(),
			"isLocal":  isLocalStr,
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
