package redis

import (
	"context"
	"encoding/json"
	"time"

	"SavingBooks/config"
	"SavingBooks/internal/domain"
	saving_regulation "SavingBooks/internal/saving-regulation"
	"SavingBooks/internal/services/redis/redis_key"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
	savingRegulationRepo saving_regulation.SavingRegulationRepository
}

func (c *Cache) SetValue(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, 5 * time.Minute).Err()
}

func (c *Cache) GetValue(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}
func(c *Cache) GetLatestSavingRegulation(ctx context.Context) (*domain.SavingRegulation, error) {
	latestRegulation := &domain.SavingRegulation{}
	err := c.GetValue(ctx, redis_key.LatestRegulation, latestRegulation)
	if err != nil {
		latestRegulation, err = c.savingRegulationRepo.GetLatestSavingRegulation(ctx)
		if err != nil {
			return nil, err
		}
		_ = c.SetValue(ctx, redis_key.LatestRegulation, latestRegulation)
	}
	return latestRegulation, nil
}

func (c *Cache) RemoveValue(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}


func NewCacheService(c *config.Configuration, savingRegulationRepo saving_regulation.SavingRegulationRepository) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis,
		Password: "",
		DB: 0,
	})
	return &Cache{redis: rdb, savingRegulationRepo: savingRegulationRepo}
}

