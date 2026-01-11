package auth

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionService_ValidateSession 测试会话验证
func TestSessionService_ValidateSession(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*SessionService, string) (*Session, error)
		sessionID string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "验证有效会话",
			setup: func(s *SessionService, userID string) (*Session, error) {
				return s.CreateSession(context.Background(), userID)
			},
			wantErr: false,
		},
		{
			name: "验证不存在的会话",
			setup: func(s *SessionService, userID string) (*Session, error) {
				return nil, nil
			},
			sessionID: "nonexistent_session_id",
			wantErr:   true,
			errMsg:    "会话不存在或已过期",
		},
		{
			name: "验证过期会话",
			setup: func(s *SessionService, userID string) (*Session, error) {
				session, _ := s.CreateSession(context.Background(), userID)
				// 手动设置为过期
				session.ExpiresAt = time.Now().Add(-1 * time.Hour)
				s.cache.Set(context.Background(), getSessionKey(session.ID), session, 0)
				return session, nil
			},
			wantErr: true,
			errMsg:  "会话不存在或已过期",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockCacheClient()
			service := NewSessionService(cache)
			ctx := context.Background()
			userID := "user_123"

			var sessionID string
			if tt.setup != nil {
				session, _ := tt.setup(service, userID)
				if session != nil {
					sessionID = session.ID
				}
			}
			if tt.sessionID != "" {
				sessionID = tt.sessionID
			}

			// Act
			valid, err := service.ValidateSession(ctx, sessionID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, valid)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, valid)
			}
		})
	}
}

// TestSessionService_UpdateSessionModel 测试更新会话模型
func TestSessionService_UpdateSessionModel(t *testing.T) {
	// Arrange
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()
	userID := "user_123"

	// 创建会话
	session, err := service.CreateSession(ctx, userID)
	require.NoError(t, err)

	// 更新会话元数据
	session.Metadata = map[string]interface{}{
		"last_action":    "read_chapter",
		"chapter_id":     "chapter_123",
		"reading_time":   3600, // 秒
		"device_type":    "mobile",
		"app_version":    "1.0.0",
	}

	// Act
	err = service.UpdateSessionModel(ctx, session)

	// Assert
	assert.NoError(t, err)

	// 验证更新
	updated, err := service.GetSession(ctx, session.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "read_chapter", updated.Metadata["last_action"])
	assert.Equal(t, "chapter_123", updated.Metadata["chapter_id"])
}

// TestSessionService_CleanupExpiredSessions 测试清理过期会话
func TestSessionService_CleanupExpiredSessions(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(*SessionService) int
		expectedCleanups  int
	}{
		{
			name: "清理多个过期会话",
			setup: func(s *SessionService) int {
				ctx := context.Background()
				userID := "user_123"

				// 创建3个会话，其中2个过期
				s1, _ := s.CreateSession(ctx, userID)
				s1.ExpiresAt = time.Now().Add(-1 * time.Hour)
				s.cache.Set(ctx, getSessionKey(s1.ID), s1, 0)

				s2, _ := s.CreateSession(ctx, userID)
				s2.ExpiresAt = time.Now().Add(-2 * time.Hour)
				s.cache.Set(ctx, getSessionKey(s2.ID), s2, 0)

				s3, _ := s.CreateSession(ctx, userID)
				// s3未过期

				return 3
			},
			expectedCleanups: 2,
		},
		{
			name: "没有过期会话",
			setup: func(s *SessionService) int {
				ctx := context.Background()
				userID := "user_123"

				// 创建2个未过期的会话
				_, _ = s.CreateSession(ctx, userID)
				_, _ = s.CreateSession(ctx, userID)

				return 2
			},
			expectedCleanups: 0,
		},
		{
			name: "清理所有会话",
			setup: func(s *SessionService) int {
				ctx := context.Background()
				userID := "user_123"

				// 创建5个会话，全部过期
				for i := 0; i < 5; i++ {
					s, _ := s.CreateSession(ctx, userID)
					s.ExpiresAt = time.Now().Add(-time.Duration(i+1) * time.Hour)
					s.cache.Set(ctx, getSessionKey(s.ID), s, 0)
				}

				return 5
			},
			expectedCleanups: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockCacheClient()
			service := NewSessionService(cache)
			ctx := context.Background()

			tt.setup(service)

			// Act
			count, err := service.CleanupExpiredSessions(ctx)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCleanups, count)
		})
	}
}

// TestSessionService_StopCleanupTask 测试停止清理任务
func TestSessionService_StopCleanupTask(t *testing.T) {
	// Arrange
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	// 启动清理任务
	service.StartCleanupTask(ctx, 1*time.Minute)
	assert.True(t, service.IsCleanupRunning())

	// Act
	err := service.StopCleanupTask()

	// Assert
	assert.NoError(t, err)
	assert.False(t, service.IsCleanupRunning())
}

// TestSessionService_ConcurrentValidation 测试并发验证会话
func TestSessionService_ConcurrentValidation(t *testing.T) {
	// Arrange
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()
	userID := "user_concurrent"

	// 创建10个会话
	sessions := make([]*Session, 10)
	for i := 0; i < 10; i++ {
		session, err := service.CreateSession(ctx, userID)
		require.NoError(t, err)
		sessions[i] = session
	}

	// Act - 并发验证
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for _, session := range sessions {
		wg.Add(1)
		go func(s *Session) {
			defer wg.Done()
			valid, err := service.ValidateSession(ctx, s.ID)
			if err != nil || !valid {
				errors <- err
			}
		}(session)
	}

	wg.Wait()
	close(errors)

	// Assert
	for err := range errors {
		assert.NoError(t, err, "并发验证不应该出错")
	}
}

// BenchmarkSessionService_CreateSession 基准测试 - 创建会话
func BenchmarkSessionService_CreateSession(b *testing.B) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateSession(ctx, "user_123")
	}
}

// BenchmarkSessionService_ValidateSession 基准测试 - 验证会话
func BenchmarkSessionService_ValidateSession(b *testing.B) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	// 预创建会话
	session, _ := service.CreateSession(ctx, "user_123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateSession(ctx, session.ID)
	}
}

// BenchmarkSessionService_GetUserSessions 基准测试 - 获取用户会话列表
func BenchmarkSessionService_GetUserSessions(b *testing.B) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()
	userID := "user_123"

	// 预创建5个会话
	for i := 0; i < 5; i++ {
		service.CreateSession(ctx, userID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetUserSessions(ctx, userID)
	}
}

// TestSessionService_DeviceLimits 表格驱动测试 - 设备限制
func TestSessionService_DeviceLimits(t *testing.T) {
	tests := []struct {
		name           string
		maxDevices     int
		createCount    int
		expectedCount  int
		shouldEvict    bool
	}{
		{
			name:          "未超过设备限制",
			maxDevices:    5,
			createCount:   3,
			expectedCount: 3,
			shouldEvict:   false,
		},
		{
			name:          "正好达到设备限制",
			maxDevices:    5,
			createCount:   5,
			expectedCount: 5,
			shouldEvict:   false,
		},
		{
			name:          "超过设备限制，触发FIFO驱逐",
			maxDevices:    3,
			createCount:   5,
			expectedCount: 3,
			shouldEvict:   true,
		},
		{
			name:          "单设备限制",
			maxDevices:    1,
			createCount:   3,
			expectedCount: 1,
			shouldEvict:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockCacheClient()
			service := NewSessionService(cache)
			service.maxSessionsPerUser = tt.maxDevices
			ctx := context.Background()
			userID := "user_limits"

			// Act
			sessions := make([]*Session, tt.createCount)
			for i := 0; i < tt.createCount; i++ {
				session, err := service.CreateSession(ctx, userID)
				require.NoError(t, err)
				sessions[i] = session
			}

			// Assert
			allSessions, err := service.GetUserSessions(ctx, userID)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(allSessions), "应该保持%d个会话", tt.expectedCount)

			if tt.shouldEvict {
				// 验证最早的会话被驱逐
				oldestSession := sessions[0]
				_, err := service.GetSession(ctx, oldestSession.ID)
				assert.Error(t, err, "最旧的会话应该被驱逐")
			}
		})
	}
}

// TestSessionService_ErrorScenarios 测试错误场景
func TestSessionService_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*SessionService, context.Context) error
		wantErr   bool
		errContains string
	}{
		{
			name: "创建会话时缓存失败",
			setup: func(s *SessionService, ctx context.Context) error {
				// 设置缓存为失败状态
				s.cache = &FailingCacheClient{}
				_, err := s.CreateSession(ctx, "user_123")
				return err
			},
			wantErr: true,
		},
		{
			name: "更新不存在的会话",
			setup: func(s *SessionService, ctx context.Context) error {
				session := &Session{
					ID:        "nonexistent",
					UserID:    "user_123",
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}
				return s.UpdateSessionModel(ctx, session)
			},
			wantErr: true,
		},
		{
			name: "销毁不存在的会话",
			setup: func(s *SessionService, ctx context.Context) error {
				return s.DestroySession(ctx, "nonexistent_session")
			},
			wantErr: false, // 不应该报错，幂等操作
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockCacheClient()
			service := NewSessionService(cache)
			ctx := context.Background()

			// Act
			err := tt.setup(service, ctx)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else if tt.name == "销毁不存在的会话" {
				// 幂等操作，不应该报错
				assert.NoError(t, err)
			}
		})
	}
}

// FailingCacheClient 模拟失败的缓存客户端
type FailingCacheClient struct{}

func (f *FailingCacheClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return assert.AnError
}

func (f *FailingCacheClient) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) Delete(ctx context.Context, key string) error {
	return assert.AnError
}

func (f *FailingCacheClient) Exists(ctx context.Context, key string) (bool, error) {
	return false, assert.AnError
}

func (f *FailingCacheClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return assert.AnError
}

func (f *FailingCacheClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) HSet(ctx context.Context, key, field string, value interface{}) error {
	return assert.AnError
}

func (f *FailingCacheClient) HGet(ctx context.Context, key, field string) (interface{}, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) HDel(ctx context.Context, key, field string) error {
	return assert.AnError
}

func (f *FailingCacheClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return false, assert.AnError
}

func (f *FailingCacheClient) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) LPush(ctx context.Context, key string, values ...string) error {
	return assert.AnError
}

func (f *FailingCacheClient) LPop(ctx context.Context, key string) (interface{}, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return assert.AnError
}

func (f *FailingCacheClient) LLen(ctx context.Context, key string) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return assert.AnError
}

func (f *FailingCacheClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return assert.AnError
}

func (f *FailingCacheClient) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return assert.AnError
}

func (f *FailingCacheClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, assert.AnError
}

func (f *FailingCacheClient) ZRem(ctx context.Context, key string, member string) error {
	return assert.AnError
}

func (f *FailingCacheClient) ZCard(ctx context.Context, key string) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) ZScore(ctx context.Context, key, member string) (float64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) Incr(ctx context.Context, key string) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) Decr(ctx context.Context, key string) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, assert.AnError
}

func (f *FailingCacheClient) Ping(ctx context.Context) error {
	return assert.AnError
}

func (f *FailingCacheClient) Close() error {
	return assert.AnError
}
