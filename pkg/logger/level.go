package logger

import "go.uber.org/zap/zapcore"

type LogLevelStr string

const (
	LevelDebugStr LogLevelStr = "debug"
)

func Str2ZapLevel(level LogLevelStr) (zapcore.Level, error) {
	l := zapcore.DebugLevel
	return l, l.Set(string(level))
}
