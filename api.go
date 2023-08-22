package zerologgrpcprovider

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// GetLogger returns provided logger from the context
func GetLogger(ctx context.Context) (*zerolog.Logger, error) {
	loggerInterface := ctx.Value(loggerCtxKey{})

	if loggerInterface == nil {
		return nil, fmt.Errorf("no logger provided to the context")
	}

	loggerValue, ok := loggerInterface.(*zerolog.Logger)

	if !ok {
		return nil, fmt.Errorf("failed to cast logger interface value to *zerolog.Logger type")
	}

	return loggerValue, nil
}

// MustGetLogger returns the logger from the context
// May panic if logger is not available
func MustGetLogger(ctx context.Context) *zerolog.Logger {
	logger, err := GetLogger(ctx)
	if err != nil {
		panic(err)
	}

	return logger
}
