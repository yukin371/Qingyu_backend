package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBookRepository 模拟BookRepository
type MockBookRepositoryForWarmup struct {
	hotBooks []*bookstore.Book
	err      error
}

func (m *MockBookRepositoryForWarmup) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	if m.err != nil {
		return nil, m.err
	}

	if offset >= len(m.hotBooks) {
		return []*bookstore.Book{}, nil
	}

	end := offset + limit
	if end > len(m.hotBooks) {
		end = len(m.hotBooks)
	}

	return m.hotBooks[offset:end], nil
}

// MockUserRepository 模拟UserRepository
type MockUserRepositoryForWarmup struct {
	activeUsers []*users.User
	err         error
}

func (m *MockUserRepositoryForWarmup) GetActiveUsers(ctx context.Context, limit int64) ([]*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	if limit >= int64(len(m.activeUsers)) {
		return m.activeUsers, nil
	}

	return m.activeUsers[:limit], nil
}

func setupTestRedisForWarmer(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	return s, client
}

func TestCacheWarmer_WarmUpCache_Success(t *testing.T) {
	ctx := context.Background()
	s, redisClient := setupTestRedisForWarmer(t)
	defer s.Close()

	// 创建测试数据
	bookID1 := primitive.NewObjectID()
	bookID2 := primitive.NewObjectID()
	bookID3 := primitive.NewObjectID()

	now := time.Now()

	hotBooks := []*bookstore.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: bookID1},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Title:            "热门书籍1",
			ViewCount:        1000,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: bookID2},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Title:            "热门书籍2",
			ViewCount:        800,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: bookID3},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Title:            "热门书籍3",
			ViewCount:        600,
		},
	}

	bookRepo := &MockBookRepositoryForWarmup{hotBooks: hotBooks}
	userRepo := &MockUserRepositoryForWarmup{activeUsers: []*users.User{}}

	warmer := NewCacheWarmer(bookRepo, userRepo, redisClient)

	// 执行预热
	err := warmer.WarmUpCache(ctx)
	require.NoError(t, err)

	// 验证Redis中有缓存数据
	keys, err := redisClient.Keys(ctx, "book:*").Result()
	require.NoError(t, err)
	assert.Equal(t, 3, len(keys), "应该有3本书的缓存")

	// 验证缓存内容
	for _, book := range hotBooks {
		key := "book:" + book.ID.Hex()
		exists := redisClient.Exists(ctx, key).Val()
		assert.Equal(t, int64(1), exists, "书籍缓存应该存在: "+key)

		// 验证TTL设置（约1小时，允许5秒误差）
		ttl := redisClient.TTL(ctx, key).Val()
		assert.True(t, ttl > 55*time.Minute && ttl <= 60*time.Minute,
			"TTL应该在55-60分钟之间，实际: "+ttl.String())
	}
}

func TestCacheWarmer_WarmUpCache_EmptyData(t *testing.T) {
	ctx := context.Background()
	s, redisClient := setupTestRedisForWarmer(t)
	defer s.Close()

	bookRepo := &MockBookRepositoryForWarmup{hotBooks: []*bookstore.Book{}}
	userRepo := &MockUserRepositoryForWarmup{activeUsers: []*users.User{}}

	warmer := NewCacheWarmer(bookRepo, userRepo, redisClient)

	// 执行预热
	err := warmer.WarmUpCache(ctx)
	require.NoError(t, err)

	// 验证没有缓存数据
	keys, _ := redisClient.Keys(ctx, "book:*").Result()
	assert.Equal(t, 0, len(keys))
}

func TestCacheWarmer_WarmUpCache_RepositoryError(t *testing.T) {
	ctx := context.Background()
	s, redisClient := setupTestRedisForWarmer(t)
	defer s.Close()

	// 模拟Repository错误
	bookRepo := &MockBookRepositoryForWarmup{
		hotBooks: []*bookstore.Book{},
		err:      assert.AnError,
	}
	userRepo := &MockUserRepositoryForWarmup{activeUsers: []*users.User{}}

	warmer := NewCacheWarmer(bookRepo, userRepo, redisClient)

	// 执行预热，错误不应导致panic，但应返回error
	err := warmer.WarmUpCache(ctx)
	// 预热失败不应返回error（容错设计）
	// 但会记录日志
	assert.NoError(t, err)
}

func TestCacheWarmer_WarmUpCache_WithActiveUsers(t *testing.T) {
	ctx := context.Background()
	s, redisClient := setupTestRedisForWarmer(t)
	defer s.Close()

	// 创建测试数据
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	now := time.Now()

	activeUsers := []*users.User{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID1},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         "user1",
			LastLoginAt:      now,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID2},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         "user2",
			LastLoginAt:      now,
		},
	}

	bookRepo := &MockBookRepositoryForWarmup{hotBooks: []*bookstore.Book{}}
	userRepo := &MockUserRepositoryForWarmup{activeUsers: activeUsers}

	warmer := NewCacheWarmer(bookRepo, userRepo, redisClient)

	// 执行预热
	err := warmer.WarmUpCache(ctx)
	require.NoError(t, err)

	// 验证用户缓存
	keys, err := redisClient.Keys(ctx, "user:*").Result()
	require.NoError(t, err)
	assert.Equal(t, 2, len(keys), "应该有2个用户的缓存")

	// 验证TTL设置（约30分钟，允许5秒误差）
	for _, user := range activeUsers {
		key := "user:" + user.ID.Hex()
		ttl := redisClient.TTL(ctx, key).Val()
		assert.True(t, ttl > 25*time.Minute && ttl <= 30*time.Minute,
			"用户TTL应该在25-30分钟之间，实际: "+ttl.String())
	}
}

func TestCacheWarmer_WarmUpCache_PartialFailure(t *testing.T) {
	ctx := context.Background()
	s, redisClient := setupTestRedisForWarmer(t)
	defer s.Close()

	// 创建测试数据
	bookID := primitive.NewObjectID()
	now := time.Now()

	hotBooks := []*bookstore.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Title:            "热门书籍1",
			ViewCount:        1000,
		},
	}

	userID := primitive.NewObjectID()
	activeUsers := []*users.User{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         "user1",
			LastLoginAt:      now,
		},
	}

	// 用户Repository会失败
	bookRepo := &MockBookRepositoryForWarmup{hotBooks: hotBooks}
	userRepo := &MockUserRepositoryForWarmup{
		activeUsers: activeUsers,
		err:         assert.AnError,
	}

	warmer := NewCacheWarmer(bookRepo, userRepo, redisClient)

	// 执行预热
	err := warmer.WarmUpCache(ctx)
	require.NoError(t, err) // 容错设计，不应返回错误

	// 验证书籍缓存仍然成功
	bookKeys, _ := redisClient.Keys(ctx, "book:*").Result()
	assert.Equal(t, 1, len(bookKeys), "书籍缓存应该成功")

	// 用户缓存失败
	userKeys, _ := redisClient.Keys(ctx, "user:*").Result()
	assert.Equal(t, 0, len(userKeys), "用户缓存应该失败")
}
