package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	slogzap "github.com/samber/slog-zap/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// logger is the default logger used by the application
var Logger *slog.Logger

// Set sets the logger configuration based on the environment
func Set(background bool, level slog.Level) {
	logRotate := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()

	core := ecszap.NewCore(encoderConfig, zapcore.AddSync(logRotate), zap.DebugLevel)
	zapLogger := zap.New(core)
	zapLogger = zapLogger.Named("zettaLogChecker")

	Logger = slog.New(
		slogzap.Option{Level: level, Logger: zapLogger, AddSource: false}.NewZapHandler(),
	)

	if !background {
		Logger = slog.New(
			slogmulti.Fanout(
				tint.NewHandler(os.Stderr, &tint.Options{
					Level:      level,
					TimeFormat: time.Kitchen,
				}),
				slogzap.Option{Level: level, Logger: zapLogger, AddSource: false}.NewZapHandler(),
			),
		)
	}

	slog.SetDefault(Logger)
}
