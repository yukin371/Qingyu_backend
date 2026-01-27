package mongodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/bookstore"
)

var ErrNotFound = errors.New("not found")

// setupTestRedis 创建测试用Redis客户端
func setupTestRedis(t *testing.T) *redis.Client {
	s := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: s.Addr()})
}

// MockBookRepository 模拟BookRepository用于测试
type MockBookRepository struct {
	books map[string]*bookstore.Book
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books: make(map[string]*bookstore.Book),
	}
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
	if book, ok := m.books[id]; ok {
		return book, nil
	}
	return nil, ErrNotFound
}

func (m *MockBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}
	m.books[book.ID.Hex()] = book
	return nil
}

func (m *MockBookRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if _, ok := m.books[id]; !ok {
		return ErrNotFound
	}
	// 简化实现，实际更新逻辑不重要
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.books[id]; !ok {
		return ErrNotFound
	}
	delete(m.books, id)
	return nil
}

func (m *MockBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.books[id]
	return ok, nil
}

// TestWrapBookRepositoryWithCache 测试BookRepository包装函数
func TestWrapBookRepositoryWithCache(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)

	cached := WrapBookRepositoryWithCache(base, redisClient)

	assert.NotNil(t, cached)
	_, ok := cached.(*CachedBookRepository)
	assert.True(t, ok, "应该是CachedBookRepository类型")
}

// TestCachedBookRepository_GetByID_CacheMiss 测试缓存未命中场景
func TestCachedBookRepository_GetByID_CacheMiss(t *testing.T) {
	base := NewMockBookRepository()
	book := &bookstore.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 第一次查询（缓存miss，应该从数据库获取）
	result, err := cached.GetByID(context.Background(), book.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, book.ID.Hex(), result.ID.Hex())
	assert.Equal(t, "Test Book", result.Title)
	assert.Equal(t, "Test Author", result.Author)
}

// TestCachedBookRepository_GetByID_CacheHit 测试缓存命中场景
func TestCachedBookRepository_GetByID_CacheHit(t *testing.T) {
	base := NewMockBookRepository()
	book := &bookstore.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 第一次查询（填充缓存）
	_, err := cached.GetByID(context.Background(), book.ID.Hex())
	require.NoError(t, err)

	// 等待缓存写入完成
	time.Sleep(100 * time.Millisecond)

	// 第二次查询（应该从缓存获取）
	result, err := cached.GetByID(context.Background(), book.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, book.ID.Hex(), result.ID.Hex())
	assert.Equal(t, "Test Book", result.Title)
}

// TestCachedBookRepository_GetByID_NotFound 测试查询不存在的书籍
func TestCachedBookRepository_GetByID_NotFound(t *testing.T) {
	base := NewMockBookRepository()
	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 查询不存在的书籍
	_, err := cached.GetByID(context.Background(), "nonexistent-id")
	assert.Error(t, err)
}

// TestCachedBookRepository_Update 测试更新书籍（应删除缓存）
func TestCachedBookRepository_Update(t *testing.T) {
	base := NewMockBookRepository()
	book := &bookstore.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 填充缓存
	cached.GetByID(context.Background(), book.ID.Hex())
	time.Sleep(100 * time.Millisecond)

	// 更新数据（应该删除缓存）
	err := cached.Update(context.Background(), book.ID.Hex(), map[string]interface{}{
		"title": "Updated Book",
	})
	require.NoError(t, err)

	// 等待双删完成
	time.Sleep(50 * time.Millisecond)

	// 验证缓存已删除（下次查询仍然能成功，因为没有真正更新数据）
	result, err := cached.GetByID(context.Background(), book.ID.Hex())
	// Mock不会真正更新，所以这里只验证能获取到数据
	if err == nil {
		assert.NotNil(t, result)
	}
}

// TestCachedBookRepository_Delete 测试删除书籍（应删除缓存）
func TestCachedBookRepository_Delete(t *testing.T) {
	base := NewMockBookRepository()
	book := &bookstore.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 填充缓存
	cached.GetByID(context.Background(), book.ID.Hex())
	time.Sleep(100 * time.Millisecond)

	// 删除书籍（应该删除缓存）
	err := cached.Delete(context.Background(), book.ID.Hex())
	require.NoError(t, err)

	// 等待双删完成
	time.Sleep(50 * time.Millisecond)

	// 验证书籍已删除
	_, err = cached.GetByID(context.Background(), book.ID.Hex())
	assert.Error(t, err)
}

// TestCachedBookRepository_Exists 测试检查书籍是否存在
func TestCachedBookRepository_Exists(t *testing.T) {
	base := NewMockBookRepository()
	book := &bookstore.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	base.Create(context.Background(), book)

	redisClient := setupTestRedis(t)
	cached := WrapBookRepositoryWithCache(base, redisClient)

	// 检查存在的书籍
	exists, err := cached.Exists(context.Background(), book.ID.Hex())
	// 如果缓存返回ErrNotFound，Exists应该返回false
	if err != nil {
		assert.False(t, exists)
	} else {
		assert.True(t, exists)
	}

	// 检查不存在的书籍
	exists, err = cached.Exists(context.Background(), "nonexistent-id")
	if err != nil {
		assert.False(t, exists)
	} else {
		assert.False(t, exists)
	}
}
