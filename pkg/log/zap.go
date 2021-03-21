package log

import (
	"os"

	"github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	//errorKey is a logger key for error
	errorKey = "error"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
	level         zapcore.Level
}

func getEncoder(config Configuration) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = getStringOrDefault(config.LogFieldNames.Time, TimeKey)
	encoderConfig.MessageKey = getStringOrDefault(config.LogFieldNames.Message, MessageKey)
	encoderConfig.LevelKey = getStringOrDefault(config.LogFieldNames.Level, LevelKey)
	encoderConfig.CallerKey = getStringOrDefault(config.LogFieldNames.Caller, CallerKey)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch config.LogFormat {
	case LogFormatJson:
		return zapcore.NewJSONEncoder(encoderConfig)
	case LogFormatLogfmt:
		return zaplogfmt.NewEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func getStringOrDefault(value, defValue string) string {
	if value != "" {
		return value
	}
	return defValue
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration) Logger {
	cores := []zapcore.Core{}

	level := getZapLevel(config.LogLevel)
	writer := zapcore.Lock(os.Stderr)
	core := zapcore.NewCore(getEncoder(config), writer, level)
	cores = append(cores, core)

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 1 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
		level:         level,
	}
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) IsDebug() bool {
	return l.level == zapcore.DebugLevel
}

func (l *zapLogger) Printf(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) IsInfo() bool {
	return l.level <= zapcore.InfoLevel
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) IsWarn() bool {
	return l.level <= zapcore.WarnLevel
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) IsError() bool {
	return l.level <= zapcore.ErrorLevel
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) IsFatal() bool {
	return l.level <= zapcore.FatalLevel
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) IsPanic() bool {
	return l.level <= zapcore.PanicLevel
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger, l.level}
}

func (l *zapLogger) WithField(key, value string) Logger {
	return l.WithFields(Fields{key: value})
}

func (l *zapLogger) WithError(err error) Logger {
	if err == nil {
		return l
	}
	return l.WithFields(Fields{errorKey: err})
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) != 0 {
		l.sugaredLogger.With(keysAndValues...).Info(msg)
	} else {
		l.sugaredLogger.Info(msg)
	}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) != 0 {
		l.sugaredLogger.With(keysAndValues...).Debug(msg)
	} else {
		l.sugaredLogger.Debug(msg)
	}
}

func (l *zapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, errorKey, err)
	}
	if len(keysAndValues) != 0 {
		l.sugaredLogger.With(keysAndValues...).Error(msg)
	} else {
		l.sugaredLogger.Error(msg)
	}
}

func (l *zapLogger) WithName(name string) Logger {
	return &zapLogger{
		sugaredLogger: l.sugaredLogger.Named(name),
		level:         l.level,
	}
}

func (l *zapLogger) Write(b []byte) (n int, err error) {
	n = len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	l.sugaredLogger.Info(string(b))
	return n, nil
}
