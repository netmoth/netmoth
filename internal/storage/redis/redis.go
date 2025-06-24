package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHandler is an interface that defines the methods to interact with Redis.
type RedisHandler interface {
	Ping() error
	Set(key string, value any, expiration time.Duration) error
	Get(key string) *redis.StringCmd
	Delete(key string) (int64, error)

	Client() *redis.Client
}

type redisHandler struct {
	ctx    context.Context
	client *redis.Client
}

// NewRedisHandler is ...
func NewRedisHandler(ctx context.Context, opts *redis.Options) RedisHandler {
	return &redisHandler{
		ctx:    ctx,
		client: redis.NewClient(opts),
	}
}

// Ping is a method of the redisHandler struct that sends a PING command to Redis using the client's Ping method.
// It returns an error if the Ping method fails, otherwise it returns nil.
func (h *redisHandler) Ping() error {
	return h.client.Ping(h.ctx).Err()
}

// Set is a method of redisHandler struct that sets the value of the given key to the provided value with an expiration time.
// It takes three parameters:
// - key: a string representing the key in Redis database
// - value: an any representing the value to be set for the given key
// - expiration: a time.Duration representing the time duration until the key expires
// The method returns an error if there's any issue encountered while setting the value for the given key.
func (h *redisHandler) Set(key string, value any, expiration time.Duration) error {
	return h.client.Set(h.ctx, key, value, expiration).Err()
}

// The following function is a method of the redisHandler struct.
// It takes in a string key and returns a pointer to a StringCmd object.
func (h *redisHandler) Get(key string) *redis.StringCmd {
	return h.client.Get(h.ctx, key)
}

// The function Delete is a method of the redisHandler struct which deletes the given key from Redis database.
// It takes a string as input parameter and returns the number of keys that were removed and an error (if any).
func (h *redisHandler) Delete(key string) (int64, error) {
	return h.client.Del(h.ctx, key).Result()
}

// The following function is named Client and belongs to the redisHandler struct.
// It returns a pointer to a redis.Client object.
// It receives a pointer to a redisHandler object as a receiver.
func (h *redisHandler) Client() *redis.Client {
	return h.client
}
