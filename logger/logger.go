package logger

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger struct {
	LogLevel  string
	SentryDns string
	Zap       *zap.Logger
}

func NewLogger(logLevel, sentryDsn string) *Logger {
	logger := Logger{
		LogLevel:  logLevel,
		SentryDns: sentryDsn,
	}
	logger.InitLogger()
	return &logger
}

func (l *Logger) InitLogger() {
	if l.SentryDns != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: l.SentryDns,
		})
	}
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch l.LogLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}
	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "zap",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	l.Zap, err = zapConfig.Build()
	if err != nil {
		panic(err)
	}
	defer l.Zap.Sync()
	stdLog := zap.RedirectStdLog(l.Zap)
	defer stdLog()
}

func (l *Logger) SendError(err error) {
	if l.SentryDns != "" {
		sentry.CaptureException(err)
		sentry.Flush(time.Second * 5)
	}
}
