package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用的Cacheable实现
type TestBook struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (b *TestBook) GetID() string {
	return b.ID
}

// Mock Repository
type MockBookRepository struct {
	data map[string]*TestBook
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		data: make(map[string]*TestBook),
	}
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*TestBook, error) {
	if book, ok := m.data[id]; ok {
		return book, nil
	}
	return nil, ErrNotFound
}

func (m *MockBookRepository) Create(ctx context.Context, book *TestBook) error {
	m.data[book.ID] = book
	return nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *TestBook) error {
	if _, ok := m.data[book.ID]; !ok {
		return ErrNotFound
	}
	m.data[book.ID] = book
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.data[id]; !ok {
		return ErrNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *MockBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.data[id]
	return ok, nil
}

// 设置测试Redis
func setupTestRedis(t *testing.T) *redis.Client {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return client
}

// 测试用例
func TestCacheable(t *testing.T) {
	book := &TestBook{ID: "123"}
	assert.Implements(t, (*Cacheable)(nil), book)
	assert.Equal(t, "123", book.GetID())
}

func TestNewCachedRepository(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 1 * time.Second,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	assert.NotNil(t, repo)
	assert.True(t, repo.enabled)
}

func TestCachedRepository_GetByID_CacheMiss(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 第一次查询（缓存miss）
	result, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	assert.Equal(t, "123", result.ID)
	assert.Equal(t, "Test Book", result.Title)
}

func TestCachedRepository_GetByID_CacheHit(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 第一次查询（填充缓存）
	repo.GetByID(context.Background(), "123")
	time.Sleep(100 * time.Millisecond) // 等待异步写入

	// 第二次查询（缓存命中）
	result, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	assert.Equal(t, "123", result.ID)
}

func TestCachedRepository_GetByID_NotFound(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 查询不存在的数据
	_, err := repo.GetByID(context.Background(), "999")
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestCachedRepository_Update(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 先填充缓存
	repo.GetByID(context.Background(), "123")
	time.Sleep(100 * time.Millisecond)

	// 更新数据
	updated := &TestBook{ID: "123", Title: "Updated Book"}
	err := repo.Update(context.Background(), updated)
	require.NoError(t, err)

	// 等待双删完成
	time.Sleep(50 * time.Millisecond)

	// 验证缓存已删除（会重新查询DB）
	result, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	assert.Equal(t, "Updated Book", result.Title)
}

func TestCachedRepository_Disabled(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           false, // 禁用缓存
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 即使填充了缓存
	repo.GetByID(context.Background(), "123")
	time.Sleep(100 * time.Millisecond)

	// 缓存应该是空的
	keys, _ := redisClient.Keys(context.Background(), "book:*").Result()
	assert.Empty(t, keys)
}

func TestCachedRepository_Create(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 创建新书籍
	book := &TestBook{ID: "123", Title: "New Book", Author: "Test Author"}
	err := repo.Create(context.Background(), book)
	require.NoError(t, err)

	// 验证书籍已创建
	result, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	assert.Equal(t, "New Book", result.Title)
	assert.Equal(t, "Test Author", result.Author)
}

func TestCachedRepository_Delete(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 先填充缓存
	repo.GetByID(context.Background(), "123")
	time.Sleep(100 * time.Millisecond)

	// 删除书籍
	err := repo.Delete(context.Background(), "123")
	require.NoError(t, err)

	// 等待双删完成
	time.Sleep(50 * time.Millisecond)

	// 验证书籍已删除
	_, err = repo.GetByID(context.Background(), "123")
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestCachedRepository_Exists(t *testing.T) {
	base := NewMockBookRepository()
	book := &TestBook{ID: "123", Title: "Test Book"}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 测试存在的书籍
	exists, err := repo.Exists(context.Background(), "123")
	require.NoError(t, err)
	assert.True(t, exists)

	// 测试不存在的书籍
	exists, err = repo.Exists(context.Background(), "999")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestCachedRepository_NullCaching(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 第一次查询不存在的数据（会缓存空值）
	_, err := repo.GetByID(context.Background(), "999")
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)

	time.Sleep(100 * time.Millisecond)

	// 第二次查询（应该从空值缓存返回）
	_, err = repo.GetByID(context.Background(), "999")
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestCachedRepository_NilConfig(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)

	// 测试nil config（应该使用默认配置）
	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", nil)

	assert.NotNil(t, repo)
	assert.True(t, repo.enabled)
	assert.NotNil(t, repo.config)
}

func TestCachedRepository_Update_NotFound(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 更新不存在的书籍
	book := &TestBook{ID: "999", Title: "Non-existent"}
	err := repo.Update(context.Background(), book)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestCachedRepository_Delete_NotFound(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	config := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 10 * time.Millisecond,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
	}

	repo := NewCachedRepository[*TestBook](base, redisClient, 3600*time.Second, "book", config)

	// 删除不存在的书籍
	err := repo.Delete(context.Background(), "999")
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}
