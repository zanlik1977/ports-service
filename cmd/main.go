package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"ports-service/pkg/domain"
	"ports-service/pkg/redis"
	"ports-service/pkg/service"
	"syscall"
)

func main() {
	rdb, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	portService := service.NewPortService(rdb)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	file, err := os.Open("ports.json")
	if err != nil {
		log.Fatalf("Failed to open ports.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var port domain.Port
		if err := decoder.Decode(&port); err != nil {
			break
		}
		if err := portService.UpsertPort(ctx, &port); err != nil {
			log.Printf("Failed to upsert port: %v", err)
		}
	}

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
