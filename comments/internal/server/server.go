package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/AleksK1NG/hotels-mocroservices/comments/config"
	commGRPC "github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment/delivery/grpc"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment/repository"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment/usecase"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
	commentsService "github.com/AleksK1NG/hotels-mocroservices/comments/proto"
)

// Server
type server struct {
	e       *echo.Echo
	logger  logger.Logger
	cfg     *config.Config
	pgxPool *pgxpool.Pool
	tracer  opentracing.Tracer
}

// NewServer
func NewServer(logger logger.Logger, cfg *config.Config, pgxPool *pgxpool.Pool, tracer opentracing.Tracer) *server {
	return &server{e: echo.New(), logger: logger, cfg: cfg, pgxPool: pgxPool, tracer: tracer}
}

// Run
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()

	commPGRepo := repository.NewCommPGRepo(s.pgxPool)
	commUC := usecase.NewCommUseCase(commPGRepo, s.logger)
	commService := commGRPC.NewCommentsService(commUC, s.logger, s.cfg, validate)

	l, err := net.Listen("tcp", s.cfg.GRPCServer.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		s.logger.Infof("Metrics server is running on port: %s", s.cfg.Metrics.Port)
		if err := router.Start(s.cfg.Metrics.Port); err != nil {
			s.logger.Error(err)
			cancel()
		}
	}()

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
		),
	)

	commentsService.RegisterCommentsServiceServer(server, commService)
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

	if err := s.e.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
