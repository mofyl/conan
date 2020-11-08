package log

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	cfg := getDefaultLog()
	logger = buildLog(cfg)
	logger.Sync()
}

func Debug(format string, args ...interface{}) {
	defaultLog(DEBUG, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	defaultLog(INFO, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	defaultLog(WARN, fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	defaultLog(ERROR, fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	defaultLog(PANIC, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	defaultLog(FATAL, fmt.Sprintf(format, args...))
}

func CloseLog() error {

	if logger == nil {
		return nil
	}
	return logger.Sync()
}

func defaultLog(level string, msg string, field ...zap.Field) {

	if ce := logger.Check(getLogLevel(level), msg); ce != nil {
		ce.Write(field...)
	}

}
