package zerologgrpcprovider

import (
	"context"

	"github.com/rs/zerolog"
)

// Context key for *zerolog.Logger structure
type loggerCtxKey struct{}

func ContextWithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}
