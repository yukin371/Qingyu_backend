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

	"Qingyu_backend/service/shared/auth"
)

// ============ Mock 实现 ============

// MockAuthServiceForAPI Mock认证服务（用于API测试）
type MockAuthServiceForAPI struct {
	users       map[string]*auth.UserInfo
	tokens      map[string]*auth.TokenClaims
	nextID      int
	shouldError bool
}

func NewMockAuthServiceForAPI() *MockAuthServiceForAPI {
	return &MockAuthServiceForAPI{
		users:  make(map[string]*auth.UserInfo),
		tokens: make(map[string]*auth.TokenClaims),
		nextID: 1,
	}
}

func (m *MockAuthServiceForAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if m.shouldError {
		return nil, &testError{msg: "注册失败"}
	}

	userID := generateUserID()
	user := &auth.UserInfo{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Roles:    []string{"reader"},
	}

	m.users[userID] = user
	token := generateToken(userID, user.Roles)

	return &auth.RegisterResponse{
		User:  user,
		Token: token,
	}, nil
}

func (m *MockAuthServiceForAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if m.shouldError {
		return nil, &testError{msg: "登录失败"}
	}

	// 查找用户
	for _, user := range m.users {
		if user.Username == req.Username {
			token := generateToken(user.ID, user.Roles)
			return &auth.LoginResponse{
				User:  user,
				Token: token,
			}, nil
		}
	}

	return nil, &testError{msg: "用户不存在或密码错误"}
}

func (m *MockAuthServiceForAPI) OAuthLogin(ctx context.Context, req *auth.OAuthLoginRequest) (*auth.LoginResponse, error) {
	return nil, &testError{msg: "OAuth登录未实现"}
}

func (m *MockAuthServiceForAPI) Logout(ctx context.Context, token string) error {
	if m.shouldError {
		return &testError{msg: "登出失败"}
	}

	delete(m.tokens, token)
	return nil
}

func (m *MockAuthServiceForAPI) RefreshToken(ctx context.Context, token string) (string, error) {
	if m.shouldError {
		return "", &testError{msg: "刷新Token失败"}
	}

	return "new_" + token, nil
}

func (m *MockAuthServiceForAPI) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	if claims, ok := m.tokens[token]; ok {
		return claims, nil
	}
	return nil, &testError{msg: "Token无效"}
}

func (m *MockAuthServiceForAPI) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return true, nil
}

func (m *MockAuthServiceForAPI) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return []string{"book.read", "book.write"}, nil
}

func (m *MockAuthServiceForAPI) HasRole(ctx context.Context, userID, role string) (bool, error) {
	return true, nil
}

func (m *MockAuthServiceForAPI) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return []string{"reader", "author"}, nil
}

func (m *MockAuthServiceForAPI) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	return nil, nil
}

func (m *MockAuthServiceForAPI) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	return nil
}

func (m *MockAuthServiceForAPI) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

func (m *MockAuthServiceForAPI) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForAPI) RemoveRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForAPI) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForAPI) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForAPI) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForAPI) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForAPI) Health(ctx context.Context) error {
	return nil
}

// ============ 测试辅助函数 ============

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func generateUserID() string {
	return "user_123"
}

func generateToken(userID string, roles []string) string {
	return "token_" + userID
}

func setupTestRouter(authAPI *AuthAPI) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 注册路由
	router.POST("/auth/register", authAPI.Register)
	router.POST("/auth/login", authAPI.Login)
	router.POST("/auth/logout", authAPI.Logout)
	router.POST("/auth/refresh", authAPI.RefreshToken)
	router.GET("/auth/permissions", authAPI.GetUserPermissions)
	router.GET("/auth/roles", authAPI.GetUserRoles)

	return router
}

func makeRequest(router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
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

// TestAuthAPI_Register 测试用户注册
func TestAuthAPI_Register(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	reqBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "StrongPass123!",
	}

	w := makeRequest(router, "POST", "/auth/register", reqBody, nil)

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

	if resp.Message != "注册成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("注册成功: %+v", resp)
}

// TestAuthAPI_Register_InvalidJSON 测试无效JSON
func TestAuthAPI_Register_InvalidJSON(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	w := makeRequest(router, "POST", "/auth/register", "invalid json", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("应该返回400错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了无效JSON")
}

// TestAuthAPI_Register_MissingFields 测试缺失字段
func TestAuthAPI_Register_MissingFields(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	reqBody := map[string]interface{}{
		"username": "testuser",
		// 缺少email和password
	}

	w := makeRequest(router, "POST", "/auth/register", reqBody, nil)

	if w.Code == http.StatusOK {
		t.Error("应该返回错误，但成功了")
	}

	t.Logf("正确拒绝了缺失字段的请求: %s", w.Body.String())
}

// TestAuthAPI_Login 测试用户登录
func TestAuthAPI_Login(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	// 先注册用户
	registerReq := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "StrongPass123!",
	}
	makeRequest(router, "POST", "/auth/register", registerReq, nil)

	// 测试登录
	loginReq := map[string]interface{}{
		"username": "testuser",
		"password": "StrongPass123!",
	}

	w := makeRequest(router, "POST", "/auth/login", loginReq, nil)

	if w.Code != http.StatusOK {
		t.Errorf("登录失败: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("响应码错误: %d", resp.Code)
	}

	data := resp.Data.(map[string]interface{})
	if data["user"] == nil {
		t.Error("响应数据中应该包含user")
	}

	t.Logf("登录成功: %+v", resp)
}

// TestAuthAPI_Login_InvalidCredentials 测试无效凭据
func TestAuthAPI_Login_InvalidCredentials(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	reqBody := map[string]interface{}{
		"username": "nonexistent",
		"password": "wrongpass",
	}

	w := makeRequest(router, "POST", "/auth/login", reqBody, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了无效凭据: %s", w.Body.String())
}

// TestAuthAPI_Logout 测试用户登出
func TestAuthAPI_Logout(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	headers := map[string]string{
		"Authorization": "Bearer test_token",
	}

	w := makeRequest(router, "POST", "/auth/logout", nil, headers)

	if w.Code != http.StatusOK {
		t.Errorf("登出失败: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Message != "登出成功" {
		t.Errorf("消息错误: %s", resp.Message)
	}

	t.Logf("登出成功")
}

// TestAuthAPI_Logout_NoToken 测试无Token登出
func TestAuthAPI_Logout_NoToken(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	w := makeRequest(router, "POST", "/auth/logout", nil, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了无Token请求")
}

// TestAuthAPI_RefreshToken 测试刷新Token
func TestAuthAPI_RefreshToken(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	headers := map[string]string{
		"Authorization": "Bearer test_token",
	}

	w := makeRequest(router, "POST", "/auth/refresh", nil, headers)

	if w.Code != http.StatusOK {
		t.Errorf("刷新Token失败: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data := resp.Data.(map[string]interface{})
	newToken := data["token"].(string)
	if newToken == "" {
		t.Error("新Token不应为空")
	}

	t.Logf("刷新Token成功: %s", newToken)
}

// TestAuthAPI_RefreshToken_NoToken 测试无Token刷新
func TestAuthAPI_RefreshToken_NoToken(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	w := makeRequest(router, "POST", "/auth/refresh", nil, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了无Token请求")
}

// TestAuthAPI_GetUserPermissions 测试获取用户权限
func TestAuthAPI_GetUserPermissions(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	_ = setupTestRouter(authAPI)

	// 创建一个带用户ID的请求
	req, _ := http.NewRequest("GET", "/auth/permissions", nil)
	w := httptest.NewRecorder()

	// 设置gin上下文
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	// 调用handler
	authAPI.GetUserPermissions(c)

	if w.Code != http.StatusOK {
		t.Errorf("获取权限失败: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("响应码错误: %d", resp.Code)
	}

	t.Logf("获取权限成功: %+v", resp)
}

// TestAuthAPI_GetUserPermissions_NotAuthenticated 测试未认证获取权限
func TestAuthAPI_GetUserPermissions_NotAuthenticated(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	w := makeRequest(router, "GET", "/auth/permissions", nil, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了未认证请求")
}

// TestAuthAPI_GetUserRoles 测试获取用户角色
func TestAuthAPI_GetUserRoles(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	// 创建一个带用户ID的请求
	req, _ := http.NewRequest("GET", "/auth/roles", nil)
	w := httptest.NewRecorder()

	// 设置gin上下文
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user_123")

	// 调用handler
	authAPI.GetUserRoles(c)

	if w.Code != http.StatusOK {
		t.Errorf("获取角色失败: %d, body: %s", w.Code, w.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("响应码错误: %d", resp.Code)
	}

	t.Logf("获取角色成功: %+v", resp)
}

// TestAuthAPI_GetUserRoles_NotAuthenticated 测试未认证获取角色
func TestAuthAPI_GetUserRoles_NotAuthenticated(t *testing.T) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	w := makeRequest(router, "GET", "/auth/roles", nil, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("应该返回401错误，实际: %d", w.Code)
	}

	t.Logf("正确拒绝了未认证请求")
}

// BenchmarkRegister 性能测试：注册接口
func BenchmarkRegister(b *testing.B) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	reqBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "StrongPass123!",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkLogin 性能测试：登录接口
func BenchmarkLogin(b *testing.B) {
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)

	// 预先注册用户
	registerReq := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "StrongPass123!",
	}
	bodyBytes, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 准备登录请求
	loginReq := map[string]interface{}{
		"username": "testuser",
		"password": "StrongPass123!",
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// ============ 集成测试辅助函数 ============

// SetupTestAuthAPIWithLogger 设置带logger的测试API
func SetupTestAuthAPIWithLogger() (*AuthAPI, *gin.Engine) {
	logger := zap.NewNop()
	mockService := NewMockAuthServiceForAPI()
	authAPI := NewAuthAPI(mockService)
	router := setupTestRouter(authAPI)
	return authAPI, router
}
