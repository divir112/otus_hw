package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	logger Logger
}

func NewInterceptor(logger Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

func (m *Interceptor) loggingMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "can't get data from ctx")
	}
	remoteAddress := p.Addr.String()
	var ip string
	if remoteAddress != "" {
		ip = strings.Split(remoteAddress, ":")[0]
	}
	date := time.Now()
	formatTime := date.Format("25/Feb/2020:19:11:24 +0600")
	method := info.FullMethod
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	duration := time.Since(date).Microseconds()
	m.logger.Info(fmt.Sprintf("%s [%s] %s %s", ip, formatTime, method, duration))

	return resp, nil
}
