package cache

import (
	"context"
	"errors"
	"io"
	"time"
)

//go:generate mockgen -package=cachefakes -destination=cachefakes/fake.go . Store

type (
	// Options represents the cache store available options
	// when using Set().
	Options struct {
		// Expiration allows to specify a global expiration
		// time when setting a value.
		Expiration time.Duration
		// Tags allows specifying associated tags to the
		// current value.
		Tags []string
	}
	// Store defines methods for interacting with the
	// caching system.
	Store interface {
		// Ping pings the Redis cache to ensure its alive.
		Ping(context.Context) error
		// Get retrieves a specific item from the cache by key. Values are
		// automatically marshalled for use with Redis.
		Get(context.Context, string, any) error
		// Set stores a singular item in memory by key, value
		// and options (tags and expiration time). Values are automatically
		// marshalled for use with Redis & Memcache.
		Set(context.Context, string, any, Options)
		// Delete removes a singular item from the cache by
		// a specific key.
		Delete(context.Context, string) error
		// Invalidate removes items from the cache via the tags passed.
		Invalidate(context.Context, []string)
		// Flush removes all items from the cache.
		Flush(context.Context)
		// Closer closes the client, releasing any open resources.
		io.Closer
	}
)

// Forever defines an infinite expiration time.
const Forever = -1

// ErrNotFound is the error that's returned by the Get() method
// when no cache key was found.
var ErrNotFound = errors.New("key not found")
