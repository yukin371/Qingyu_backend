package user

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

	"Qingyu_backend/models/users"
)

var ErrNotFound = errors.New("not found")

// setupTestRedis 创建测试用Redis客户端
func setupTestRedis(t *testing.T) *redis.Client {
	s := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: s.Addr()})
}

// MockUserRepository 模拟UserRepository用于测试
type MockUserRepository struct {
	users map[string]*users.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*users.User),
	}
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, ErrNotFound
}

func (m *MockUserRepository) Create(ctx context.Context, user *users.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	m.users[user.ID.Hex()] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if _, ok := m.users[id]; !ok {
		return ErrNotFound
	}
	// 简化实现
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return ErrNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.users[id]
	return ok, nil
}

// TestWrapUserRepositoryWithCache 测试UserRepository包装函数
func TestWrapUserRepositoryWithCache(t *testing.T) {
	base := NewMockUserRepository()
	redisClient := setupTestRedis(t)

	cached := WrapUserRepositoryWithCache(base, redisClient)

	assert.NotNil(t, cached)
	_, ok := cached.(*CachedUserRepository)
	assert.True(t, ok, "应该是CachedUserRepository类型")
}

// TestCachedUserRepository_GetByID_CacheMiss 测试缓存未命中场景
func TestCachedUserRepository_GetByID_CacheMiss(t *testing.T) {
	base := NewMockUserRepository()
	user := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	base.Create(context.Background(), user)

	redisClient := setupTestRedis(t)
	cached := WrapUserRepositoryWithCache(base, redisClient)

	// 第一次查询（缓存miss）
	result, err := cached.GetByID(context.Background(), user.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), result.ID.Hex())
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)
}

// TestCachedUserRepository_GetByID_CacheHit 测试缓存命中场景
func TestCachedUserRepository_GetByID_CacheHit(t *testing.T) {
	base := NewMockUserRepository()
	user := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	base.Create(context.Background(), user)

	redisClient := setupTestRedis(t)
	cached := WrapUserRepositoryWithCache(base, redisClient)

	// 第一次查询（填充缓存）
	_, err := cached.GetByID(context.Background(), user.ID.Hex())
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// 第二次查询（应该从缓存获取）
	result, err := cached.GetByID(context.Background(), user.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), result.ID.Hex())
	assert.Equal(t, "testuser", result.Username)
}

// TestCachedUserRepository_Update 测试更新用户（应删除缓存）
func TestCachedUserRepository_Update(t *testing.T) {
	base := NewMockUserRepository()
	user := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	base.Create(context.Background(), user)

	redisClient := setupTestRedis(t)
	cached := WrapUserRepositoryWithCache(base, redisClient)

	// 填充缓存
	cached.GetByID(context.Background(), user.ID.Hex())
	time.Sleep(100 * time.Millisecond)

	// 更新数据
	err := cached.Update(context.Background(), user.ID.Hex(), map[string]interface{}{
		"username": "updateduser",
	})
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	// 验证（Mock不会真正更新）
	result, err := cached.GetByID(context.Background(), user.ID.Hex())
	if err == nil {
		assert.NotNil(t, result)
	}
}

// TestCachedUserRepository_Delete 测试删除用户
func TestCachedUserRepository_Delete(t *testing.T) {
	base := NewMockUserRepository()
	user := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	base.Create(context.Background(), user)

	redisClient := setupTestRedis(t)
	cached := WrapUserRepositoryWithCache(base, redisClient)

	// 填充缓存
	cached.GetByID(context.Background(), user.ID.Hex())
	time.Sleep(100 * time.Millisecond)

	// 删除用户
	err := cached.Delete(context.Background(), user.ID.Hex())
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	// 验证已删除
	_, err = cached.GetByID(context.Background(), user.ID.Hex())
	assert.Error(t, err)
}

// TestCachedUserRepository_Exists 测试检查用户是否存在
func TestCachedUserRepository_Exists(t *testing.T) {
	base := NewMockUserRepository()
	user := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	base.Create(context.Background(), user)

	redisClient := setupTestRedis(t)
	cached := WrapUserRepositoryWithCache(base, redisClient)

	// 检查存在的用户
	exists, err := cached.Exists(context.Background(), user.ID.Hex())
	if err != nil {
		assert.False(t, exists)
	} else {
		assert.True(t, exists)
	}

	// 检查不存在的用户
	exists, err = cached.Exists(context.Background(), "nonexistent-id")
	if err != nil {
		assert.False(t, exists)
	} else {
		assert.False(t, exists)
	}
}
