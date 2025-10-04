package logger

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger is a logger that uses uber-go/zap.
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger initializes a new zap logger.
func NewZapLogger() (log.Logger, func(), error) {
	config := zap.NewProductionConfig()
	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure output format
	config.Encoding = "json"
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Add service name field
	config.InitialFields = map[string]interface{}{
		"service": "api-gateway",
	}

	logger, err := config.Build()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		logger.Sync()
	}

	return &zapLogger{logger: logger}, cleanup, nil
}

// Log implements the Kratos logger interface.
func (l *zapLogger) Log(level log.Level, keyvals ...interface{}) error {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.logger.Warn("Keyvalues must appear in pairs", zap.Any("keyvals", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(keyvals[i].(string), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.logger.Debug("", data...)
	case log.LevelInfo:
		l.logger.Info("", data...)
	case log.LevelWarn:
		l.logger.Warn("", data...)
	case log.LevelError:
		l.logger.Error("", data...)
	case log.LevelFatal:
		l.logger.Fatal("", data...)
	}
	return nil
}

var ProviderSet = wire.NewSet(NewZapLogger)
