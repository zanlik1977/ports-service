package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zanlik1977/ports-service/pkg/redis"
	"github.com/zanlik1977/ports-service/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdatePort(t *testing.T) {
	rdb, err := redis.NewRedisClient()
	assert.NoError(t, err)

	portService := service.NewPortService(rdb)

	// Set up the test server
	server := httptest.NewServer(http.HandlerFunc(portService.HandlePorts))
	defer server.Close()

	port := map[string]interface{}{
		"id":          "AEAJM",
		"name":        "Ajman",
		"city":        "Ajman",
		"country":     "United Arab Emirates",
		"alias":       []string{},
		"regions":     []string{},
		"coordinates": []float64{55.5136433, 25.4052165},
		"province":    "Ajman",
		"timezone":    "Asia/Dubai",
		"unlocs":      []string{"AEAJM"},
		"code":        "52000",
	}

	portJSON, err := json.Marshal(port)
	assert.NoError(t, err)

	resp, err := http.Post(server.URL+"/ports", "application/json", bytes.NewBuffer(portJSON))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the port was upserted
	data, err := rdb.Get(context.Background(), "AEAJM").Result()
	assert.NoError(t, err)

	var retrievedPort map[string]interface{}
	err = json.Unmarshal([]byte(data), &retrievedPort)
	assert.NoError(t, err)
	assert.Equal(t, port["name"], retrievedPort["name"])
	assert.Equal(t, port["city"], retrievedPort["city"])
	assert.Equal(t, port["country"], retrievedPort["country"])
}

func TestInvalidPortData(t *testing.T) {
	rdb, err := redis.NewRedisClient()
	assert.NoError(t, err)

	portService := service.NewPortService(rdb)

	// Set up the test server
	server := httptest.NewServer(http.HandlerFunc(portService.HandlePorts))
	defer server.Close()

	// Test empty port data
	resp, err := http.Post(server.URL+"/ports", "application/json", bytes.NewBuffer([]byte("{}")))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test missing ID
	port := map[string]interface{}{
		"name": "Invalid Port",
	}
	portJSON, err := json.Marshal(port)
	assert.NoError(t, err)

	resp, err = http.Post(server.URL+"/ports", "application/json", bytes.NewBuffer(portJSON))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
