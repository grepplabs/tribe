package log

import (
	"context"
	"os"
)

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//DebugLevel has verbose message
	DebugLevel = "debug"
	//InfoLevel is default log level
	InfoLevel = "info"
	//WarnLevel is for logging messages about possible issues
	WarnLevel = "warn"
	//WithError is for logging errors
	ErrorLevel = "error"
	// PanicLevel log a message and panic.
	PanicLevel = "panic"
	//FatalLevel is for logging fatal messages. The system shuts down after logging the message.
	FatalLevel = "fatal"
)

const (
	//TimeKey is a logger key for time
	TimeKey = "ts"
	//MessageKey is a logger key for message
	MessageKey = "msg"
	//LevelKey is a logger key for logging level
	LevelKey = "level"
	//CallerKey ia a logger key for caller/invoking function
	CallerKey = "caller"
)

const (
	// LogFormatJson is a format for json logging
	LogFormatJson = "json"
	// LogFormatPlain is a format for plain-text logging
	LogFormatPlain = "plain"
	// LogFormatLogfmt is a format for logfmt logging
	LogFormatLogfmt = "logfmt"
)

const (
	EnvLogFormat           = "LOG_FORMAT"
	EnvLogLevel            = "LOG_LEVEL"
	EnvLogFieldNameTime    = "LOG_FIELD_NAME_TIME"
	EnvLogFieldNameMessage = "LOG_FIELD_NAME_MESSAGE"
	EnvLogFieldNameCaller  = "LOG_FIELD_NAME_CALLER"
	EnvLogFieldNameLevel   = "LOG_FIELD_NAME_LEVEL"
)

//Logger is our contract for the logger
type Logger interface {
	Write(p []byte) (n int, err error)

	Printf(format string, args ...interface{})

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	WithField(key, value string) Logger

	WithError(err error) Logger

	WithName(name string) Logger

	IsDebug() bool

	IsInfo() bool

	IsWarn() bool

	IsError() bool

	IsPanic() bool

	IsFatal() bool

	Info(msg string, keysAndValues ...interface{})

	Debug(msg string, keysAndValues ...interface{})

	Error(err error, msg string, keysAndValues ...interface{})
}

type LogFieldNames struct {
	Time    string
	Message string
	Level   string
	Caller  string
}

// Configuration stores the config for the logger
type Configuration struct {
	LogFormat     string
	LogLevel      string
	LogFieldNames LogFieldNames
}

//NewLogger returns an instance of logger
func NewLogger(config Configuration) Logger {
	return newZapLogger(config)
}

var DefaultLogger = NewDefaultLogger()

//NewDefaultLogger returns an instance of logger with default parameters
func NewDefaultLogger() Logger {
	config := Configuration{
		LogFormat: getEnv(EnvLogFormat, LogFormatLogfmt),
		LogLevel:  getEnv(EnvLogLevel, InfoLevel),
		LogFieldNames: LogFieldNames{
			Time:    getEnv(EnvLogFieldNameTime, TimeKey),
			Message: getEnv(EnvLogFieldNameMessage, MessageKey),
			Level:   getEnv(EnvLogFieldNameLevel, LevelKey),
			Caller:  getEnv(EnvLogFieldNameCaller, CallerKey),
		},
	}
	return newZapLogger(config)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	DefaultLogger.Panicf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return DefaultLogger.WithFields(keyValues)
}

func WithField(key, value string) Logger {
	return DefaultLogger.WithField(key, value)
}

func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

func WithName(name string) Logger {
	return DefaultLogger.WithName(name)
}

func Info(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Info(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debug(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	DefaultLogger.Error(err, msg, keysAndValues...)
}

type contextKey struct{}

// FromContext returns a Logger constructed from ctx or nil if no logger details are found.
func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return nil
}

// FromContextOrDiscard returns a Logger constructed from ctx or a default logger if no logger details are found.
func FromContextOrDefault(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return DefaultLogger
}

// NewContext returns a new context derived from ctx that embeds the Logger.
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
