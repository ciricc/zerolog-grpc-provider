package zerologgrpcprovider

import (
	"context"

	"github.com/rs/zerolog"
)

type loggerModifier func(context.Context, zerolog.Logger) (zerolog.Logger, error)

type loggerModifiers struct {
	modifiers []loggerModifier
}

// AddFields modifies the logger by adding custom fields.
//
// It accepts a provider function that returns a map of string to any type,
// which represents the custom fields to be added to the logger.
func (c *loggerModifiers) AddFields(
	provider func(ctx context.Context) map[string]interface{},
) *loggerModifiers {
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

func (c *loggerModifiers) getModifiedLogger(ctx context.Context, logger zerolog.Logger) (zerolog.Logger, error) {
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
