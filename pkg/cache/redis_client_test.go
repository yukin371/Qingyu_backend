package cache

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	args := m.Called(ctx, keys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRedisClient) MSet(ctx context.Context, pairs ...interface{}) error {
	args := m.Called(ctx, pairs)
	return args.Error(0)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) GetClient() interface{} {
	args := m.Called()
	return args.Get(0)
}

// ========== 单元测试 ==========

func TestDefaultRedisConfig(t *testing.T) {
	cfg := config.DefaultRedisConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 6379, cfg.Port)
	assert.Equal(t, 0, cfg.DB)
	assert.Equal(t, 10, cfg.PoolSize)
	assert.Equal(t, 5, cfg.MinIdleConns)
}

func TestWrapRedisError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: nil,
		},
		{
			name:     "timeout error",
			input:    assert.AnError,
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapRedisError(tt.input)
			if tt.expected == nil {
				assert.NoError(t, result)
			}
		})
	}
}

// TestMockRedisClient 测试Mock客户端
func TestMockRedisClient(t *testing.T) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()

	// 测试Get
	mockClient.On("Get", ctx, "test_key").Return("test_value", nil)
	val, err := mockClient.Get(ctx, "test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)

	// 测试Set
	mockClient.On("Set", ctx, "test_key", "test_value", time.Minute).Return(nil)
	err = mockClient.Set(ctx, "test_key", "test_value", time.Minute)
	assert.NoError(t, err)

	// 测试Delete
	mockClient.On("Delete", ctx, []string{"test_key"}).Return(nil)
	err = mockClient.Delete(ctx, "test_key")
	assert.NoError(t, err)

	// 测试Exists
	mockClient.On("Exists", ctx, []string{"test_key"}).Return(int64(1), nil)
	count, err := mockClient.Exists(ctx, "test_key")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	mockClient.AssertExpectations(t)
}

// TestMockRedisClient_HashOperations 测试Hash操作
func TestMockRedisClient_HashOperations(t *testing.T) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()

	// HSet
	mockClient.On("HSet", ctx, "hash_key", mock.Anything).Return(nil)
	err := mockClient.HSet(ctx, "hash_key", "field1", "value1")
	assert.NoError(t, err)

	// HGet
	mockClient.On("HGet", ctx, "hash_key", "field1").Return("value1", nil)
	val, err := mockClient.HGet(ctx, "hash_key", "field1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// HGetAll
	expected := map[string]string{"field1": "value1", "field2": "value2"}
	mockClient.On("HGetAll", ctx, "hash_key").Return(expected, nil)
	result, err := mockClient.HGetAll(ctx, "hash_key")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockClient.AssertExpectations(t)
}

// TestMockRedisClient_SetOperations 测试Set操作
func TestMockRedisClient_SetOperations(t *testing.T) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()

	// SAdd
	mockClient.On("SAdd", ctx, "set_key", mock.Anything).Return(nil)
	err := mockClient.SAdd(ctx, "set_key", "member1", "member2")
	assert.NoError(t, err)

	// SMembers
	expected := []string{"member1", "member2"}
	mockClient.On("SMembers", ctx, "set_key").Return(expected, nil)
	members, err := mockClient.SMembers(ctx, "set_key")
	assert.NoError(t, err)
	assert.Equal(t, expected, members)

	mockClient.AssertExpectations(t)
}

// TestMockRedisClient_AtomicOperations 测试原子操作
func TestMockRedisClient_AtomicOperations(t *testing.T) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()

	// Incr
	mockClient.On("Incr", ctx, "counter").Return(int64(1), nil)
	val, err := mockClient.Incr(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	// IncrBy
	mockClient.On("IncrBy", ctx, "counter", int64(5)).Return(int64(6), nil)
	val, err = mockClient.IncrBy(ctx, "counter", 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), val)

	// Decr
	mockClient.On("Decr", ctx, "counter").Return(int64(5), nil)
	val, err = mockClient.Decr(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	mockClient.AssertExpectations(t)
}

// TestMockRedisClient_Ping 测试Ping
func TestMockRedisClient_Ping(t *testing.T) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()

	mockClient.On("Ping", ctx).Return(nil)
	err := mockClient.Ping(ctx)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

// Benchmark测试
func BenchmarkMockRedisClient_Get(b *testing.B) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()
	mockClient.On("Get", ctx, "test_key").Return("test_value", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockClient.Get(ctx, "test_key")
	}
}

func BenchmarkMockRedisClient_Set(b *testing.B) {
	mockClient := new(MockRedisClient)
	ctx := context.Background()
	mockClient.On("Set", ctx, "test_key", "test_value", time.Minute).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockClient.Set(ctx, "test_key", "test_value", time.Minute)
	}
}
