package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext(t *testing.T) {
	a := assert.New(t)

	ctx := context.Background()

	out := FromContext(ctx)
	a.Nil(out)

	out = FromContextOrDefault(ctx)
	a.Same(DefaultLogger, out)

	logger := NewLogger(Configuration{LogFormat: LogFormatPlain, LogLevel: "debug"})
	lctx := NewContext(ctx, logger)

	out = FromContext(lctx)
	a.Same(logger, out)

	out = FromContextOrDefault(lctx)
	a.Same(logger, out)

}
