package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	traceutils "github.com/opentracing-contrib/go-grpc"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/config"
	crfRepository "github.com/AleksK1NG/hotels-mocroservices/sessions/internal/csrf/repository"
	csrfUseCase "github.com/AleksK1NG/hotels-mocroservices/sessions/internal/csrf/usecase"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/interceptors"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session/delivery"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session/repository"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session/usecase"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/pkg/logger"
	sessionService "github.com/AleksK1NG/hotels-mocroservices/sessions/proto"
)

// Server
type Server struct {
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
	tracer    opentracing.Tracer
}

// NewServer
func NewSessionsServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client, tracer opentracing.Tracer) *Server {
	return &Server{logger: logger, cfg: cfg, redisConn: redisConn, tracer: tracer}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	im := interceptors.NewInterceptorManager(s.logger, s.cfg)
	sessionRedisRepo := repository.NewSessionRedisRepo(s.redisConn, s.cfg.GRPCServer.SessionPrefix, time.Duration(s.cfg.GRPCServer.SessionExpire)*time.Minute)
	sessionUC := usecase.NewSessionUseCase(sessionRedisRepo)
	csrfRepository := crfRepository.NewCsrfRepository(s.redisConn, s.cfg.GRPCServer.CSRFPrefix, time.Duration(s.cfg.GRPCServer.CsrfExpire)*time.Minute)
	csrfUC := csrfUseCase.NewCsrfUseCase(csrfRepository)

	router := echo.New()
	router.GET("", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(":7070"); err != nil {
			s.logger.Errorf("router.Start metrics: %v", err)
			cancel()
		}
		s.logger.Infof("Metrics available on: %v", s.cfg.Metrics.URL)
	}()

	l, err := net.Listen("tcp", s.cfg.GRPCServer.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.GRPCServer.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.GRPCServer.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.GRPCServer.MaxConnectionAge * time.Minute,
		Time:              s.cfg.GRPCServer.Timeout * time.Minute,
	}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			traceutils.OpenTracingServerInterceptor(s.tracer),
			im.Logger,
		),
	)

	sessGRPCService := delivery.NewSessionsService(s.logger, sessionUC, csrfUC)
	sessionService.RegisterAuthorizationServiceServer(server, sessGRPCService)
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

	if err := router.Shutdown(ctx); err != nil {
		s.logger.Errorf("Metrics router.Shutdown: %v", err)
	}

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
