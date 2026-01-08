package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	authModel "Qingyu_backend/models/auth"
	"Qingyu_backend/service/shared/auth"
)

// ============ Mock 实现 ============

// MockOAuthServiceForAPI Mock OAuth服务（用于API测试）
type MockOAuthServiceForAPI struct {
	authURL        string
	sessions       map[string]*authModel.OAuthSession
	linkedAccounts map[string][]*authModel.OAuthAccount // userID -> accounts
	shouldError    bool
}

func NewMockOAuthServiceForAPI() *MockOAuthServiceForAPI {
	return &MockOAuthServiceForAPI{
		authURL:        "https://accounts.google.com/o/oauth2/v2/auth?redirect_uri=http://localhost:3000/callback",
		sessions:       make(map[string]*authModel.OAuthSession),
		linkedAccounts: make(map[string][]*authModel.OAuthAccount),
		shouldError:    false,
	}
}

func (m *MockOAuthServiceForAPI) GetAuthURL(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error) {
	if m.shouldError {
		return "", &testError{msg: "获取授权URL失败"}
	}

	session := &authModel.OAuthSession{
		ID:          "session_" + state,
		State:       state,
		Provider:    provider,
		RedirectURI: redirectURI,
		LinkMode:    linkMode,
	}
	if linkMode && len(userID) > 0 {
		session.UserID = userID[0]
	}

	m.sessions[state] = session
	return m.authURL + "&state=" + state, nil
}

func (m *MockOAuthServiceForAPI) ExchangeCode(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
	if m.shouldError {
		return nil, nil, &testError{msg: "交换授权码失败"}
	}

	session, ok := m.sessions[state]
	if !ok {
		return nil, nil, &testError{msg: "无效的state"}
	}

	token := &oauth2.Token{
		AccessToken:  "test_access_token_" + state,
		RefreshToken: "test_refresh_token_" + state,
		Expiry:       nil,
	}

	return token, session, nil
}

func (m *MockOAuthServiceForAPI) GetUserInfo(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
	if m.shouldError {
		return nil, &testError{msg: "获取用户信息失败"}
	}

	return &authModel.UserIdentity{
		Provider:      provider,
		ProviderID:    "provider_user_123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Avatar:        "https://example.com/avatar.jpg",
		Username:      "testuser",
	}, nil
}

func (m *MockOAuthServiceForAPI) LinkAccount(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
	if m.shouldError {
		return nil, &testError{msg: "绑定账号失败"}
	}

	account := &authModel.OAuthAccount{
		ID:             "oauth_" + userID,
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: identity.ProviderID,
		Email:          identity.Email,
		Username:       identity.Username,
		Avatar:         identity.Avatar,
		IsPrimary:      true,
	}

	m.linkedAccounts[userID] = append(m.linkedAccounts[userID], account)
	return account, nil
}

func (m *MockOAuthServiceForAPI) UnlinkAccount(ctx context.Context, userID, accountID string) error {
	if m.shouldError {
		return &testError{msg: "解绑账号失败"}
	}

	accounts := m.linkedAccounts[userID]
	newAccounts := make([]*authModel.OAuthAccount, 0, len(accounts))
	for _, acc := range accounts {
		if acc.ID != accountID {
			newAccounts = append(newAccounts, acc)
		}
	}
	m.linkedAccounts[userID] = newAccounts
	return nil
}

func (m *MockOAuthServiceForAPI) GetLinkedAccounts(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	if m.shouldError {
		return nil, &testError{msg: "获取绑定账号失败"}
	}

	accounts := m.linkedAccounts[userID]
	if accounts == nil {
		return []*authModel.OAuthAccount{}, nil
	}
	return accounts, nil
}

func (m *MockOAuthServiceForAPI) SetPrimaryAccount(ctx context.Context, userID, accountID string) error {
	if m.shouldError {
		return &testError{msg: "设置主账号失败"}
	}

	accounts := m.linkedAccounts[userID]
	for _, acc := range accounts {
		acc.IsPrimary = (acc.ID == accountID)
	}
	return nil
}

// MockAuthServiceForOAuth Mock AuthService（用于OAuth API测试）
type MockAuthServiceForOAuth struct {
	oauthLoginResp *auth.LoginResponse
	shouldError    bool
}

func NewMockAuthServiceForOAuth() *MockAuthServiceForOAuth {
	return &MockAuthServiceForOAuth{
		oauthLoginResp: &auth.LoginResponse{
			User: &auth.UserInfo{
				ID:       "user_123",
				Username: "testuser",
				Email:    "test@example.com",
				Roles:    []string{"reader"},
			},
			Token: "test_jwt_token",
		},
		shouldError: false,
	}
}

func (m *MockAuthServiceForOAuth) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) OAuthLogin(ctx context.Context, req *auth.OAuthLoginRequest) (*auth.LoginResponse, error) {
	if m.shouldError {
		return nil, &testError{msg: "OAuth登录失败"}
	}
	return m.oauthLoginResp, nil
}

func (m *MockAuthServiceForOAuth) Logout(ctx context.Context, token string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) RefreshToken(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (m *MockAuthServiceForOAuth) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return true, nil
}

func (m *MockAuthServiceForOAuth) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (m *MockAuthServiceForOAuth) HasRole(ctx context.Context, userID, role string) (bool, error) {
	return true, nil
}

func (m *MockAuthServiceForOAuth) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (m *MockAuthServiceForOAuth) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	return nil
}

func (m *MockAuthServiceForOAuth) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) RemoveRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForOAuth) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForOAuth) Health(ctx context.Context) error {
	return nil
}

// ============ 测试辅助函数 ============

func setupTestOAuthRouter(oauthAPI *OAuthAPI) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/oauth/:provider/authorize", oauthAPI.GetAuthorizeURL)
	router.POST("/oauth/:provider/callback", oauthAPI.HandleCallback)
	router.GET("/oauth/accounts", oauthAPI.GetLinkedAccounts)
	router.DELETE("/oauth/accounts/:accountID", oauthAPI.UnlinkAccount)
	router.PUT("/oauth/accounts/:accountID/primary", oauthAPI.SetPrimaryAccount)

	return router
}

func makeOAuthRequest(router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// ============ 测试用例 ============

// TestOAuthAPI_GetAuthorizeURL_Google 测试获取Google授权URL
func TestOAuthAPI_GetAuthorizeURL_Google(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	reqBody := map[string]interface{}{
		"redirect_uri": "http://localhost:3000/callback",
		"state":        "test_state_123",
	}

	w := makeOAuthRequest(router, "POST", "/oauth/google/authorize", reqBody, nil)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("响应码错误: %d", resp.Code)
	}

	data := resp.Data.(map[string]interface{})
	if data["authorize_url"] == nil {
		t.Error("响应数据中应该包含authorize_url")
	}

	t.Logf("获取授权URL成功: %+v", resp)
}

// TestOAuthAPI_GetAuthorizeURL_GitHub 测试获取GitHub授权URL
func TestOAuthAPI_GetAuthorizeURL_GitHub(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	reqBody := map[string]interface{}{
		"redirect_uri": "http://localhost:3000/callback",
		"state":        "test_state_456",
	}

	w := makeOAuthRequest(router, "POST", "/oauth/github/authorize", reqBody, nil)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "获取授权URL成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("获取GitHub授权URL成功")
}

// TestOAuthAPI_GetAuthorizeURL_MissingRedirectURI 测试缺失redirect_uri
func TestOAuthAPI_GetAuthorizeURL_MissingRedirectURI(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	reqBody := map[string]interface{}{
		"state": "test_state",
	}

	w := makeOAuthRequest(router, "POST", "/oauth/google/authorize", reqBody, nil)

	if w.Code == http.StatusOK {
		t.Error("应该返回错误，但成功了")
	}

	t.Logf("正确拒绝了缺失redirect_uri的请求")
}

// TestOAuthAPI_HandleCallback 测试处理OAuth回调（登录模式）
func TestOAuthAPI_HandleCallback(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 先获取授权URL创建会话
	authReq := map[string]interface{}{
		"redirect_uri": "http://localhost:3000/callback",
		"state":        "test_state_789",
	}
	makeOAuthRequest(router, "POST", "/oauth/google/authorize", authReq, nil)

	// 测试回调
	callbackReq := map[string]interface{}{
		"code":  "test_authorization_code",
		"state": "test_state_789",
	}

	w := makeOAuthRequest(router, "POST", "/oauth/google/callback", callbackReq, nil)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "OAuth登录成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("OAuth回调处理成功: %+v", resp)
}

// TestOAuthAPI_HandleCallback_LinkMode 测试处理OAuth回调（绑定模式）
func TestOAuthAPI_HandleCallback_LinkMode(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 创建绑定模式的会话
	authReq := map[string]interface{}{
		"redirect_uri": "http://localhost:3000/callback",
		"state":        "test_state_link",
	}

	// 设置用户ID（模拟已登录）
	req, _ := http.NewRequest("POST", "/oauth/google/authorize", nil)
	req.Header.Set("Content-Type", "application/json")
	bodyBytes, _ := json.Marshal(authReq)
	req.Body = httptest.NewRequest("POST", "/oauth/google/authorize", bytes.NewBuffer(bodyBytes)).Body

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	oauthAPI.GetAuthorizeURL(c)

	// 测试回调（绑定模式）
	callbackReq := map[string]interface{}{
		"code":  "test_authorization_code",
		"state": "test_state_link",
	}

	w2 := makeOAuthRequest(router, "POST", "/oauth/google/callback", callbackReq, nil)

	if w2.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w2.Code, w2.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "账号绑定成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("OAuth绑定模式回调成功")
}

// TestOAuthAPI_HandleCallback_InvalidState 测试无效的state
func TestOAuthAPI_HandleCallback_InvalidState(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	callbackReq := map[string]interface{}{
		"code":  "test_authorization_code",
		"state": "invalid_state",
	}

	w := makeOAuthRequest(router, "POST", "/oauth/google/callback", callbackReq, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("应该返回400错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了无效state: %s", w.Body.String())
}

// TestOAuthAPI_GetLinkedAccounts 测试获取绑定账号列表
func TestOAuthAPI_GetLinkedAccounts(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 添加绑定账号
	account := &authModel.OAuthAccount{
		ID:             "oauth_123",
		UserID:         "user_123",
		Provider:       authModel.OAuthProviderGoogle,
		ProviderUserID: "google_user_123",
		Email:          "test@example.com",
		IsPrimary:      true,
	}
	mockOAuthService.linkedAccounts["user_123"] = append(mockOAuthService.linkedAccounts["user_123"], account)

	// 创建带用户ID的请求
	req, _ := http.NewRequest("GET", "/oauth/accounts", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	oauthAPI.GetLinkedAccounts(c)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("响应码错误: %d", resp.Code)
	}

	t.Logf("获取绑定账号成功: %+v", resp)
}

// TestOAuthAPI_GetLinkedAccounts_NotAuthenticated 测试未认证获取绑定账号
func TestOAuthAPI_GetLinkedAccounts_NotAuthenticated(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	w := makeOAuthRequest(router, "GET", "/oauth/accounts", nil, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了未认证请求")
}

// TestOAuthAPI_UnlinkAccount 测试解绑账号
func TestOAuthAPI_UnlinkAccount(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 添加绑定账号
	account := &authModel.OAuthAccount{
		ID:             "oauth_123",
		UserID:         "user_123",
		Provider:       authModel.OAuthProviderGoogle,
		ProviderUserID: "google_user_123",
		Email:          "test@example.com",
		IsPrimary:      true,
	}
	mockOAuthService.linkedAccounts["user_123"] = append(mockOAuthService.linkedAccounts["user_123"], account)

	// 添加第二个账号（避免只有一个账号的错误）
	account2 := &authModel.OAuthAccount{
		ID:             "oauth_456",
		UserID:         "user_123",
		Provider:       authModel.OAuthProviderGitHub,
		ProviderUserID: "github_user_123",
		Email:          "test@github.com",
		IsPrimary:      false,
	}
	mockOAuthService.linkedAccounts["user_123"] = append(mockOAuthService.linkedAccounts["user_123"], account2)

	// 创建带用户ID的请求
	req, _ := http.NewRequest("DELETE", "/oauth/accounts/oauth_456", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	oauthAPI.UnlinkAccount(c)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "解绑成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("解绑账号成功")
}

// TestOAuthAPI_SetPrimaryAccount 测试设置主账号
func TestOAuthAPI_SetPrimaryAccount(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 添加绑定账号
	account := &authModel.OAuthAccount{
		ID:             "oauth_123",
		UserID:         "user_123",
		Provider:       authModel.OAuthProviderGoogle,
		ProviderUserID: "google_user_123",
		Email:          "test@example.com",
		IsPrimary:      false,
	}
	mockOAuthService.linkedAccounts["user_123"] = append(mockOAuthService.linkedAccounts["user_123"], account)

	// 创建带用户ID的请求
	req, _ := http.NewRequest("PUT", "/oauth/accounts/oauth_123/primary", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	oauthAPI.SetPrimaryAccount(c)

	if w.Code != http.StatusOK {
		t.Errorf("状态码错误: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "设置成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	// 验证账号已设为主账号
	accounts, _ := mockOAuthService.GetLinkedAccounts(context.Background(), "user_123")
	if !accounts[0].IsPrimary {
		t.Error("账号应该已设为主账号")
	}

	t.Logf("设置主账号成功")
}

// TestOAuthAPI_SetPrimaryAccount_MissingAccountID 测试缺失账号ID
func TestOAuthAPI_SetPrimaryAccount_MissingAccountID(t *testing.T) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	// 创建带用户ID的请求
	req, _ := http.NewRequest("PUT", "/oauth/accounts//primary", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	oauthAPI.SetPrimaryAccount(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("应该返回400错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了缺失账号ID的请求")
}

// BenchmarkGetAuthorizeURL 性能测试：获取授权URL
func BenchmarkGetAuthorizeURL(b *testing.B) {
	mockOAuthService := NewMockOAuthServiceForAPI()
	mockAuthService := NewMockAuthServiceForOAuth()
	logger := zap.NewNop()

	oauthAPI := NewOAuthAPI(mockOAuthService, mockAuthService, logger)
	router := setupTestOAuthRouter(oauthAPI)

	reqBody := map[string]interface{}{
		"redirect_uri": "http://localhost:3000/callback",
		"state":        "test_state",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/oauth/google/authorize", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
