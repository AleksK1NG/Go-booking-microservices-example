package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	userGRPC "github.com/AleksK1NG/hotels-mocroservices/user/internal/user/delivery/grpc"
	userHandlers "github.com/AleksK1NG/hotels-mocroservices/user/internal/user/delivery/http"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user/repository"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user/usecase"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
	userGRPCService "github.com/AleksK1NG/hotels-mocroservices/user/proto/user"
)

const (
	certFile       = "ssl/server.crt"
	keyFile        = "ssl/server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server
type Server struct {
	echo      *echo.Echo
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
	pgxPool   *pgxpool.Pool
}

// NewServer
func NewServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client, pgxPool *pgxpool.Pool) *Server {
	return &Server{logger: logger, cfg: cfg, redisConn: redisConn, pgxPool: pgxPool, echo: echo.New()}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()
	v1 := s.echo.Group("/api/v1")
	usersGroup := v1.Group("/users")

	userPGRepository := repository.NewUserPGRepository(s.pgxPool)
	userUseCase := usecase.NewUserUseCase(userPGRepository)
	uh := userHandlers.NewUserHandlers(usersGroup, userUseCase, s.logger, validate)
	uh.MapUserRoutes()

	v1.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})
	v1.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	s.MapRoutes()

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.HttpServer.Port)
		s.echo.Server.ReadTimeout = time.Second * s.cfg.HttpServer.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.HttpServer.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		if err := s.echo.StartTLS(s.cfg.HttpServer.Port, certFile, keyFile); err != nil {
			s.logger.Fatalf("Error starting TLS Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HttpServer.PprofPort)
		if err := http.ListenAndServe(s.cfg.HttpServer.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
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
		// grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	userService := userGRPC.NewUserService(userUseCase, s.logger)
	userGRPCService.RegisterUserServiceServer(server, userService)
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
