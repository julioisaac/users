package cache

import (
	"context"
	"time"
)

type Cache interface {
	RemoveAll(ctx context.Context, keys ...string) (int64, error)
	Add(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}
