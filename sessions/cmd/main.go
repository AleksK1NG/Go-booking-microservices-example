package main

import (
	"log"
	"os"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/config"
)

func main() {
	log.Println("Starting sessions server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	log.Printf("%-v", cfg)

}
