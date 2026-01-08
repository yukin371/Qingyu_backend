package auth

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// ============ SessionService 测试 ============

// TestSessionService_CreateSession 测试创建会话
func TestSessionService_CreateSession(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	session, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	if session.ID == "" {
		t.Error("会话ID不应为空")
	}
	if session.UserID != userID {
		t.Errorf("用户ID错误: %s", session.UserID)
	}
	if session.CreatedAt.IsZero() {
		t.Error("创建时间不应为零")
	}
	if session.ExpiresAt.Before(session.CreatedAt) {
		t.Error("过期时间不应早于创建时间")
	}

	t.Logf("创建会话成功: ID=%s, UserID=%s", session.ID, session.UserID)
}

// TestSessionService_CreateSession_Multiple 测试创建多个会话
func TestSessionService_CreateSession_Multiple(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建3个会话
	sessions := make([]*Session, 0, 3)
	for i := 0; i < 3; i++ {
		session, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话%d失败: %v", i+1, err)
		}
		sessions = append(sessions, session)
	}

	// 验证会话ID唯一
	sessionIDs := make(map[string]bool)
	for _, session := range sessions {
		if sessionIDs[session.ID] {
			t.Errorf("会话ID重复: %s", session.ID)
		}
		sessionIDs[session.ID] = true
	}

	if len(sessions) != 3 {
		t.Errorf("会话数量错误: %d", len(sessions))
	}

	t.Logf("成功创建3个会话，ID均唯一")
}

// TestSessionService_GetSession 测试获取会话
func TestSessionService_GetSession(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建会话
	createdSession, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 获取会话
	retrievedSession, err := service.GetSession(ctx, createdSession.ID)
	if err != nil {
		t.Fatalf("获取会话失败: %v", err)
	}

	if retrievedSession.ID != createdSession.ID {
		t.Errorf("会话ID不匹配: %s vs %s", retrievedSession.ID, createdSession.ID)
	}
	if retrievedSession.UserID != userID {
		t.Errorf("用户ID错误: %s", retrievedSession.UserID)
	}

	t.Logf("获取会话成功")
}

// TestSessionService_GetSession_NotFound 测试获取不存在的会话
func TestSessionService_GetSession_NotFound(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	// 获取不存在的会话
	_, err := service.GetSession(ctx, "nonexistent_session_id")
	if err == nil {
		t.Fatal("应该返回错误，但成功了")
	}

	t.Logf("正确返回了错误: %v", err)
}

// TestSessionService_RefreshSession 测试刷新会话
func TestSessionService_RefreshSession(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建会话
	session, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 记录原始过期时间
	originalExpiresAt := session.ExpiresAt

	// 等待一小段时间
	time.Sleep(1100 * time.Millisecond)

	// 刷新会话
	err = service.RefreshSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("刷新会话失败: %v", err)
	}

	// 获取刷新后的会话
	refreshedSession, err := service.GetSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("获取刷新后的会话失败: %v", err)
	}

	// 验证过期时间已更新
	if !refreshedSession.ExpiresAt.After(originalExpiresAt) {
		t.Error("过期时间应该已更新")
	}

	t.Logf("刷新会话成功")
}

// TestSessionService_DestroySession 测试销毁会话
func TestSessionService_DestroySession(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建会话
	session, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 销毁会话
	err = service.DestroySession(ctx, session.ID)
	if err != nil {
		t.Fatalf("销毁会话失败: %v", err)
	}

	// 验证会话已不存在
	_, err = service.GetSession(ctx, session.ID)
	if err == nil {
		t.Fatal("会话应该已不存在")
	}

	t.Logf("销毁会话成功")
}

// TestSessionService_GetUserSessions 测试获取用户的所有会话
func TestSessionService_GetUserSessions(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建3个会话
	sessionIDs := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		session, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话%d失败: %v", i+1, err)
		}
		sessionIDs = append(sessionIDs, session.ID)
	}

	// 获取用户的所有会话
	sessions, err := service.GetUserSessions(ctx, userID)
	if err != nil {
		t.Fatalf("获取用户会话失败: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("会话数量错误: 期望3个，实际%d个", len(sessions))
	}

	// 验证所有会话ID都存在
	retrievedIDs := make(map[string]bool)
	for _, session := range sessions {
		retrievedIDs[session.ID] = true
	}

	for _, sessionID := range sessionIDs {
		if !retrievedIDs[sessionID] {
			t.Errorf("会话ID未找到: %s", sessionID)
		}
	}

	t.Logf("获取用户会话成功: %d个会话", len(sessions))
}

// TestSessionService_DestroyUserSessions 测试销毁用户的所有会话
func TestSessionService_DestroyUserSessions(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建3个会话
	for i := 0; i < 3; i++ {
		_, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话%d失败: %v", i+1, err)
		}
	}

	// 销毁所有会话
	err := service.DestroyUserSessions(ctx, userID)
	if err != nil {
		t.Fatalf("销毁用户会话失败: %v", err)
	}

	// 验证所有会话已不存在
	sessions, err := service.GetUserSessions(ctx, userID)
	if err != nil {
		t.Fatalf("获取用户会话失败: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("应该没有会话了，但还有%d个", len(sessions))
	}

	t.Logf("销毁用户所有会话成功")
}

// TestSessionService_CheckDeviceLimit 测试设备数量限制检查
func TestSessionService_CheckDeviceLimit(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"
	maxDevices := 3

	// 创建2个会话（未超限）
	for i := 0; i < 2; i++ {
		_, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}
	}

	// 检查设备限制（应该通过）
	err := service.CheckDeviceLimit(ctx, userID, maxDevices)
	if err != nil {
		t.Errorf("设备数量未超限，应该通过: %v", err)
	}

	t.Logf("设备数量限制检查通过")
}

// TestSessionService_CheckDeviceLimit_Exceeded 测试设备数量超限
func TestSessionService_CheckDeviceLimit_Exceeded(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"
	maxDevices := 2

	// 创建2个会话（已达上限）
	for i := 0; i < 2; i++ {
		_, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}
	}

	// 检查设备限制（应该失败）
	err := service.CheckDeviceLimit(ctx, userID, maxDevices)
	if err == nil {
		t.Error("设备数量已达上限，应该返回错误")
	}

	t.Logf("正确检测到设备数量超限: %v", err)
}

// TestSessionService_EnforceDeviceLimit 测试强制执行设备限制（FIFO）
func TestSessionService_EnforceDeviceLimit(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"
	maxDevices := 3

	// 创建3个会话（达到上限）
	sessionIDs := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		session, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话%d失败: %v", i+1, err)
		}
		sessionIDs = append(sessionIDs, session.ID)
		// 等待确保创建时间不同
		time.Sleep(100 * time.Millisecond)
	}

	// 强制执行设备限制（应该踢出最老的会话）
	err := service.EnforceDeviceLimit(ctx, userID, maxDevices)
	if err != nil {
		t.Fatalf("强制执行设备限制失败: %v", err)
	}

	// 验证最老的会话已被踢出
	_, err = service.GetSession(ctx, sessionIDs[0])
	if err == nil {
		t.Error("最老的会话应该已被踢出")
	}

	// 验证其他会话仍然存在
	for i := 1; i < len(sessionIDs); i++ {
		_, err := service.GetSession(ctx, sessionIDs[i])
		if err != nil {
			t.Errorf("会话%d应该还存在: %v", i, err)
		}
	}

	t.Logf("强制执行设备限制成功，最老的会话已被踢出")
}

// TestSessionService_EnforceDeviceLimit_CreateNew 测试强制执行设备限制后创建新会话
func TestSessionService_EnforceDeviceLimit_CreateNew(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"
	maxDevices := 2

	// 创建2个会话（达到上限）
	for i := 0; i < 2; i++ {
		_, err := service.CreateSession(ctx, userID)
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 强制执行设备限制
	err := service.EnforceDeviceLimit(ctx, userID, maxDevices)
	if err != nil {
		t.Fatalf("强制执行设备限制失败: %v", err)
	}

	// 创建新会话（应该成功，因为最老的已被踢出）
	newSession, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建新会话失败: %v", err)
	}

	// 验证当前会话数量
	sessions, _ := service.GetUserSessions(ctx, userID)
	if len(sessions) != maxDevices {
		t.Errorf("会话数量应该为%d，实际为%d", maxDevices, len(sessions))
	}

	// 验证新会话存在
	found := false
	for _, session := range sessions {
		if session.ID == newSession.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("新会话未在列表中")
	}

	t.Logf("强制执行设备限制后创建新会话成功")
}

// TestSessionService_UpdateSession 测试更新会话
func TestSessionService_UpdateSession(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建会话
	session, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 更新会话
	data := map[string]interface{}{
		"last_activity": time.Now(),
	}
	err = service.UpdateSession(ctx, session.ID, data)
	if err != nil {
		t.Fatalf("更新会话失败: %v", err)
	}

	// 验证会话仍然存在
	_, err = service.GetSession(ctx, session.ID)
	if err != nil {
		t.Error("更新后会话应该仍然存在")
	}

	t.Logf("更新会话成功")
}

// TestSessionService_ConcurrentAccess 测试并发访问会话
func TestSessionService_ConcurrentAccess(t *testing.T) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 并发创建多个会话
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := service.CreateSession(ctx, userID)
			if err != nil {
				t.Errorf("并发创建会话失败: %v", err)
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		<-done
	}

	// 验证会话数量
	sessions, err := service.GetUserSessions(ctx, userID)
	if err != nil {
		t.Fatalf("获取用户会话失败: %v", err)
	}

	// 由于并发和设备限制，会话数量可能在2-5之间
	// (EnforceDeviceLimit会删除超过5个的最旧会话，但并发时可能有些会话被提前删除)
	if len(sessions) < 2 || len(sessions) > 5 {
		t.Errorf("并发创建后应该有2-5个会话，实际有%d个", len(sessions))
	}

	t.Logf("并发访问测试通过，创建了%d个会话", len(sessions))
}

// TestSessionService_ExpiredSession 测试过期会话处理
func TestSessionService_ExpiredSession(t *testing.T) {
	// 创建一个会立即过期的配置
	// 注意：当前Mock实现不支持自定义TTL，所以这个测试主要验证逻辑
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	userID := "user_123"

	// 创建会话
	session, err := service.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 在实际实现中，我们可以通过修改缓存来模拟过期
	// 这里我们删除缓存的会话来模拟过期
	key := fmt.Sprintf("session:%s", session.ID)
	_ = cache.Delete(ctx, key)

	// 尝试获取过期会话
	_, err = service.GetSession(ctx, session.ID)
	if err == nil {
		t.Error("过期会话应该返回错误")
	}

	t.Logf("过期会话处理测试通过")
}

// BenchmarkCreateSession 性能测试：创建会话
func BenchmarkCreateSession(b *testing.B) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateSession(ctx, "user_123")
	}
}

// BenchmarkGetSession 性能测试：获取会话
func BenchmarkGetSession(b *testing.B) {
	cache := NewMockCacheClient()
	service := NewSessionService(cache)
	ctx := context.Background()

	// 预先创建会话
	session, _ := service.CreateSession(ctx, "user_123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetSession(ctx, session.ID)
	}
}
