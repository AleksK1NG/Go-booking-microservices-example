package main

import (
	"log"
	"os"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/config"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
)

func main() {
	log.Println("Starting hotels microservice")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s",
		cfg.GRPCServer.AppVersion,
		cfg.Logger.Level,
		cfg.GRPCServer.Mode,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.GRPCServer.AppVersion)

	log.Printf("CFG: %-v", cfg)
}
