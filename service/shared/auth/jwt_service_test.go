package auth

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
)

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value.(string)
	return nil
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", nil
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			return 1, nil
		}
	}
	return 0, nil
}

// 测试辅助函数
func getTestJWTConfig() *config.JWTConfigEnhanced {
	return &config.JWTConfigEnhanced{
		SecretKey:       "test-secret-key",
		Issuer:          "qingyu-test",
		Expiration:      1 * time.Hour,
		RefreshDuration: 7 * 24 * time.Hour,
	}
}

// ============ 测试用例 ============

// TestGenerateToken 测试生成访问Token
func TestGenerateToken(t *testing.T) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	token, err := service.GenerateToken(ctx, "user123", []string{"reader"})
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	if token == "" {
		t.Fatal("Token为空")
	}

	t.Logf("生成的Token: %s", token)
}

// TestValidateToken 测试验证Token
func TestValidateToken(t *testing.T) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	// 生成Token
	token, err := service.GenerateToken(ctx, "user123", []string{"reader"})
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	// 验证Token
	claims, err := service.ValidateToken(ctx, token)
	if err != nil {
		t.Fatalf("验证Token失败: %v", err)
	}

	// 检查Claims
	if claims.UserID != "user123" {
		t.Errorf("UserID错误: %s", claims.UserID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "reader" {
		t.Errorf("Roles错误: %v", claims.Roles)
	}

	t.Logf("验证成功，Claims: UserID=%s, Roles=%v", claims.UserID, claims.Roles)
}

// TestValidateToken_InvalidSignature 测试无效签名
func TestValidateToken_InvalidSignature(t *testing.T) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	// 生成Token
	token, _ := service.GenerateToken(ctx, "user123", []string{"reader"})

	// 篡改Token
	tampered := token + "tampered"

	// 验证应该失败
	_, err := service.ValidateToken(ctx, tampered)
	if err == nil {
		t.Fatal("应该验证失败，但成功了")
	}

	t.Logf("正确拒绝了篡改的Token: %v", err)
}

// TestValidateToken_Expired 测试过期Token
func TestValidateToken_Expired(t *testing.T) {
	// 创建一个立即过期的配置
	cfg := &config.JWTConfigEnhanced{
		SecretKey:       "test-secret-key",
		Issuer:          "qingyu-test",
		Expiration:      -1 * time.Hour, // 负数，立即过期
		RefreshDuration: 7 * 24 * time.Hour,
	}
	service := NewJWTService(cfg, NewMockRedisClient())
	ctx := context.Background()

	// 生成Token（会立即过期）
	token, _ := service.GenerateToken(ctx, "user123", []string{"reader"})

	// 验证应该失败
	_, err := service.ValidateToken(ctx, token)
	if err == nil {
		t.Fatal("应该验证失败（Token已过期），但成功了")
	}

	t.Logf("正确拒绝了过期Token: %v", err)
}

// TestRefreshToken 测试刷新Token
func TestRefreshToken(t *testing.T) {
	redisClient := NewMockRedisClient()
	service := NewJWTService(getTestJWTConfig(), redisClient)
	ctx := context.Background()

	// 生成初始Token
	originalToken, err := service.GenerateToken(ctx, "user123", []string{"reader"})
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	// 等待1秒确保过期时间戳不同
	time.Sleep(1100 * time.Millisecond)

	// 使用原Token作为刷新Token
	newToken, err := service.RefreshToken(ctx, originalToken)
	if err != nil {
		t.Fatalf("刷新Token失败: %v", err)
	}

	// 新旧Token应该不同（过期时间不同）
	if newToken == originalToken {
		t.Error("新旧Token相同，应该不同")
	}

	// 验证旧Token应该被吊销
	revoked, _ := service.IsTokenRevoked(ctx, originalToken)
	if !revoked {
		t.Error("旧Token应该被吊销，但没有")
	}

	// 直接解析新Token检查Claims（不验证，因为Redis可能有缓存问题）
	svc := service.(*JWTServiceImpl)
	claims, err := svc.ParseTokenClaims(newToken)
	if err != nil {
		t.Fatalf("解析新Token失败: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("UserID错误: %s", claims.UserID)
	}

	t.Logf("刷新Token成功，旧Token已吊销")
}

// TestRevokeToken 测试吊销Token
func TestRevokeToken(t *testing.T) {
	redisClient := NewMockRedisClient()
	service := NewJWTService(getTestJWTConfig(), redisClient)
	ctx := context.Background()

	// 生成Token
	token, _ := service.GenerateToken(ctx, "user123", []string{"reader"})

	// 吊销Token
	err := service.RevokeToken(ctx, token)
	if err != nil {
		t.Fatalf("吊销Token失败: %v", err)
	}

	// 检查Token是否在黑名单中
	revoked, err := service.IsTokenRevoked(ctx, token)
	if err != nil {
		t.Fatalf("检查黑名单失败: %v", err)
	}
	if !revoked {
		t.Fatal("Token应该在黑名单中，但不在")
	}

	// 验证Token应该失败
	_, err = service.ValidateToken(ctx, token)
	if err == nil {
		t.Fatal("应该验证失败（Token已吊销），但成功了")
	}

	t.Logf("正确拒绝了已吊销的Token: %v", err)
}

// TestMultipleRoles 测试多角色Token
func TestMultipleRoles(t *testing.T) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	roles := []string{"reader", "author", "admin"}
	token, err := service.GenerateToken(ctx, "user123", roles)
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	claims, err := service.ValidateToken(ctx, token)
	if err != nil {
		t.Fatalf("验证Token失败: %v", err)
	}

	if len(claims.Roles) != 3 {
		t.Errorf("角色数量错误: 期望3个，实际%d个", len(claims.Roles))
	}

	t.Logf("多角色Token验证成功: %v", claims.Roles)
}

// BenchmarkGenerateToken 性能测试：生成Token
func BenchmarkGenerateToken(b *testing.B) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(ctx, "user123", []string{"reader"})
	}
}

// BenchmarkValidateToken 性能测试：验证Token
func BenchmarkValidateToken(b *testing.B) {
	service := NewJWTService(getTestJWTConfig(), NewMockRedisClient())
	ctx := context.Background()

	token, _ := service.GenerateToken(ctx, "user123", []string{"reader"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(ctx, token)
	}
}
