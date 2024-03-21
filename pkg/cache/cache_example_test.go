package cache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/env"
)

func ExampleStore() {
	var store cache.Store
	if env.IsDevelopment() {
		store = cache.NewInMemory(time.Hour)
	} else {
		store = cache.NewRedis(&redis.Options{
			Addr:     "localhost:6379",
			Password: "passsword",
			DB:       0,
		})
	}

	// Set a value in the cache with a specific key and options
	err := store.Set(context.Background(), "key", "value", cache.Options{
		Expiration: time.Minute * 30,
		Tags:       []string{"tag1", "tag2"},
	})
	if err != nil {
		fmt.Println("Error setting value in MemCache:", err)
		return
	}

	// Retrieve the value from the cache
	var value string
	err = store.Get(context.Background(), "key", &value)
	if err != nil {
		fmt.Println("Error getting value from MemCache:", err)
		return
	}
	fmt.Println(value)
	// Output: value
}
