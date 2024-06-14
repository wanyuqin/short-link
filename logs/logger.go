package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"short-link/config"
)

var _logger *zap.Logger

func InitializeLogger() {
	var (
		writeSyncer zapcore.WriteSyncer
		logLevel    zapcore.Level
	)
	logCfg := config.GetConfig().Application.Logger
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	switch logCfg.StdType {
	case "std":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "file":
		if logCfg.FilePath == "" {
			panic("log file path cannot be null")
		}
		file, err := os.OpenFile(logCfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		writeSyncer = zapcore.AddSync(file)
	default:
		panic("unknown log std type eg: std/file")
	}
	switch logCfg.Level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel

	}
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	_logger = zap.New(core, zap.AddCaller())
}

func Error(err error, msg string, field ...zap.Field) {
	field = append(field, zap.Any("error", err))
	_logger.Error(msg, field...)
}

func Info(msg string, field ...zap.Field) {
	_logger.Info(msg, field...)
}
func Debug(msg string, field ...zap.Field) {
	_logger.Debug(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	_logger.Warn(msg, field...)
}

func Fatal(msg string, field ...zap.Field) {
	_logger.Fatal(msg, field...)
}
