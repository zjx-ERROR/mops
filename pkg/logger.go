package pkg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(level zapcore.Level, pathToLogFile string, maxSize, maxBackups, maxAge int, compress bool) *zap.SugaredLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   pathToLogFile,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}),
		zap.NewAtomicLevelAt(level),
	)
	logger := zap.New(core)
	defer logger.Sync()
	return logger.Sugar()
}
