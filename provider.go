package zerologgrpcprovider

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ZerologGrpcProvider interface {
	// UnaryInterceptor returns interceptor compatible with grpc api for provide zerolog logger
	UnaryInterceptor() grpc.UnaryServerInterceptor
	// StreamInterceptor returns interceptor compatible with grpc api for provide zerolog logger
	StreamInterceptor() grpc.StreamServerInterceptor
	// WithModifiers returns object for customizing zerolog logger
	WithModifiers() *loggerModifiers
}

type zerologGrpcProvider struct {
	options   *Options
	modifiers *loggerModifiers
}

func New(opts ...Option) (ZerologGrpcProvider, error) {
	defaultLogger := zerolog.New(os.Stdout)

	options := Options{
		logRequests:                  true,
		useRequestId:                 true,
		provideRequestFieldsToLogger: true,
		logErrors:                    true,
		requestLogger:                &defaultLogger,
	}

	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, fmt.Errorf("set option error: %w", err)
		}
	}

	return &zerologGrpcProvider{
		options: &options,
		modifiers: &loggerModifiers{
			modifiers: []loggerModifier{},
		},
	}, nil
}

func (z *zerologGrpcProvider) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error,
	) {
		initialLogger, err := z.getLogger(ctx)
		if err != nil {
			return nil, err
		}

		loggerCtx := initialLogger.With().Bool(grpcUnaryInterceptorFieldName, true)
		if z.options.useRequestId {
			loggerCtx = z.loggerWithRequestId(loggerCtx)
		}

		logger := loggerCtx.Logger()
		loggerWithRequestFields := loggerCtx

		if z.options.provideRequestFieldsToLogger {
			requestProtoMessage, ok := req.(proto.Message)
			if !ok {
				return nil, ErrFailedToCastProtoMessage
			}

			requestMap, err := protobufToMap(requestProtoMessage)
			if err != nil {
				return nil, fmt.Errorf("failed to cast protobuf message into map: %w", err)
			}

			if z.options.requestValueModifier != nil {
				err = z.modifyRequestValues(requestMap)
				if err != nil {
					return nil, err
				}
			}

			loggerWithRequestFields = loggerCtx.Fields(map[string]interface{}{
				grpcRequestFieldName: requestMap,
				grpcMethodFieldName:  info.FullMethod,
				grpcServerFieldName:  info.Server,
			})
		}

		if z.options.logRequests {
			logger := loggerWithRequestFields.Logger()
			(&logger).Debug().Msg("new unary request")
		}

		res, err := handler(contextWithLogger(ctx, &logger), req)
		if err != nil && z.options.logErrors {
			(&logger).Err(err).Msg("unary request error")

			return res, err
		}

		return res, err
	}
}

func (z *zerologGrpcProvider) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		initialLogger, err := z.getLogger(ctx)
		if err != nil {
			return err
		}

		loggerCtx := initialLogger.With().Bool(grpcUnaryInterceptorFieldName, false)
		if z.options.useRequestId {
			loggerCtx = z.loggerWithRequestId(loggerCtx)
		}

		if z.options.provideRequestFieldsToLogger {
			loggerCtx = loggerCtx.Fields(map[string]interface{}{
				grpcStreamInfoFieldName: info,
			})
		}

		logger := loggerCtx.Logger()

		if z.options.logRequests {
			(&logger).Debug().Msg("new stream request")
		}

		wrapper := serverStreamWrapper{
			ServerStream: ss,
			ctx:          contextWithLogger(ctx, &logger),
		}

		err = handler(srv, &wrapper)
		if err != nil && z.options.logErrors {
			(&logger).Err(err).Msg("stream request error")

			return err
		}

		return err
	}
}

func (z *zerologGrpcProvider) modifyRequestValues(requestMap map[string]any) error {
	for k, v := range requestMap {
		mapValue, isMap := v.(map[string]any)
		if isMap {
			err := z.modifyRequestValues(mapValue)
			if err != nil {
				return err
			}
		}

		newValue, err := z.options.requestValueModifier(k, v)
		if err != nil {
			return fmt.Errorf("failed to modify request value (key=%s): %w", k, err)
		}

		requestMap[k] = newValue
	}

	return nil
}

func (z *zerologGrpcProvider) loggerWithRequestId(ctx zerolog.Context) zerolog.Context {
	return ctx.Str(grpcRequestIdFieldName, uuid.NewString())
}

func (z *zerologGrpcProvider) getLogger(ctx context.Context) (*zerolog.Logger, error) {
	newLogger, err := z.modifiers.getModifiedLogger(ctx, *z.options.requestLogger)
	if err != nil {
		return z.options.requestLogger, err
	}

	return &newLogger, nil
}

func (z *zerologGrpcProvider) WithModifiers() *loggerModifiers {
	return z.modifiers
}
