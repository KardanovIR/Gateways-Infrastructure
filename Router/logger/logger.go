package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	With(args ...interface{}) *logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

const loggerKey string = "loggerCtxKey"

type Level string

const (
	DEBUG = Level("DEBUG")
	INFO  = Level("INFO")
	WARN  = Level("WARN")
	ERROR = Level("ERROR")
)

var def ILogger = nil

type logger struct {
	*zap.SugaredLogger
}

func Init(json bool, level Level) (ILogger, error) {
	l, err := newLogger(json, level)
	if err != nil {
		return nil, err
	}
	setDefaultLogger(l)
	return l, nil
}

func ToContext(ctx context.Context, l ILogger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns logger from context with previously added fields.
// If context has not logger returns default logger
func FromContext(ctx context.Context) ILogger {
	if l, ok := ctx.Value(loggerKey).(ILogger); ok {
		return l
	} else {
		return GetDefaultLogger()
	}
}

func GetDefaultLogger() ILogger {
	if def != nil {
		return def
	}
	l, err := newLogger(false, INFO)
	if err != nil {
		panic(fmt.Errorf("can't create default logger: %s", err))
	}
	return l
}

func newLogger(json bool, level Level) (ILogger, error) {

	logConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(convertLevel(level)),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "timestamp",
			CallerKey:    "service",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder, // todo nano seconds
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	if !json {
		logConf.Encoding = "console"
		logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, err := logConf.Build()
	if err != nil {
		return nil, err
	}
	return &logger{l.Sugar()}, nil
}

func convertLevel(l Level) zapcore.Level {
	switch l {
	case ERROR:
		return zapcore.ErrorLevel
	case WARN:
		return zapcore.WarnLevel
	case DEBUG:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func setDefaultLogger(l ILogger) {
	def = l
}

// With creates a child logger and adds a variadic number of fields to the logging context.
// It accepts a mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value..
// Fields added to the child don't affect the parent, and vice versa.
func (l *logger) With(args ...interface{}) *logger {
	return &logger{l.SugaredLogger.With(args)}
}

var _ ILogger = (*logger)(nil)
