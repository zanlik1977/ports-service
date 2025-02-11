package tests

import (
	"context"
	"encoding/json"
	"ports-service/pkg/redis"
	"ports-service/pkg/service"
	"ports-service/pkg/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsertPort(t *testing.T) {
	rdb, err := redis.NewRedisClient()
	assert.NoError(t, err)

	portService := service.NewPortService(rdb)

	port := &domain.Port{ID: "AEAJM", Name: "Ajman"}
	err = portService.UpsertPort(context.Background(), port)
	assert.NoError(t, err)

	// Retrieve the port to verify it was upserted
	var retrievedPort domain.Port
	data, err := rdb.Get(context.Background(), port.ID).Result()
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(data), &retrievedPort)
	assert.NoError(t, err)
	assert.Equal(t, port.Name, retrievedPort.Name)
}
