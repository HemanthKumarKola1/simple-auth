package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RevokedTokensCache struct {
	rdb redis.Client
}

func NewRevokedTokensCache(rdbClient *redis.Client) *RevokedTokensCache {
	return &RevokedTokensCache{rdb: *rdbClient}
}

func (t *RevokedTokensCache) SetWithTTL(key string, value string, ttl time.Duration) error {
	t.rdb.Set("hellow", "wassup", 0)
	fmt.Println(key)
	return t.rdb.SetNX(key, value, ttl).Err()
}

func (t *RevokedTokensCache) Get(key string) (string, error) {
	fmt.Println(key)
	return t.rdb.Get(key).Result()
}
