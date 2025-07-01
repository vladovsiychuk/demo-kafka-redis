package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func PartnerKey(partner, date string) string {
	return fmt.Sprintf("%s:%s", partner, date)
}

type State struct {
	Clicks      int     `json:"clicks"`
	Cost        float64 `json:"cost"`
	Date        string  `json:"date"`
	Impressions int     `json:"impressions"`
	Installs    int     `json:"installs"`
}

type Store interface {
	Get(ctx context.Context, partner, date string) (*State, error)
	Set(ctx context.Context, partner, date string, state *State) error
	Close() error
}

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		ttl: 30 * 24 * time.Hour,
	}
}

func (r *RedisStore) Get(ctx context.Context, partner, date string) (*State, error) {
	key := PartnerKey(partner, date)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var state State
	if err := json.Unmarshal([]byte(val), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *RedisStore) Set(ctx context.Context, partner, date string, state *State) error {
	key := PartnerKey(partner, date)
	bytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, r.ttl).Err()
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}
