package log

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZapLogger(t *testing.T) {
	a := assert.New(t)

	type LevelEnabled struct {
		debug bool
		info  bool
		warn  bool
		error bool
		panic bool
		fatal bool
	}

	tt := []struct {
		name         string
		logger       Logger
		levelEnabled LevelEnabled
	}{
		{"default logger", NewDefaultLogger(), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"default logger with field", NewDefaultLogger().WithField("tag", "value"), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"default logger with error", NewDefaultLogger().WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"default logger with nil error", NewDefaultLogger().WithError(nil), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"plain debug", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "debug"}), LevelEnabled{
			debug: true, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"plain info", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "info"}), LevelEnabled{
			debug: false, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"plain warn", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "warn"}), LevelEnabled{
			debug: false, info: false, warn: true, error: true, panic: true, fatal: true}},
		{"plain error", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "error"}), LevelEnabled{
			debug: false, info: false, warn: false, error: true, panic: true, fatal: true}},
		{"plain panic", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "panic"}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: true, fatal: true}},
		{"plain fatal", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "fatal"}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: false, fatal: true}},
		{"plain with error", NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "info"}).WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"json debug", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: DebugLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: true, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"json info", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: InfoLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: false, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"json warn", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: WarnLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: false, info: false, warn: true, error: true, panic: true, fatal: true}},
		{"json error", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: ErrorLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: false, info: false, warn: false, error: true, panic: true, fatal: true}},
		{"json panic", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: PanicLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: true, fatal: true}},
		{"json fatal", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: FatalLevel}).WithFields(Fields{"tag": "value"}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: false, fatal: true}},
		{"json with error", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: "info"}).WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"json changed field names", NewLogger(Configuration{LogFormat: LogFormatJson, LogLevel: "info", LogFieldNames: LogFieldNames{
			Time: "time", Message: "message", Level: "lvl", Caller: "call",
		}}).WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"logfmt debug", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: DebugLevel}), LevelEnabled{
			debug: true, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"logfmt info", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: InfoLevel}), LevelEnabled{
			debug: false, info: true, warn: true, error: true, panic: true, fatal: true}},
		{"logfmt warn", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: WarnLevel}), LevelEnabled{
			debug: false, info: false, warn: true, error: true, panic: true, fatal: true}},
		{"logfmt error", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: ErrorLevel}), LevelEnabled{
			debug: false, info: false, warn: false, error: true, panic: true, fatal: true}},
		{"logfmt panic", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: PanicLevel}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: true, fatal: true}},
		{"logfmt fatal", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: FatalLevel}), LevelEnabled{
			debug: false, info: false, warn: false, error: false, panic: false, fatal: true}},
		{"logfmt with error", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: "info"}).WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
		{"logfmt changed field names", NewLogger(Configuration{LogFormat: LogFormatLogfmt, LogLevel: "info", LogFieldNames: LogFieldNames{
			Time: "time", Message: "message", Level: "lvl", Caller: "call",
		}}).WithError(errors.New("my error")), LevelEnabled{
			debug: false, info: true, warn: true, error: true, fatal: true, panic: true}},
	}
	{
		for _, tc := range tt {
			tc.logger.Debugf("DebugLevel log '%s'", tc.name)
			tc.logger.Infof("InfoLevel log '%s'", tc.name)
			tc.logger.Warnf("WarnLevel log '%s'", tc.name)
			tc.logger.Errorf("WithError log '%s'", tc.name)

			_, err := tc.logger.Write([]byte(fmt.Sprintf("Write interface'%s'", tc.name)))
			a.Nil(err)

			a.Equal(tc.levelEnabled.debug, tc.logger.IsDebug())
			a.Equal(tc.levelEnabled.info, tc.logger.IsInfo())
			a.Equal(tc.levelEnabled.warn, tc.logger.IsWarn())
			a.Equal(tc.levelEnabled.error, tc.logger.IsError())
			a.Equal(tc.levelEnabled.fatal, tc.logger.IsFatal())
			a.Equal(tc.levelEnabled.panic, tc.logger.IsPanic())
		}
	}
}

func TestDefaultLogger(t *testing.T) {
	const name = "default"

	Debugf("DebugLevel log '%s'", name)
	Infof("InfoLevel log '%s'", name)
	Printf("Printf log '%s'", name)
	Warnf("WarnLevel log '%s'", name)
	Errorf("ErrorLevel log '%s'", name)

	WithFields(Fields{"tag": "value"}).Infof("InfoLevel log '%s'", name)
	WithField("tag", "value").Infof("InfoLevel log '%s'", name)
	WithError(errors.New("my error")).Infof("InfoLevel log '%s'", name)

	WithName("new-named").Infof("New named logger")

	Debug("Debug message", "key1", "value1")
	Info("Info message", "key1", "value1")
	Error(nil, "Error message", "key1", "value1")
	Error(errors.New("test error"), "Error message", "key1", "value1")
}
