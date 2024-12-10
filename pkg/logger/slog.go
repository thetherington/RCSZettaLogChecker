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
var logger *slog.Logger

type Logger struct {
	Background bool
	Level      slog.Level
	Filename   string
	ZettaHost  string
	Version    string
}

type Option func(*Logger)

// Set sets the logger configuration based on the environment
func Set(options ...Option) {
	o := &Logger{
		Background: false,
		Level:      slog.LevelDebug,
		Filename:   "log/app.log",
		ZettaHost:  "",
		Version:    "",
	}

	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(o)
	}

	logRotate := &lumberjack.Logger{
		Filename:   o.Filename,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()

	core := ecszap.NewCore(encoderConfig, zapcore.AddSync(logRotate), zap.DebugLevel)
	zapLogger := zap.New(core)
	zapLogger = zapLogger.Named("zettaLogChecker")

	logger = slog.New(
		slogzap.Option{Level: o.Level, Logger: zapLogger, AddSource: false}.NewZapHandler(),
	)

	if !o.Background {
		logger = slog.New(
			slogmulti.Fanout(
				tint.NewHandler(os.Stderr, &tint.Options{
					Level:      o.Level,
					TimeFormat: time.Kitchen,
				}),
				slogzap.Option{Level: o.Level, Logger: zapLogger, AddSource: false}.NewZapHandler(),
			),
		)
	}

	if o.ZettaHost != "" {
		logger = logger.With("Zetta_Host", o.ZettaHost)
	}

	if o.Version != "" {
		logger = logger.With("release", o.Version)
	}

	slog.SetDefault(logger)
}

// WithBackground is a functional option to enable background only logging
func WithBackground() Option {
	return func(l *Logger) {
		l.Background = true
	}
}

// WithFileName is a functional option to specify a file option (default log/app.log)
func WithFileName(fname string) Option {
	return func(l *Logger) {
		l.Filename = fname
	}
}

// WithLevel specifies the logging level
func WithLevel(level slog.Level) Option {
	return func(l *Logger) {
		l.Level = level
	}
}

// WithZettaHost includes the zetta host address to every log
func WithZettaHost(host string) Option {
	return func(l *Logger) {
		l.ZettaHost = host
	}
}

// WithVersion includes the app version to every log
func WithVersion(version string) Option {
	return func(l *Logger) {
		l.Version = version
	}
}
