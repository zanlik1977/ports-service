package service

import (
	"context"
	"encoding/json"
	"ports-service/pkg/domain"

	"github.com/go-redis/redis/v8"
)

type PortService struct {
	rdb *redis.Client
}

func NewPortService(rdb *redis.Client) *PortService {
	return &PortService{rdb: rdb}
}

func (s *PortService) UpsertPort(ctx context.Context, port *domain.Port) error {
	data, err := json.Marshal(port)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, port.ID, data, 0).Err()
}

