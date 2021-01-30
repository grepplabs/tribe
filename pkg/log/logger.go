package log

import "os"

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//WithError is for logging errors
	Error = "error"
	// Panic log a message and panic.
	Panic = "panic"
	//Fatal is for logging fatal messages. The system shuts down after logging the message.
	Fatal = "fatal"
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
	// ErrorKey is a logger key for message
	ErrorKey = "err"
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
	EnvLogFieldNameError   = "LOG_FIELD_NAME_ERROR"
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

	IsDebug() bool

	IsInfo() bool

	IsWarn() bool

	IsError() bool

	IsPanic() bool

	IsFatal() bool
}

type LogFieldNames struct {
	Time    string
	Message string
	Level   string
	Caller  string
	Error   string
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
		LogLevel:  getEnv(EnvLogLevel, Info),
		LogFieldNames: LogFieldNames{
			Time:    getEnv(EnvLogFieldNameTime, TimeKey),
			Message: getEnv(EnvLogFieldNameMessage, MessageKey),
			Level:   getEnv(EnvLogFieldNameLevel, LevelKey),
			Caller:  getEnv(EnvLogFieldNameCaller, CallerKey),
			Error:   getEnv(EnvLogFieldNameError, ErrorKey),
		},
	}
	return newZapLogger(config)
}

func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
