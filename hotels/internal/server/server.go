package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/config"
	hotelsGrpc "github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels/delivery/grpc"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels/repository"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels/usecase"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
	hotelsService "github.com/AleksK1NG/hotels-mocroservices/hotels/proto/hotels"
)

var (
	zapLogger  *zap.Logger
	customFunc grpc_zap.CodeToLevel
)

// Server
type Server struct {
	echo      *echo.Echo
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
	pgxPool   *pgxpool.Pool
	tracer    opentracing.Tracer
}

// NewServer
func NewServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client, pgxPool *pgxpool.Pool, tracer opentracing.Tracer) *Server {
	return &Server{logger: logger, cfg: cfg, redisConn: redisConn, pgxPool: pgxPool, echo: echo.New(), tracer: tracer}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()
	hotelsPGRepo := repository.NewHotelsPGRepository(s.pgxPool)
	hotelsUC := usecase.NewHotelsUC(hotelsPGRepo, s.logger)

	l, err := net.Listen("tcp", s.cfg.GRPCServer.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	// opts := []grpc_zap.Option{
	// 	grpc_zap.WithLevels(customFunc),
	// 	grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
	// 		return zap.Int64("grpc.time_ns", duration.Nanoseconds())
	// 	}),
	// }
	// grpc_zap.ReplaceGrpcLoggerV2(zapLogger)

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.GRPCServer.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.GRPCServer.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.GRPCServer.MaxConnectionAge * time.Minute,
		Time:              s.cfg.GRPCServer.Timeout * time.Minute,
	}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
			// im.Logger,
		),
		// grpc_middleware.WithUnaryServerChain(
		// 	grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		// 	grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
		// 	grpc_ctxtags.UnaryServerInterceptor(),
		// ),
	)

	hotelsGRPCService := hotelsGrpc.NewHotelsService(hotelsUC, s.logger, validate)
	hotelsService.RegisterHotelsServiceServer(server, hotelsGRPCService)
	grpc_prometheus.Register(server)

	go func() {
		s.logger.Infof("GRPC Server is listening on port: %v", s.cfg.GRPCServer.Port)
		s.logger.Fatal(server.Serve(l))
	}()

	if s.cfg.GRPCServer.Mode != "Production" {
		reflection.Register(server)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	s.logger.Info("Server Exited Properly")

	if err := s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
