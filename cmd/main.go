package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zanlik1977/ports-service/pkg/redis"
	"github.com/zanlik1977/ports-service/pkg/service"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rdb, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	portService := service.NewPortService(rdb)

	http.HandleFunc("/ports", portService.HandlePorts)

	server := &http.Server{Addr: ":8080"}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	log.Println("Server started on :8080")

	// Load ports from file
	go loadPortsIncrementally()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")

	if err := server.Close(); err != nil {
		log.Fatalf("Server Close: %v", err)
	}
}

func loadPortsIncrementally() {
	file, err := os.Open("/app/ports.json")
	if err != nil {
		log.Fatalf("Failed to open ports.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read the opening brace
	if _, err := decoder.Token(); err != nil {
		log.Fatalf("Failed to read JSON: %v", err)
	}

	for decoder.More() {
		// Read the port ID (key)
		token, err := decoder.Token()
		if err != nil {
			log.Printf("Failed to read port ID: %v", err)
			break
		}
		id, ok := token.(string)
		if !ok {
			log.Printf("Expected string for port ID, got %T", token)
			continue
		}

		// Read the port details (value)
		var port map[string]interface{}
		if err := decoder.Decode(&port); err != nil {
			log.Printf("Failed to decode port %s: %v", id, err)
			continue
		}

		// Ensure ID is part of the port data
		port["id"] = id
		portJSON, err := json.Marshal(port)
		if err != nil {
			log.Printf("Failed to marshal port %s: %v", id, err)
			continue
		}

		// send POST request
		resp, err := http.Post("http://localhost:8080/ports", "application/json", bytes.NewBuffer(portJSON))
		if err != nil {
			log.Printf("Failed to send port %s: %v", id, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Failed to process port %s, statusCode %d, response: %s", id, resp.StatusCode, string(body))
		} else {
			fmt.Printf("Port %s processed successfully\n", id)
		}
	}

	// Read the closing brace
	if _, err := decoder.Token(); err != nil {
		log.Fatalf("Failed to read JSON: %v", err)
	}
}
