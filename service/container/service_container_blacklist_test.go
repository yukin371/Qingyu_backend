package container

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"Qingyu_backend/service/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthRedisClientBlacklist_NilClientReturnsNil(t *testing.T) {
	assert.Nil(t, newAuthRedisClientBlacklist(nil))
}

func TestAuthRedisClientBlacklist_AddCheckRemove(t *testing.T) {
	ctx := context.Background()
	client := auth.NewInMemoryTokenBlacklist()
	defer func() {
		_ = client.Close()
	}()

	blacklist := newAuthRedisClientBlacklist(client)
	require.NotNil(t, blacklist)

	token := "token-123"
	require.NoError(t, blacklist.Add(ctx, token, time.Minute))

	revoked, err := blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.True(t, revoked)

	require.NoError(t, blacklist.Remove(ctx, token))

	revoked, err = blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.False(t, revoked)
}

func TestAuthRedisClientBlacklist_UsesHashedKey(t *testing.T) {
	ctx := context.Background()
	client := auth.NewInMemoryTokenBlacklist()
	defer func() {
		_ = client.Close()
	}()

	blacklist := newAuthRedisClientBlacklist(client)
	require.NotNil(t, blacklist)

	token := "token-456"
	require.NoError(t, blacklist.Add(ctx, token, time.Minute))

	h := sha256.New()
	_, _ = h.Write([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	key := fmt.Sprintf("token:blacklist:%s", tokenHash[:32])

	exists, err := client.Exists(ctx, key)
	require.NoError(t, err)
	assert.EqualValues(t, 1, exists)
}
