package zerologgrpcprovider

import (
	"context"

	"github.com/rs/zerolog"
)

type loggerModification func(context.Context, zerolog.Logger) (zerolog.Logger, error)

type loggerModifier struct {
	modifiers []loggerModification
}

// CustomFields modifies the logger by adding custom fields.
//
// It accepts a provider function that returns a map of string to any type,
// which represents the custom fields to be added to the logger.
func (c *loggerModifier) CustomFields(provider func(ctx context.Context) map[string]any) *loggerModifier {
	modifier := func(ctx context.Context, logger zerolog.Logger) (zerolog.Logger, error) {
		fields := provider(ctx)
		if fields == nil {
			return logger, nil
		}

		loggerCtx := logger.With().Fields(fields)
		newLogger := loggerCtx.Logger()

		return newLogger, nil
	}

	c.modifiers = append(c.modifiers, modifier)

	return c
}

func (c *loggerModifier) getModifiedLogger(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	newLogger := logger

	for _, modifier := range c.modifiers {
		modifiedLogger, err := modifier(ctx, logger)
		if err != nil {
			continue
		}

		newLogger = modifiedLogger
	}

	return newLogger
}
