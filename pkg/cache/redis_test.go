package cache

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/pkg/cache/internal"
)

// Setup is a helper to obtain a mock cache store for testing.
func Setup(t *testing.T, mf func(m *internal.MockRedisStore)) *Redis {
	ctrl := gomock.NewController(t)
	m := internal.NewMockRedisStore(ctrl)
	if mf != nil {
		mf(m)
	}
	return &Redis{
		client: m,
		mtx:    &sync.Mutex{},
	}
}

func TestCache_Flush(t *testing.T) {
	c := Setup(t, func(m *internal.MockRedisStore) {
		m.EXPECT().FlushAll(context.TODO()).AnyTimes()
	})
	c.Flush(context.TODO())
}

// testCacheStruct represents a struct for working with
// JSON values within the cache store.
type testCacheStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

var (
	// key is the test key used for Redis testing.
	key = "key"
	// tag is the test tag used for Redis testing.
	tag = "tag"
	// value is the test value to match against testing for
	// get and set test methods, to see if it's marshalling
	// properly.
	value = testCacheStruct{
		Name:  "name",
		Value: 1,
	}
	// options are the default testing set options.

	ctx = context.TODO()
)

func TestNew(t *testing.T) {
	t.Parallel()

	got := NewRedis(&redis.Options{})
	assert.NotNil(t, got.client)
	assert.NotNil(t, got.mtx)
}

func TestPing(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		mock    func(m *internal.MockRedisStore)
		wantErr bool
	}{
		"Success": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Ping(ctx).
					Return(redis.NewStatusCmd(ctx, nil))
			},
			false,
		},
		"Error": {
			func(m *internal.MockRedisStore) {
				cmd := redis.NewStatusCmd(ctx, nil)
				cmd.SetErr(errors.New("ping error"))
				m.EXPECT().
					Ping(ctx).
					Return(cmd)
			},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			err := c.Ping(ctx)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		mock    func(m *internal.MockRedisStore)
		wantErr bool
	}{
		"Success": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Close().
					Return(nil)
			},
			false,
		},
		"Error": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Close().
					Return(errors.New("close error"))
			},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			err := c.Close()
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		mock    func(m *internal.MockRedisStore)
		wantErr bool
		want    any
	}{
		"Success": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Get(ctx, key).
					Return(redis.NewStringResult(`{"name": "name", "value": 1}`, nil))
			},
			false,
			testCacheStruct{
				Name:  "name",
				Value: 1,
			},
		},
		"Redis Error": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Get(ctx, key).
					Return(redis.NewStringResult("", errors.New("redis error")))
			},
			true,
			testCacheStruct{},
		},
		"Decode Error": {
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Get(ctx, key).
					Return(redis.NewStringResult("wrong", nil))
			},
			true,
			testCacheStruct{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			got := testCacheStruct{}
			err := c.Get(ctx, key, &got)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		value   any
		mock    func(m *internal.MockRedisStore)
		wantErr bool
	}{
		"Success": {
			key,
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Set(ctx, key, gomock.Any(), gomock.Any()).
					Return(redis.NewStatusCmd(ctx, nil))

				m.EXPECT().
					SAdd(ctx, tag, key).
					Return(redis.NewIntCmd(ctx, ""))

				m.EXPECT().
					Expire(ctx, tag, 720*time.Hour).
					Return(redis.NewBoolCmd(ctx, true))
			},
			false,
		},
		"Redis Error": {
			"key",
			func(m *internal.MockRedisStore) {
				cmd := redis.NewStatusCmd(ctx, nil)
				cmd.SetErr(errors.New("redis error"))

				m.EXPECT().
					Set(ctx, "key", gomock.Any(), gomock.Any()).
					Return(cmd)
			},
			true,
		},
		"Encode Error": {
			make(chan string),
			func(m *internal.MockRedisStore) {
			},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			err := c.Set(ctx, key, test.value, Options{
				Expiration: -1,
				Tags:       []string{tag},
			})
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestCache_Delete(t *testing.T) {
	tt := map[string]struct {
		value   any
		mock    func(m *internal.MockRedisStore)
		wantErr bool
	}{
		"Success": {
			value,
			func(m *internal.MockRedisStore) {
				m.EXPECT().
					Del(ctx, key).
					Return(redis.NewIntCmd(ctx, nil))
			},
			false,
		},
		"Redis Error": {
			value,
			func(m *internal.MockRedisStore) {
				cmd := redis.NewIntCmd(ctx, nil)
				cmd.SetErr(errors.New("delete error"))

				m.EXPECT().
					Del(ctx, key).
					Return(cmd)
			},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			err := c.Delete(ctx, key)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestCache_Invalidate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []string
		mock  func(m *internal.MockRedisStore)
	}{
		"Success": {
			[]string{tag},
			func(m *internal.MockRedisStore) {
				cmd := redis.NewStringSliceCmd(ctx)
				cmd.SetVal([]string{key})

				m.EXPECT().
					SMembers(ctx, gomock.Any()).
					Return(cmd)

				m.EXPECT().
					Del(ctx, gomock.Any()).
					Return(redis.NewIntCmd(ctx, nil)).
					AnyTimes()
			},
		},
		"Nil Tags": {
			nil,
			nil,
		},
		"SMembers Error": {
			[]string{tag},
			func(m *internal.MockRedisStore) {
				cmd := redis.NewStringSliceCmd(ctx)
				cmd.SetErr(errors.New("err"))

				m.EXPECT().
					SMembers(ctx, gomock.Any()).
					Return(cmd)
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			c := Setup(t, test.mock)
			c.Invalidate(ctx, test.input)
		})
	}
}
