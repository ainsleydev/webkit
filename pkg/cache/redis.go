package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"

	"github.com/ainsleydev/webkit/pkg/cache/internal"
)

// Redis defines the methods for interacting with the
// cache layer.
type Redis struct {
	client internal.RedisStore
	mtx    *sync.Mutex
}

// NewRedis returns a new instance of the Redis cache store.
func NewRedis(opts *redis.Options) *Redis {
	return &Redis{
		client: redis.NewClient(opts),
		mtx:    &sync.Mutex{},
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *Redis) Get(ctx context.Context, key string, v any) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return ErrNotFound
	}

	err = json.Unmarshal([]byte(result), v)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, options Options) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	buf, err := json.Marshal(value)
	if err != nil {
		slog.Error("marshalling cache value: " + err.Error())
		return
	}

	err = r.client.Set(ctx, key, buf, options.Expiration).Err()
	if err != nil {
		slog.Error("setting cache value: " + err.Error())
		return
	}

	if len(options.Tags) > 0 {
		r.setTags(ctx, key, options.Tags)
	}
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Invalidate(ctx context.Context, tags []string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if len(tags) == 0 {
		return
	}

	for _, tag := range tags {
		cacheKeys, err := r.client.SMembers(ctx, tag).Result()
		if err != nil {
			continue
		}

		for _, cacheKey := range cacheKeys {
			r.client.Del(ctx, cacheKey)
		}

		r.client.Del(ctx, tag)
	}
}

func (r *Redis) Flush(ctx context.Context) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.client.FlushAll(ctx)
}

func (r *Redis) Close() error {
	return r.client.Close()
}

// setTags sets SMembers in the redis store for caching.
func (r *Redis) setTags(ctx context.Context, key any, tags []string) {
	for _, tag := range tags {
		r.client.SAdd(ctx, tag, key.(string))
		r.client.Expire(ctx, tag, 720*time.Hour)
	}
}
