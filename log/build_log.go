package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"conan/config"
	"conan/utils"
)

type LogLEVEL = string

var (
	DEBUG LogLEVEL = "DEBUG"
	INFO  LogLEVEL = "INFO"
	WARN  LogLEVEL = "WARN"
	ERROR LogLEVEL = "ERROR"
	PANIC LogLEVEL = "PANIC"
	FATAL LogLEVEL = "FATAL"
)

const (
	JSON   = "JSON"
	STDOUT = "CONSOLE"
)

type LogConfig struct {
	OutFile   bool   `ini:"outFile"`
	Level     string `ini:"level"`
	LogFormat string `ini:"logFormat"`
	FilePath  string `ini:"filePath"`
}

func getDefaultLog() LogConfig {
	cfg := LogConfig{}
	err := config.Decode(config.CONFIG_FILE_NAME, "log", &cfg)
	if err != nil {
		panic("build log config fail package is log , method is getDefaultLog")
	}
	return cfg
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case PANIC:
		return zapcore.PanicLevel
	case FATAL:
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}

func buildLog(config LogConfig) *zap.Logger {

	/*

				TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	*/

	coreCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		NameKey:        "logger",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	autoLevel := zap.NewAtomicLevel()
	autoLevel.SetLevel(getLogLevel(config.Level))

	encoder := getLogEncoder(config.LogFormat, coreCfg)

	writeSync := getWriteSync(config.OutFile, config.FilePath)
	core := zapcore.NewCore(encoder, writeSync, autoLevel)
	return zap.New(core)
}

func buildLog_v2(config LogConfig) *zap.SugaredLogger {
	coreCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	autoLevel := zap.NewAtomicLevel()
	autoLevel.SetLevel(getLogLevel(config.Level))

	encoder := getLogEncoder(config.LogFormat, coreCfg)

	writeSync := getWriteSync(config.OutFile, config.FilePath)
	core := zapcore.NewCore(encoder, writeSync, autoLevel)
	log := zap.New(core)
	log.Sugar()

	return log.Sugar()
}

func getWriteSync(outFile bool, filePath string) zapcore.WriteSyncer {
	if outFile {
		if filePath == "" {
			filePath = utils.GetProDirAbs() + "logs" + string(filepath.Separator) + "conan.log"
			//filePath = "../logs/conan.log"
		}
		hook := lumberjack.Logger{
			Filename:   filePath, // 日志文件路径
			MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 30,       // 日志文件最多保存多少个备份
			MaxAge:     7,        // 文件最多保存多少天
			Compress:   true,     // 是否压缩
		}
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

}

func getLogEncoder(logFormat string, cfg zapcore.EncoderConfig) zapcore.Encoder {
	if logFormat == JSON {
		return zapcore.NewJSONEncoder(cfg)
	} else if logFormat == STDOUT {
		return zapcore.NewConsoleEncoder(cfg)
	}
	return zapcore.NewJSONEncoder(cfg)
}
