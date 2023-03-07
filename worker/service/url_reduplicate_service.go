package service

import (
	"context"
	"crawler-distributed/support/redissupport"
	"github.com/redis/go-redis/v9"
)

// 使用内存去重
var visitedUrls = make(map[string]bool)

func IsDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}

// ReduplicateWithRedis 使用 Redis 去重
func ReduplicateWithRedis(ctx context.Context, client *redis.Client, url string) bool {
	value := redissupport.GetFromRedis(client, url, ctx)
	if value == "" {
		redissupport.SaveToRedis(client, url, "1", ctx)
		return false
	}
	return true
}
