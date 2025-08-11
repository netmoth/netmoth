package redis

import (
	"context"
	"os"
	"testing"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

func TestRedisHandler_RealServerOptional(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("set TEST_REDIS_ADDR to run redis integration test")
	}
	h := NewRedisHandler(
		context.Background(),
		&redisv9.Options{Addr: addr},
	)
	if err := h.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	key := "netmoth:test:key"
	if err := h.Set(key, "v", time.Second); err != nil {
		t.Fatalf("set: %v", err)
	}
	if got, err := h.Get(key).Result(); err != nil || got != "v" {
		t.Fatalf("get: %v %q", err, got)
	}
	if n, err := h.Delete(key); err != nil || n != 1 {
		t.Fatalf("del: %v %d", err, n)
	}
}
