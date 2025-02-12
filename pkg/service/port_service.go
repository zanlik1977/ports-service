package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type PortService struct {
	rdb *redis.Client
}

func NewPortService(rdb *redis.Client) *PortService {
	return &PortService{rdb: rdb}
}

func (s *PortService) HandlePorts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createOrUpdatePort(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *PortService) createOrUpdatePort(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Decode into a map
	var port map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&port); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check for the port ID
	id, ok := port["id"].(string)
	if !ok || id == "" {
		http.Error(w, "Missing port ID", http.StatusBadRequest)
		return
	}

	// Marshal the map back into JSON
	data, err := json.Marshal(port)
	if err != nil {
		http.Error(w, "Failed to process port", http.StatusInternalServerError)
		return
	}

	// Store in Redis using the ID as the key
	if err := s.rdb.Set(ctx, id, data, 0).Err(); err != nil {
		http.Error(w, "Failed to save port", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
