package interceptors

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/AleksK1NG/hotels-mocroservices/comments/config"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
)

// InterceptorManager
type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
	tracer opentracing.Tracer
}

// InterceptorManager constructor
func NewInterceptorManager(logger logger.Logger, cfg *config.Config, tracer opentracing.Tracer) *InterceptorManager {
	return &InterceptorManager{logger: logger, cfg: cfg, tracer: tracer}
}

// Logger Interceptor
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}

// GetInterceptor
func (im *InterceptorManager) GetInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		im.logger.Infof("call=%v req=%#v reply=%#v time=%v err=%v",
			method, req, reply, time.Since(start), err)
		return err
	}
}

// GetTracer
func (im *InterceptorManager) GetTracer() opentracing.Tracer {
	return im.tracer
}
