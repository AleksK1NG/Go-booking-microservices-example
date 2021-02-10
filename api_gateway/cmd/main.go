package main

import (
	"log"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/config"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/server"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/jaeger"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/redis"
)

// @title API Gateway
// @version 1.0
// @description API Gateway
// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru
// @BasePath /api/v1
func main() {
	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting API Gateway")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s",
		cfg.HttpServer.AppVersion,
		cfg.Logger.Level,
		cfg.Logger.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.HttpServer.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	redisClient := redis.NewRedisClient(cfg)
	appLogger.Infof("Redis connected: %-v", redisClient.PoolStats())

	s := server.NewServer(appLogger, cfg, redisClient, tracer)
	appLogger.Fatal(s.Run())
}
