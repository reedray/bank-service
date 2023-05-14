package storage

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/reedray/bank-service/config/converter"
	"strconv"
	"time"
)
import "context"

type RedisRepository struct {
	client    *redis.Client
	expiresAt time.Duration
}

func NewRedis(cfg *converter.Config) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("can`t connect to a database %w", err)
	}

	return &RedisRepository{
		client:    client,
		expiresAt: cfg.Redis.ExpiresAt,
	}, nil
}

// GetExchangeRates takes a pair of currency codes as a string
// for example EUR:USD and returns a ratio
func (r *RedisRepository) GetExchangeRates(ctx context.Context, currenciesCodeKey string) (float64, error) {
	result, err := r.client.Get(ctx, currenciesCodeKey).Result()
	if err != nil || err == redis.Nil {
		return 0, err
	}
	res, err := strconv.ParseFloat(result, 64)
	if err != nil {
		return 0, fmt.Errorf("can`t parse string to float %w", err)
	}
	return res, nil

}

// SetExchangeRates takes a pair of currency codes as a string
// for example EUR:USD and sets a ratio for them
func (r *RedisRepository) SetExchangeRates(ctx context.Context, currencyCode string, ratio float64) error {
	fmt.Println("IN REPO SET")
	err := r.client.Set(ctx, currencyCode, ratio, r.expiresAt).Err()
	if err != nil {
		return fmt.Errorf("can`t set a value to a database %w", err)
	}

	return nil
}
