package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// In a real enterprise system, log files would be rotated or streamed to ELK/Datadog.
	logFile, _ := os.OpenFile("enterprise-erp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.InfoLevel
	if os.Getenv("ENV") == "development" {
		defaultLogLevel = zapcore.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(log) // Make it globally accessible via zap.L()
}

func Get() *zap.Logger {
	if log == nil {
		InitLogger()
	}
	return log
}
