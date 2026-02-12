package redis

import (
	"context"
	"encoding/json"
	"time"

	"nusatek-backend/internal/domain"

	"github.com/redis/go-redis/v9"
)

type propertyCacheRepository struct {
	Client *redis.Client
}

func NewPropertyCacheRepository(client *redis.Client) domain.PropertyCacheRepository {
	return &propertyCacheRepository{Client: client}
}

func (r *propertyCacheRepository) Get(ctx context.Context, key string) (*domain.Property, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var p domain.Property
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *propertyCacheRepository) Set(ctx context.Context, key string, p *domain.Property, ttl time.Duration) error {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, jsonBytes, ttl).Err()
}

func (r *propertyCacheRepository) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
