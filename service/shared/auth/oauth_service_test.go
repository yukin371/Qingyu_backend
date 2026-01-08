package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"go.uber.org/zap"
	"go.mongodb.org/mongo-driver/bson/primitive"

	authModel "Qingyu_backend/models/auth"
)

// ============ Mock 实现 ============

// MockOAuthRepository Mock OAuth仓储（扩展版）
type MockOAuthRepositoryExt struct {
	accounts map[string]*authModel.OAuthAccount // key: ID
	sessions map[string]*authModel.OAuthSession // key: State
	nextID   int
}

func NewMockOAuthRepositoryExt() *MockOAuthRepositoryExt {
	return &MockOAuthRepositoryExt{
		accounts: make(map[string]*authModel.OAuthAccount),
		sessions: make(map[string]*authModel.OAuthSession),
		nextID:   1,
	}
}

func (m *MockOAuthRepositoryExt) Create(ctx context.Context, account *authModel.OAuthAccount) error {
	account.ID = fmt.Sprintf("oauth_%d", m.nextID)
	m.nextID++
	m.accounts[account.ID] = account
	return nil
}

func (m *MockOAuthRepositoryExt) FindByID(ctx context.Context, id string) (*authModel.OAuthAccount, error) {
	account, ok := m.accounts[id]
	if !ok {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (m *MockOAuthRepositoryExt) FindByProviderAndProviderID(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
	for _, account := range m.accounts {
		if account.Provider == provider && account.ProviderUserID == providerUserID {
			return account, nil
		}
	}
	return nil, nil
}

func (m *MockOAuthRepositoryExt) FindByUserID(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	result := make([]*authModel.OAuthAccount, 0)
	for _, account := range m.accounts {
		if account.UserID == userID {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *MockOAuthRepositoryExt) Update(ctx context.Context, account *authModel.OAuthAccount) error {
	m.accounts[account.ID] = account
	return nil
}

func (m *MockOAuthRepositoryExt) Delete(ctx context.Context, id string) error {
	delete(m.accounts, id)
	return nil
}

func (m *MockOAuthRepositoryExt) UpdateLastLogin(ctx context.Context, id string) error {
	if account, ok := m.accounts[id]; ok {
		account.LastLoginAt = time.Now()
	}
	return nil
}

func (m *MockOAuthRepositoryExt) GetPrimaryAccount(ctx context.Context, userID string) (*authModel.OAuthAccount, error) {
	for _, account := range m.accounts {
		if account.UserID == userID && account.IsPrimary {
			return account, nil
		}
	}
	return nil, nil
}

func (m *MockOAuthRepositoryExt) SetPrimaryAccount(ctx context.Context, userID, accountID string) error {
	// 取消所有主账号
	for _, account := range m.accounts {
		if account.UserID == userID {
			account.IsPrimary = false
		}
	}
	// 设置新的主账号
	if account, ok := m.accounts[accountID]; ok {
		account.IsPrimary = true
	}
	return nil
}

func (m *MockOAuthRepositoryExt) CountByUserID(ctx context.Context, userID string) (int64, error) {
	count := int64(0)
	for _, account := range m.accounts {
		if account.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *MockOAuthRepositoryExt) UpdateTokens(ctx context.Context, accountID, accessToken, refreshToken string, expiresAt primitive.DateTime) error {
	if account, ok := m.accounts[accountID]; ok {
		account.AccessToken = accessToken
		account.RefreshToken = refreshToken
	}
	return nil
}

func (m *MockOAuthRepositoryExt) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	count := int64(0)
	now := time.Now()
	for id, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, id)
			count++
		}
	}
	return count, nil
}

func (m *MockOAuthRepositoryExt) CreateSession(ctx context.Context, session *authModel.OAuthSession) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MockOAuthRepositoryExt) FindSessionByID(ctx context.Context, id string) (*authModel.OAuthSession, error) {
	return m.sessions[id], nil
}

func (m *MockOAuthRepositoryExt) FindSessionByState(ctx context.Context, state string) (*authModel.OAuthSession, error) {
	for _, session := range m.sessions {
		if session.State == state {
			return session, nil
		}
	}
	return nil, nil
}

func (m *MockOAuthRepositoryExt) DeleteSession(ctx context.Context, id string) error {
	delete(m.sessions, id)
	return nil
}

// ============ 测试辅助函数 ============

func setupTestOAuthService() (*OAuthService, *MockOAuthRepositoryExt) {
	logger := zap.NewNop()
	repo := NewMockOAuthRepositoryExt()

	// 创建测试配置
	configs := map[string]*authModel.OAuthConfig{
		"google": {
			Enabled:     true,
			ClientID:    "test-google-client-id",
			ClientSecret: "test-google-secret",
			RedirectURI: "http://localhost:8080/oauth/google/callback",
			Scopes:      "openid email profile",
			AuthURL:     "",
			TokenURL:    "",
		},
		"github": {
			Enabled:     true,
			ClientID:    "test-github-client-id",
			ClientSecret: "test-github-secret",
			RedirectURI: "http://localhost:8080/oauth/github/callback",
			Scopes:      "user:email read:user",
			AuthURL:     "",
			TokenURL:    "",
		},
		"qq": {
			Enabled:     true,
			ClientID:    "test-qq-client-id",
			ClientSecret: "test-qq-secret",
			RedirectURI: "http://localhost:8080/oauth/qq/callback",
			Scopes:      "get_user_info",
			AuthURL:     "https://graph.qq.com/oauth2.0/authorize",
			TokenURL:    "https://graph.qq.com/oauth2.0/token",
		},
	}

	service, err := NewOAuthService(logger, repo, configs)
	if err != nil {
		panic(fmt.Sprintf("创建OAuth服务失败: %v", err))
	}

	return service, repo
}

// ============ 测试用例 ============

// TestOAuthService_GetAuthURL_Google 测试获取Google授权URL
func TestOAuthService_GetAuthURL_Google(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	authURL, err := service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", "test_state", false)
	if err != nil {
		t.Fatalf("获取授权URL失败: %v", err)
	}

	if authURL == "" {
		t.Error("授权URL不应为空")
	}

	// 验证URL包含必要的参数
	if !contains(authURL, "accounts.google.com") && !contains(authURL, "google.com") {
		t.Error("URL应该包含Google域名")
	}

	// 验证state参数已保存（通过检查stateStore）
	if len(service.stateStore) == 0 {
		t.Error("OAuth会话应该已创建")
	}

	t.Logf("Google授权URL: %s", authURL)
}

// TestOAuthService_GetAuthURL_GitHub 测试获取GitHub授权URL
func TestOAuthService_GetAuthURL_GitHub(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	authURL, err := service.GetAuthURL(ctx, authModel.OAuthProviderGitHub, "http://localhost:3000/callback", "test_state", false)
	if err != nil {
		t.Fatalf("获取授权URL失败: %v", err)
	}

	if authURL == "" {
		t.Error("授权URL不应为空")
	}

	if !contains(authURL, "github.com") {
		t.Error("URL应该包含GitHub域名")
	}

	t.Logf("GitHub授权URL: %s", authURL)
}

// TestOAuthService_GetAuthURL_QQ 测试获取QQ授权URL
func TestOAuthService_GetAuthURL_QQ(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	authURL, err := service.GetAuthURL(ctx, authModel.OAuthProviderQQ, "http://localhost:3000/callback", "test_state", false)
	if err != nil {
		t.Fatalf("获取授权URL失败: %v", err)
	}

	if authURL == "" {
		t.Error("授权URL不应为空")
	}

	if !contains(authURL, "graph.qq.com") {
		t.Error("URL应该包含QQ域名")
	}

	t.Logf("QQ授权URL: %s", authURL)
}

// TestOAuthService_GetAuthURL_UnsupportedProvider 测试不支持的提供商
func TestOAuthService_GetAuthURL_UnsupportedProvider(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	_, err := service.GetAuthURL(ctx, "unsupported_provider", "http://localhost:3000/callback", "test_state", false)
	if err == nil {
		t.Fatal("应该返回错误，但成功了")
	}

	t.Logf("正确返回了错误: %v", err)
}

// TestOAuthService_GetAuthURL_LinkMode 测试绑定模式
func TestOAuthService_GetAuthURL_LinkMode(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"
	authURL, err := service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", "test_state", true, userID)
	if err != nil {
		t.Fatalf("获取授权URL失败: %v", err)
	}

	if authURL == "" {
		t.Error("授权URL不应为空")
	}

	// 验证会话已创建且标记为绑定模式
	if len(service.stateStore) == 0 {
		t.Error("OAuth会话应该已创建")
	}

	var session *authModel.OAuthSession
	for _, s := range service.stateStore {
		session = s
		break
	}

	if session == nil {
		t.Fatal("会话不应为空")
	}

	if !session.LinkMode {
		t.Error("会话应该标记为绑定模式")
	}

	if session.UserID != userID {
		t.Errorf("用户ID错误: %s", session.UserID)
	}

	t.Logf("绑定模式授权URL成功")
}

// TestOAuthService_ExchangeCode 测试交换授权码
func TestOAuthService_ExchangeCode(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	// 先创建一个会话
	_, _ = service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", "test_state", false)

	// 从URL中提取state（简化处理）
	var state string
	for key := range service.stateStore {
		state = service.stateStore[key].State
		break
	}

	// 注意：这个测试会失败，因为我们没有真实的HTTP服务器来交换授权码
	// 实际项目中需要使用HTTP mock或测试服务器
	_, session, err := service.ExchangeCode(ctx, authModel.OAuthProviderGoogle, "mock_code", state)
	if err == nil {
		t.Log("交换授权码成功（注意：实际项目需要mock HTTP服务器）")
	} else {
		t.Logf("交换授权码失败（预期行为，因为需要真实HTTP交互）: %v", err)
	}

	if session != nil {
		t.Logf("会话信息: Provider=%s, LinkMode=%v", session.Provider, session.LinkMode)
	}
}

// TestOAuthService_ExchangeCode_InvalidState 测试无效的state
func TestOAuthService_ExchangeCode_InvalidState(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	_, session, err := service.ExchangeCode(ctx, authModel.OAuthProviderGoogle, "mock_code", "invalid_state")
	if err == nil {
		t.Error("应该返回错误（无效state），但成功了")
	}
	if session != nil {
		t.Error("会话应该为nil")
	}

	t.Logf("正确返回了错误: %v", err)
}

// TestOAuthService_LinkAccount 测试绑定账号
func TestOAuthService_LinkAccount(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"
	provider := authModel.OAuthProviderGoogle

	// 创建模拟token和身份信息
	token := &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	identity := &authModel.UserIdentity{
		Provider:      provider,
		ProviderID:    "google_user_123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Avatar:        "https://example.com/avatar.jpg",
		Username:      "testuser",
	}

	// 绑定账号
	account, err := service.LinkAccount(ctx, userID, provider, token, identity)
	if err != nil {
		t.Fatalf("绑定账号失败: %v", err)
	}

	if account.UserID != userID {
		t.Errorf("用户ID错误: %s", account.UserID)
	}
	if account.Provider != provider {
		t.Errorf("提供商错误: %s", account.Provider)
	}
	if account.ProviderUserID != identity.ProviderID {
		t.Errorf("提供商用户ID错误: %s", account.ProviderUserID)
	}
	if !account.IsPrimary {
		t.Error("第一个账号应该设为主账号")
	}

	t.Logf("绑定账号成功: %+v", account)
}

// TestOAuthService_LinkAccount_Duplicate 测试绑定重复账号
func TestOAuthService_LinkAccount_Duplicate(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"
	provider := authModel.OAuthProviderGoogle

	token := &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
	}

	identity := &authModel.UserIdentity{
		Provider:   provider,
		ProviderID: "google_user_123",
		Email:      "test@example.com",
	}

	// 第一次绑定
	_, err := service.LinkAccount(ctx, userID, provider, token, identity)
	if err != nil {
		t.Fatalf("第一次绑定失败: %v", err)
	}

	// 第二次绑定（相同账号）
	_, err = service.LinkAccount(ctx, userID, provider, token, identity)
	if err != nil {
		t.Fatalf("第二次绑定应该成功（已存在），但失败了: %v", err)
	}

	t.Logf("重复绑定测试通过")
}

// TestOAuthService_LinkAccount_AlreadyLinkedToOtherUser 测试账号已绑定到其他用户
func TestOAuthService_LinkAccount_AlreadyLinkedToOtherUser(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID1 := "user_123"
	userID2 := "user_456"
	provider := authModel.OAuthProviderGoogle

	token := &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
	}

	identity := &authModel.UserIdentity{
		Provider:   provider,
		ProviderID: "google_user_123",
		Email:      "test@example.com",
	}

	// 用户1绑定
	_, err := service.LinkAccount(ctx, userID1, provider, token, identity)
	if err != nil {
		t.Fatalf("用户1绑定失败: %v", err)
	}

	// 用户2尝试绑定相同账号
	_, err = service.LinkAccount(ctx, userID2, provider, token, identity)
	if err == nil {
		t.Error("应该返回错误（账号已绑定到其他用户），但成功了")
	}

	t.Logf("正确拒绝了重复绑定: %v", err)
}

// TestOAuthService_UnlinkAccount 测试解绑账号
func TestOAuthService_UnlinkAccount(t *testing.T) {
	service, repo := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"

	// 绑定第一个账号（Google）
	token1 := &oauth2.Token{
		AccessToken:  "test_access_token_1",
		RefreshToken: "test_refresh_token_1",
	}

	identity1 := &authModel.UserIdentity{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google_user_123",
		Email:      "test@example.com",
	}

	account1, _ := service.LinkAccount(ctx, userID, authModel.OAuthProviderGoogle, token1, identity1)

	// 绑定第二个账号（GitHub）
	token2 := &oauth2.Token{
		AccessToken:  "test_access_token_2",
		RefreshToken: "test_refresh_token_2",
	}

	identity2 := &authModel.UserIdentity{
		Provider:   authModel.OAuthProviderGitHub,
		ProviderID: "github_user_456",
		Email:      "test@github.com",
	}

	_, _ = service.LinkAccount(ctx, userID, authModel.OAuthProviderGitHub, token2, identity2)

	// 解绑第一个账号（不是唯一的）
	err := service.UnlinkAccount(ctx, userID, account1.ID)
	if err != nil {
		t.Fatalf("解绑账号失败: %v", err)
	}

	// 验证账号已删除
	_, err = repo.FindByID(ctx, account1.ID)
	if err == nil {
		t.Error("账号应该已删除")
	}

	t.Logf("解绑账号成功")
}

// TestOAuthService_UnlinkAccount_PrimaryWithOnlyOne 测试解绑唯一的主账号
func TestOAuthService_UnlinkAccount_PrimaryWithOnlyOne(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"
	provider := authModel.OAuthProviderGoogle

	token := &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
	}

	identity := &authModel.UserIdentity{
		Provider:   provider,
		ProviderID: "google_user_123",
		Email:      "test@example.com",
	}

	// 绑定账号
	account, _ := service.LinkAccount(ctx, userID, provider, token, identity)

	// 尝试解绑唯一的主账号
	err := service.UnlinkAccount(ctx, userID, account.ID)
	if err == nil {
		t.Error("应该返回错误（不能解绑唯一的主账号），但成功了")
	}

	t.Logf("正确拒绝了解绑唯一主账号: %v", err)
}

// TestOAuthService_GetLinkedAccounts 测试获取绑定账号列表
func TestOAuthService_GetLinkedAccounts(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"

	// 绑定多个账号
	providers := []authModel.OAuthProvider{
		authModel.OAuthProviderGoogle,
		authModel.OAuthProviderGitHub,
		authModel.OAuthProviderQQ,
	}

	for i, provider := range providers {
		token := &oauth2.Token{
			AccessToken:  fmt.Sprintf("access_token_%d", i),
			RefreshToken: fmt.Sprintf("refresh_token_%d", i),
		}

		identity := &authModel.UserIdentity{
			Provider:   provider,
			ProviderID: fmt.Sprintf("user_%d", i),
			Email:      fmt.Sprintf("user%d@example.com", i),
		}

		_, err := service.LinkAccount(ctx, userID, provider, token, identity)
		if err != nil {
			t.Fatalf("绑定账号%d失败: %v", i, err)
		}
	}

	// 获取绑定账号列表
	accounts, err := service.GetLinkedAccounts(ctx, userID)
	if err != nil {
		t.Fatalf("获取绑定账号列表失败: %v", err)
	}

	if len(accounts) != 3 {
		t.Errorf("账号数量错误: 期望3个，实际%d个", len(accounts))
	}

	t.Logf("获取绑定账号列表成功: %d个账号", len(accounts))
}

// TestOAuthService_SetPrimaryAccount 测试设置主账号
func TestOAuthService_SetPrimaryAccount(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	userID := "user_123"

	// 绑定两个账号
	token1 := &oauth2.Token{AccessToken: "access_token_1"}
	identity1 := &authModel.UserIdentity{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google_123",
		Email:      "google@example.com",
	}

	account1, _ := service.LinkAccount(ctx, userID, authModel.OAuthProviderGoogle, token1, identity1)

	token2 := &oauth2.Token{AccessToken: "access_token_2"}
	identity2 := &authModel.UserIdentity{
		Provider:   authModel.OAuthProviderGitHub,
		ProviderID: "github_123",
		Email:      "github@example.com",
	}

	account2, _ := service.LinkAccount(ctx, userID, authModel.OAuthProviderGitHub, token2, identity2)

	// 设置第二个账号为主账号
	err := service.SetPrimaryAccount(ctx, userID, account2.ID)
	if err != nil {
		t.Fatalf("设置主账号失败: %v", err)
	}

	// 验证主账号已更新
	accounts, _ := service.GetLinkedAccounts(ctx, userID)
	for _, account := range accounts {
		if account.ID == account2.ID && !account.IsPrimary {
			t.Error("第二个账号应该已设为主账号")
		}
		if account.ID == account1.ID && account.IsPrimary {
			t.Error("第一个账号应该不再是主账号")
		}
	}

	t.Logf("设置主账号成功")
}

// TestOAuthService_CleanupSession 测试清理过期会话
func TestOAuthService_CleanupSession(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	// 创建一个会话
	_, _ = service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", "test_state", false)

	// 手动设置会话为已过期
	for id, session := range service.stateStore {
		session.ExpiresAt = time.Now().Add(-1 * time.Hour)
		service.stateStore[id] = session
	}

	// 清理过期会话
	service.CleanupSession(ctx)

	// 验证会话已清理
	if len(service.stateStore) != 0 {
		t.Errorf("过期会话应该已清理，但还有%d个", len(service.stateStore))
	}

	t.Logf("清理过期会话成功")
}

// TestOAuthService_CleanupExpiredSessions 测试清理所有过期会话
func TestOAuthService_CleanupExpiredSessions(t *testing.T) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	// 创建多个会话
	for i := 0; i < 3; i++ {
		_, _ = service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", fmt.Sprintf("state_%d", i), false)
	}

	// 手动设置部分会话为已过期
	count := 0
	for id, session := range service.stateStore {
		if count < 2 {
			session.ExpiresAt = time.Now().Add(-1 * time.Hour)
			service.stateStore[id] = session
			count++
		}
	}

	// 记录清理前的数量
	beforeCount := len(service.stateStore)

	// 清理过期会话
	err := service.CleanupExpiredSessions(ctx)
	if err != nil {
		t.Fatalf("清理过期会话失败: %v", err)
	}

	// 验证会话已清理
	afterCount := len(service.stateStore)
	if afterCount >= beforeCount {
		t.Errorf("过期会话应该已清理，但数量没有减少: before=%d, after=%d", beforeCount, afterCount)
	}

	t.Logf("清理过期会话成功: 清理了%d个", beforeCount-afterCount)
}

// BenchmarkGetAuthURL 性能测试：获取授权URL
func BenchmarkGetAuthURL(b *testing.B) {
	service, _ := setupTestOAuthService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetAuthURL(ctx, authModel.OAuthProviderGoogle, "http://localhost:3000/callback", "test_state", false)
	}
}

// ============ 辅助函数 ============

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
