package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	traceutils "github.com/opentracing-contrib/go-grpc"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	grpcImg "github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image/delivery/grpc"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image/delivery/rabbitmq"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image/repository"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image/usecase"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	imageService "github.com/AleksK1NG/hotels-mocroservices/images-microservice/proto/image"
)

type Server struct {
	logger  logger.Logger
	cfg     *config.Config
	tracer  opentracing.Tracer
	pgxPool *pgxpool.Pool
	s3      *s3.S3
}

func NewServer(logger logger.Logger, cfg *config.Config, tracer opentracing.Tracer, pgxPool *pgxpool.Pool, s3 *s3.S3) *Server {
	return &Server{logger: logger, cfg: cfg, tracer: tracer, pgxPool: pgxPool, s3: s3}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	imagePublisher, err := rabbitmq.NewImagePublisher(s.cfg, s.logger)
	if err != nil {
		return errors.Wrap(err, "NewImagePublisher")
	}
	uploadedChan, err := imagePublisher.CreateExchangeAndQueue("images", "uploaded", "uploaded")
	if err != nil {
		return errors.Wrap(err, "imagePublisher.CreateExchangeAndQueue")
	}
	defer uploadedChan.Close()

	imagePGRepo := repository.NewImagePGRepository(s.pgxPool)
	imageAWSRepo := repository.NewImageAWSRepository(s.cfg, s.s3)
	imageUC := usecase.NewImageUseCase(imagePGRepo, imageAWSRepo, s.logger, imagePublisher)

	imageConsumer := rabbitmq.NewImageConsumer(s.logger, s.cfg, imageUC)
	if err := imageConsumer.Initialize(); err != nil {
		return errors.Wrap(err, "imageConsumer.Initialize")
	}
	imageConsumer.RunConsumers(ctx, cancel)
	defer imageConsumer.CloseChannels()

	l, err := net.Listen("tcp", s.cfg.GRPCServer.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	router := echo.New()
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(s.cfg.Metrics.URL); err != nil {
			s.logger.Errorf("router.Start metrics: %v", err)
			cancel()
		}
		s.logger.Infof("Metrics available on: %v", s.cfg.Metrics.URL)
	}()

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
		),
	)

	imgService := grpcImg.NewImageService(s.cfg, s.logger, imageUC)
	imageService.RegisterImageServiceServer(server, imgService)
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

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
