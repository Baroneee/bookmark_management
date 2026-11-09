package main

import (
	"context"
	"time"

	"github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
)

// main exercises the configured Redis clients.
func main() {
	ctx := context.Background()

	rdb, _ := redis.NewClient("")

	rdb.Set(ctx, "123", "google.com", time.Hour)

	cacheRedisDB, _ := redis.NewClient("CACHE")
	cacheRedisDB.Set(ctx, "hjk", "huhuhuhuhu", time.Hour)
}
