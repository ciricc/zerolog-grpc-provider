package zerologgrpcprovider

import (
	"context"

	"google.golang.org/grpc"
)

type serverStreamWrapper struct {
	grpc.ServerStream

	//nolint:containedctx
	ctx context.Context
}

func (w *serverStreamWrapper) Context() context.Context { return w.ctx }
