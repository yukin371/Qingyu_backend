package cache

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/shared/cache/strategies"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisCacheService_WithStrategyManager_AppliesDefaultTTL(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})
	manager := strategies.NewStrategyManager(3 * time.Second)
	service := NewRedisCacheServiceWithStrategyManager(redisClient, manager)
	ctx := context.Background()

	err := service.Set(ctx, "book:hot:1", "v", 0)
	require.NoError(t, err)

	ttl, err := service.TTL(ctx, "book:hot:1")
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, 3*time.Second)
}
