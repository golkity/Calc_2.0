package tokens

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
}

func NewStore(addr string) *Store {
	return &Store{redis.NewClient(&redis.Options{Addr: addr})}
}

func (s *Store) Save(ctx context.Context, token string, ttl time.Duration) error {
	return s.rdb.Set(ctx, token, "1", ttl).Err()
}

func (s *Store) Exists(ctx context.Context, token string) (bool, error) {
	res, err := s.rdb.Exists(ctx, token).Result()
	return res == 1, err
}

func (s *Store) Delete(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, token).Err()
}
