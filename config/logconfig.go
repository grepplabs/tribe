package config

import (
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/spf13/pflag"
)

type LogConfig struct {
	flagBase
	log.Configuration
}

func NewLogConfig() *LogConfig {
	return &LogConfig{}
}

func (c *LogConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.LogLevel, "log-level", log.InfoLevel, "Log filtering One of: [fatal, error, warn, info, debug]")
		c.flagSet.StringVar(&c.LogFormat, "log-format", log.LogFormatLogfmt, "Log format to use. One of: [logfmt, json, plain]")
		c.flagSet.StringVar(&c.LogFieldNames.Time, "log-field-name-time", log.TimeKey, "Log time field name")
		c.flagSet.StringVar(&c.LogFieldNames.Message, "log-field-name-message", log.MessageKey, "Log message field name")
		c.flagSet.StringVar(&c.LogFieldNames.Caller, "log-field-name-caller", log.CallerKey, "Log caller field name")
		c.flagSet.StringVar(&c.LogFieldNames.Level, "log-field-name-level", log.LevelKey, "Log time field name")
	}
	return c.flagSet
}
