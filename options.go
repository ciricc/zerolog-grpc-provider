package zerologgrpcprovider

import "github.com/rs/zerolog"

type Options struct {
	// Whether need to log requests or not
	logRequests bool
	// Whether need to log errors or not
	logErrors bool
	// Whether need push request fields to log or not
	provideRequestFieldsToLogger bool
	// Whether need generate request id in the log or not
	useRequestId bool
	// Zerolog logger
	requestLogger *zerolog.Logger
}

type Option func(opts *Options) error

// When loggin is true, the logger will print messages by self into zerolog output.
// Default value is true
func WithLogRequests(logging bool) Option {
	return func(opts *Options) error {
		opts.logRequests = logging

		return nil
	}
}

// When logging errors is true, the logger will print errors after request completion.
// Default value is true
func WithLogErrors(logging bool) Option {
	return func(opts *Options) error {
		opts.logErrors = logging

		return nil
	}
}

// When provideFields is enabled, provider will add into zerolog context some request fields
// Like grpcMethod, grpcServer information etc.
func WithProvideRequestFieldsToLogger(provideFields bool) Option {
	return func(opts *Options) error {
		opts.provideRequestFieldsToLogger = provideFields

		return nil
	}
}

// When use request id is enabled, interceptor will generate a request identifier in the UUID format
// and will add this into zerolog context fields list
func WithUseRequestId(useRequestId bool) Option {
	return func(opts *Options) error {
		opts.useRequestId = useRequestId

		return nil
	}
}

// WithLogger changes the default zerolog logger
func WithLogger(logger *zerolog.Logger) Option {
	return func(opts *Options) error {
		opts.requestLogger = logger

		return nil
	}
}
