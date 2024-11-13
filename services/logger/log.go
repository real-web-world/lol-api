package logger

import (
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/pkg/bdk"
	"math"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ = Debug
	_ = Info
	_ = Warn
	_ = Error
)

func Debug(msg string, keysAndValues ...interface{}) {
	global.Logger.Debugw(msg, keysAndValues...)
}
func Info(msg string, keysAndValues ...interface{}) {
	global.Logger.Infow(msg, keysAndValues...)
}
func Warn(msg string, keysAndValues ...interface{}) {
	go sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelWarning)
		scope.SetExtras(zapFields2SentryExtra(keysAndValues))
		sentry.CaptureMessage(msg)
	})
	global.Logger.Warnw(msg, keysAndValues...)
}
func Error(msg string, keysAndValues ...interface{}) {
	go sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.SetExtras(zapFields2SentryExtra(keysAndValues))
		sentry.CaptureMessage(msg)
	})
	global.Logger.Errorw(msg, keysAndValues...)
}
func zapFields2SentryExtra(fields ...any) map[string]any {
	extraData := map[string]any{}
	for _, v := range fields {
		f, ok := v.(zap.Field)
		if !ok {
			continue
		}
		key := f.Key
		var val any
		switch f.Type {
		case zapcore.BoolType:
			val = f.Integer == 1
		case zapcore.ByteStringType:
			val = bdk.Bytes2Str(f.Interface.([]byte))
		case zapcore.DurationType:
			val = time.Duration(f.Integer)
		case zapcore.Float64Type:
			val = math.Float64frombits(uint64(f.Integer))
		case zapcore.Float32Type:
			val = math.Float32frombits(uint32(f.Integer))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type,
			zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
			val = f.Integer
		case zapcore.StringType:
			val = f.String
		case zapcore.TimeType:
			if f.Interface != nil {
				val = time.Unix(0, f.Integer).In(f.Interface.(*time.Location))
			} else {
				val = time.Unix(0, f.Integer)
			}
		case zapcore.TimeFullType:
			val = f.Interface.(time.Time)
		case zapcore.UintptrType:
			val = uintptr(f.Integer)
		case zapcore.StringerType:
			val = f.Interface
		case zapcore.ErrorType:
			val = f.Interface.(error)
		default:
			val = "unknown log val"
		}
		extraData[key] = val
	}
	return extraData
}
