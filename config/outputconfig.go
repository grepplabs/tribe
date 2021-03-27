package config

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/yamlpc"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"os"
)

var outputProducers = map[string]runtime.Producer{
	"json": runtime.JSONProducer(),
	"yaml": yamlpc.YAMLProducer(),
}

type OutputConfig struct {
	flagBase
	Format string
}

func NewOutputConfig() *OutputConfig {
	return &OutputConfig{}
}

func (c *OutputConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVarP(&c.Format, "output", "o", "json", "Output format. One of: json|yaml")
	}
	return c.flagSet
}

func (c *OutputConfig) Validate() error {
	switch c.Format {
	case "json", "yaml":
		return nil
	default:
		return errors.Errorf("Unsupported output format: %s", c.Format)
	}
}

// MustGetProducer returns the producer or calls os.Exit(1) if the output is unknown
func (c *OutputConfig) MustGetProducer() runtime.Producer {
	producer, ok := outputProducers[c.Format]
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to match a printer suitable for the output format: %s\n", c.Format)
		os.Exit(1)
	}
	return producer
}
