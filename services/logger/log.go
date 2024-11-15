package logger

import (
	"github.com/real-web-world/lol-api/global"
	"go.uber.org/zap/zapcore"
)

var (
	_ = Debug
	_ = Info
	_ = Warn
	_ = Error
)

func Debug(msg string, keysAndValues ...interface{}) {
	log(zapcore.DebugLevel, msg, keysAndValues...)
}
func Info(msg string, keysAndValues ...interface{}) {
	log(zapcore.InfoLevel, msg, keysAndValues...)
}
func Warn(msg string, keysAndValues ...interface{}) {
	log(zapcore.WarnLevel, msg, keysAndValues...)
}
func Error(msg string, keysAndValues ...interface{}) {
	log(zapcore.ErrorLevel, msg, keysAndValues...)
}
func log(lvl zapcore.Level, msg string, keysAndValues ...any) {
	global.Logger.Logw(lvl, msg, keysAndValues...)
}
