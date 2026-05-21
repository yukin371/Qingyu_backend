package auth

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRedisClient struct {
	values map[string]string
	ttls   map[string]time.Duration
}

func newStubRedisClient() *stubRedisClient {
	return &stubRedisClient{
		values: make(map[string]string),
		ttls:   make(map[string]time.Duration),
	}
}

func (s *stubRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	s.values[key] = value.(string)
	s.ttls[key] = expiration
	return nil
}

func (s *stubRedisClient) Get(ctx context.Context, key string) (string, error) {
	return s.values[key], nil
}

func (s *stubRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	var count int64
	for _, key := range keys {
		if _, ok := s.values[key]; ok {
			count++
		}
	}
	return count, nil
}

func testJWTConfig() *config.JWTConfigEnhanced {
	return &config.JWTConfigEnhanced{
		SecretKey:       "unit-test-secret",
		Issuer:          "qingyu-backend-test",
		Expiration:      time.Hour,
		RefreshDuration: 24 * time.Hour,
	}
}

func TestJWTService_RevokeTokenMarksTokenAsRevoked(t *testing.T) {
	ctx := context.Background()
	redisClient := newStubRedisClient()
	service := NewJWTService(testJWTConfig(), redisClient)

	token, err := service.GenerateToken(ctx, "user-1", []string{"reader"})
	require.NoError(t, err)

	err = service.RevokeToken(ctx, token)
	require.NoError(t, err)

	revoked, err := service.IsTokenRevoked(ctx, token)
	require.NoError(t, err)
	assert.True(t, revoked)
	assert.Len(t, redisClient.values, 1)

	for _, ttl := range redisClient.ttls {
		assert.Greater(t, ttl, time.Duration(0))
	}
}

func TestJWTService_ValidateTokenRejectsRevokedToken(t *testing.T) {
	ctx := context.Background()
	redisClient := newStubRedisClient()
	service := NewJWTService(testJWTConfig(), redisClient)

	token, err := service.GenerateToken(ctx, "user-2", []string{"writer"})
	require.NoError(t, err)
	require.NoError(t, service.RevokeToken(ctx, token))

	claims, err := service.ValidateToken(ctx, token)
	require.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token已被吊销")
}

func TestInMemoryTokenBlacklistBasicFlow(t *testing.T) {
	ctx := context.Background()
	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() {
		_ = blacklist.Close()
	})

	require.NoError(t, blacklist.Set(ctx, "token:blacklist:test", "revoked", time.Minute))

	value, err := blacklist.Get(ctx, "token:blacklist:test")
	require.NoError(t, err)
	assert.Equal(t, "revoked", value)

	exists, err := blacklist.Exists(ctx, "token:blacklist:test")
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	require.NoError(t, blacklist.Delete(ctx, "token:blacklist:test"))

	exists, err = blacklist.Exists(ctx, "token:blacklist:test")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}
