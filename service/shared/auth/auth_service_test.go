package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	authModel "Qingyu_backend/models/auth"
	usersModel "Qingyu_backend/models/users"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/config"
)

// ============ Mock 实现 ============

// MockUserService Mock用户服务
type MockUserService struct {
	users       map[string]*usersModel.User
	nextID      int
	shouldError bool
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:  make(map[string]*usersModel.User),
		nextID: 1,
	}
}

func (m *MockUserService) CreateUser(ctx context.Context, req *userServiceInterface.CreateUserRequest) (*userServiceInterface.CreateUserResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf("user service error")
	}

	user := &usersModel.User{
		ID:       fmt.Sprintf("user_%d", m.nextID),
		Username: req.Username,
		Email:    req.Email,
	}
	m.nextID++
	m.users[user.ID] = user

	return &userServiceInterface.CreateUserResponse{
		User: user,
	}, nil
}

func (m *MockUserService) LoginUser(ctx context.Context, req *userServiceInterface.LoginUserRequest) (*userServiceInterface.LoginUserResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf("login failed")
	}

	// 查找用户
	for _, user := range m.users {
		if user.Username == req.Username {
			return &userServiceInterface.LoginUserResponse{
				User: user,
			}, nil
		}
	}

	return nil, fmt.Errorf("用户不存在或密码错误")
}

func (m *MockUserService) RegisterUser(ctx context.Context, req *userServiceInterface.RegisterUserRequest) (*userServiceInterface.RegisterUserResponse, error) {
	return nil, nil
}

func (m *MockUserService) LogoutUser(ctx context.Context, req *userServiceInterface.LogoutUserRequest) (*userServiceInterface.LogoutUserResponse, error) {
	return nil, nil
}

func (m *MockUserService) ValidateToken(ctx context.Context, req *userServiceInterface.ValidateTokenRequest) (*userServiceInterface.ValidateTokenResponse, error) {
	return nil, nil
}

func (m *MockUserService) UpdateLastLogin(ctx context.Context, req *userServiceInterface.UpdateLastLoginRequest) (*userServiceInterface.UpdateLastLoginResponse, error) {
	return nil, nil
}

func (m *MockUserService) UpdatePassword(ctx context.Context, req *userServiceInterface.UpdatePasswordRequest) (*userServiceInterface.UpdatePasswordResponse, error) {
	return nil, nil
}

func (m *MockUserService) RequestPasswordReset(ctx context.Context, req *userServiceInterface.RequestPasswordResetRequest) (*userServiceInterface.RequestPasswordResetResponse, error) {
	return nil, nil
}

func (m *MockUserService) ConfirmPasswordReset(ctx context.Context, req *userServiceInterface.ConfirmPasswordResetRequest) (*userServiceInterface.ConfirmPasswordResetResponse, error) {
	return nil, nil
}

func (m *MockUserService) SendEmailVerification(ctx context.Context, req *userServiceInterface.SendEmailVerificationRequest) (*userServiceInterface.SendEmailVerificationResponse, error) {
	return nil, nil
}

func (m *MockUserService) AssignRole(ctx context.Context, req *userServiceInterface.AssignRoleRequest) (*userServiceInterface.AssignRoleResponse, error) {
	return nil, nil
}

func (m *MockUserService) RemoveRole(ctx context.Context, req *userServiceInterface.RemoveRoleRequest) (*userServiceInterface.RemoveRoleResponse, error) {
	return nil, nil
}

func (m *MockUserService) GetUserRoles(ctx context.Context, req *userServiceInterface.GetUserRolesRequest) (*userServiceInterface.GetUserRolesResponse, error) {
	return nil, nil
}

func (m *MockUserService) GetUserPermissions(ctx context.Context, req *userServiceInterface.GetUserPermissionsRequest) (*userServiceInterface.GetUserPermissionsResponse, error) {
	return nil, nil
}

func (m *MockUserService) Initialize(ctx context.Context) error {
	return nil
}

func (m *MockUserService) Health(ctx context.Context) error {
	return nil
}

func (m *MockUserService) Close(ctx context.Context) error {
	return nil
}

func (m *MockUserService) GetServiceName() string {
	return "MockUserService"
}

func (m *MockUserService) GetVersion() string {
	return "v1.0.0"
}

func (m *MockUserService) GetUser(ctx context.Context, req *userServiceInterface.GetUserRequest) (*userServiceInterface.GetUserResponse, error) {
	user, ok := m.users[req.ID]
	if !ok {
		return nil, fmt.Errorf("用户不存在")
	}
	return &userServiceInterface.GetUserResponse{User: user}, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, req *userServiceInterface.UpdateUserRequest) (*userServiceInterface.UpdateUserResponse, error) {
	return nil, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, req *userServiceInterface.DeleteUserRequest) (*userServiceInterface.DeleteUserResponse, error) {
	return nil, nil
}

func (m *MockUserService) ListUsers(ctx context.Context, req *userServiceInterface.ListUsersRequest) (*userServiceInterface.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockUserService) ResetPassword(ctx context.Context, req *userServiceInterface.ResetPasswordRequest) (*userServiceInterface.ResetPasswordResponse, error) {
	return nil, nil
}

func (m *MockUserService) VerifyEmail(ctx context.Context, req *userServiceInterface.VerifyEmailRequest) (*userServiceInterface.VerifyEmailResponse, error) {
	return nil, nil
}

// MockOAuthRepository Mock OAuth仓储
type MockOAuthRepository struct {
	accounts map[string]*authModel.OAuthAccount // key: provider_providerID
	nextID   int
}

func NewMockOAuthRepository() *MockOAuthRepository {
	return &MockOAuthRepository{
		accounts: make(map[string]*authModel.OAuthAccount),
		nextID:   1,
	}
}

func (m *MockOAuthRepository) Create(ctx context.Context, account *authModel.OAuthAccount) error {
	account.ID = fmt.Sprintf("oauth_%d", m.nextID)
	m.nextID++
	key := fmt.Sprintf("%s_%s", account.Provider, account.ProviderUserID)
	m.accounts[key] = account
	return nil
}

func (m *MockOAuthRepository) FindByProviderAndProviderID(ctx context.Context, provider authModel.OAuthProvider, providerID string) (*authModel.OAuthAccount, error) {
	key := fmt.Sprintf("%s_%s", provider, providerID)
	return m.accounts[key], nil
}

func (m *MockOAuthRepository) FindByUserID(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	result := make([]*authModel.OAuthAccount, 0)
	for _, account := range m.accounts {
		if account.UserID == userID {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *MockOAuthRepository) Update(ctx context.Context, account *authModel.OAuthAccount) error {
	key := fmt.Sprintf("%s_%s", account.Provider, account.ProviderUserID)
	m.accounts[key] = account
	return nil
}

func (m *MockOAuthRepository) Delete(ctx context.Context, id string) error {
	for key, account := range m.accounts {
		if account.ID == id {
			delete(m.accounts, key)
			return nil
		}
	}
	return nil
}

func (m *MockOAuthRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return nil
}

func (m *MockOAuthRepository) UpdateTokens(ctx context.Context, id string, accessToken, refreshToken string, expiresAt primitive.DateTime) error {
	return nil
}

func (m *MockOAuthRepository) SetPrimaryAccount(ctx context.Context, userID string, accountID string) error {
	return nil
}

func (m *MockOAuthRepository) GetPrimaryAccount(ctx context.Context, userID string) (*authModel.OAuthAccount, error) {
	return nil, nil
}

func (m *MockOAuthRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *MockOAuthRepository) CreateSession(ctx context.Context, session *authModel.OAuthSession) error {
	return nil
}

func (m *MockOAuthRepository) FindSessionByID(ctx context.Context, id string) (*authModel.OAuthSession, error) {
	return nil, nil
}

func (m *MockOAuthRepository) FindSessionByState(ctx context.Context, state string) (*authModel.OAuthSession, error) {
	return nil, nil
}

func (m *MockOAuthRepository) DeleteSession(ctx context.Context, id string) error {
	return nil
}

func (m *MockOAuthRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockOAuthRepository) FindByID(ctx context.Context, id string) (*authModel.OAuthAccount, error) {
	for _, account := range m.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return nil, nil
}

// MockSessionService Mock会话服务
type MockSessionService struct {
	sessions map[string]*Session
	nextID   int
}

func NewMockSessionService() *MockSessionService {
	return &MockSessionService{
		sessions: make(map[string]*Session),
		nextID:   1,
	}
}

func (m *MockSessionService) CreateSession(ctx context.Context, userID string) (*Session, error) {
	session := &Session{
		ID:        fmt.Sprintf("session_%d", m.nextID),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	m.nextID++
	m.sessions[session.ID] = session
	return session, nil
}

func (m *MockSessionService) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return m.sessions[sessionID], nil
}

func (m *MockSessionService) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	return nil
}

func (m *MockSessionService) DestroySession(ctx context.Context, sessionID string) error {
	delete(m.sessions, sessionID)
	return nil
}

func (m *MockSessionService) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockSessionService) CheckDeviceLimit(ctx context.Context, userID string, maxDevices int) error {
	return nil
}

func (m *MockSessionService) EnforceDeviceLimit(ctx context.Context, userID string, maxDevices int) error {
	return nil
}

func (m *MockSessionService) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	return nil, nil
}

func (m *MockSessionService) DestroyUserSessions(ctx context.Context, userID string) error {
	return nil
}

// ============ 测试辅助函数 ============

func setupTestAuthService() (*AuthServiceImpl, *MockAuthRepository, *MockUserService, *MockOAuthRepository, *MockSessionService) {
	authRepo := NewMockAuthRepository()
	userService := NewMockUserService()
	oauthRepo := NewMockOAuthRepository()
	sessionService := NewMockSessionService()
	redisClient := NewMockRedisClient()
	cacheClient := NewMockCacheClient()

	jwtConfig := &config.JWTConfigEnhanced{
		SecretKey:       "test-secret-key-for-unit-testing",
		Issuer:          "qingyu-test",
		Expiration:      1 * time.Hour,
		RefreshDuration: 7 * 24 * time.Hour,
	}

	jwtService := NewJWTService(jwtConfig, redisClient)
	roleService := NewRoleService(authRepo)
	permissionService := NewPermissionService(authRepo, cacheClient)

	authService := NewAuthService(
		jwtService,
		roleService,
		permissionService,
		authRepo,
		oauthRepo,
		userService,
		sessionService,
	).(*AuthServiceImpl)

	// 初始化服务
	ctx := context.Background()
	_ = authService.Initialize(ctx)

	// 创建默认角色
	readerRole := &authModel.Role{
		ID:          "role_reader",
		Name:        "reader",
		Description: "读者角色",
		Permissions: []string{"book.read"},
		IsSystem:    true,
	}
	authRepo.roles["role_reader"] = readerRole

	return authService, authRepo, userService, oauthRepo, sessionService
}

// ============ 测试用例 ============

// TestRegister 测试用户注册
func TestRegister(t *testing.T) {
	authService, _, userService, _, _ := setupTestAuthService()
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "StrongPass123!",
		Role:     "reader",
	}

	resp, err := authService.Register(ctx, req)
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	if resp.User.Username != "testuser" {
		t.Errorf("用户名错误: %s", resp.User.Username)
	}
	if resp.Token == "" {
		t.Error("Token不应为空")
	}

	// 验证用户已创建
	if len(userService.users) != 1 {
		t.Error("用户未正确创建")
	}

	t.Logf("注册成功: %+v", resp.User)
}

// TestRegister_WeakPassword 测试弱密码注册
func TestRegister_WeakPassword(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "123", // 弱密码
	}

	_, err := authService.Register(ctx, req)
	if err == nil {
		t.Fatal("弱密码应该被拒绝，但成功了")
	}

	t.Logf("正确拒绝了弱密码: %v", err)
}

// TestRegister_NoPassword 测试无密码注册
func TestRegister_NoPassword(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "",
	}

	_, err := authService.Register(ctx, req)
	if err == nil {
		t.Fatal("空密码应该被拒绝，但成功了")
	}

	t.Logf("正确拒绝了空密码: %v", err)
}

// TestLogin 测试用户登录
func TestLogin(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 先注册用户
	registerReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}
	authService.Register(ctx, registerReq)

	// 测试登录
	loginReq := &LoginRequest{
		Username: "testuser",
		Password: "StrongPass123!",
	}

	resp, err := authService.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	if resp.User.Username != "testuser" {
		t.Errorf("用户名错误: %s", resp.User.Username)
	}
	if resp.Token == "" {
		t.Error("Token不应为空")
	}

	// 验证会话已创建
	if len(authRepo.userRoles) == 0 {
		t.Error("用户角色未分配")
	}

	t.Logf("登录成功: %+v", resp.User)
}

// TestLogin_InvalidCredentials 测试无效凭据登录
func TestLogin_InvalidCredentials(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	loginReq := &LoginRequest{
		Username: "nonexistent",
		Password: "wrongpass",
	}

	_, err := authService.Login(ctx, loginReq)
	if err == nil {
		t.Fatal("无效凭据应该被拒绝，但成功了")
	}

	t.Logf("正确拒绝了无效凭据: %v", err)
}

// TestOAuthLogin_ExistingUser 测试OAuth登录（已存在用户）
func TestOAuthLogin_ExistingUser(t *testing.T) {
	authService, authRepo, userService, oauthRepo, _ := setupTestAuthService()
	ctx := context.Background()

	// 先创建OAuth账号
	user := &usersModel.User{
		ID:       "user_existing",
		Username: "existinguser",
		Email:    "existing@example.com",
	}
	userService.users[user.ID] = user

	oauthAccount := &authModel.OAuthAccount{
		ID:             "oauth_1",
		UserID:         user.ID,
		Provider:       authModel.OAuthProviderGoogle,
		ProviderUserID: "google123",
		Email:          user.Email,
		IsPrimary:      true,
		LastLoginAt:    time.Now(),
	}
	key := fmt.Sprintf("%s_%s", oauthAccount.Provider, oauthAccount.ProviderUserID)
	oauthRepo.accounts[key] = oauthAccount

	// 分配角色
	authRepo.userRoles[user.ID] = []string{"role_reader"}

	// OAuth登录
	req := &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google123",
		Email:      "existing@example.com",
		Name:       "Existing User",
	}

	resp, err := authService.OAuthLogin(ctx, req)
	if err != nil {
		t.Fatalf("OAuth登录失败: %v", err)
	}

	if resp.User.ID != "user_existing" {
		t.Errorf("用户ID错误: %s", resp.User.ID)
	}
	if resp.Token == "" {
		t.Error("Token不应为空")
	}

	t.Logf("OAuth登录成功（已存在用户）: %+v", resp.User)
}

// TestOAuthLogin_NewUser 测试OAuth登录（新用户）
func TestOAuthLogin_NewUser(t *testing.T) {
	authService, _, _, oauthRepo, _ := setupTestAuthService()
	ctx := context.Background()

	req := &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google_new_123",
		Email:      "newuser@example.com",
		Name:       "New User",
	}

	resp, err := authService.OAuthLogin(ctx, req)
	if err != nil {
		t.Fatalf("OAuth登录失败: %v", err)
	}

	if resp.User.Email != "newuser@example.com" {
		t.Errorf("邮箱错误: %s", resp.User.Email)
	}
	if resp.Token == "" {
		t.Error("Token不应为空")
	}

	// 验证OAuth账号已创建
	key := fmt.Sprintf("%s_%s", req.Provider, req.ProviderID)
	if oauthRepo.accounts[key] == nil {
		t.Error("OAuth账号未创建")
	}

	t.Logf("OAuth登录成功（新用户）: %+v", resp.User)
}

// TestOAuthLogin_GitHub 测试GitHub OAuth登录
func TestOAuthLogin_GitHub(t *testing.T) {
	authService, _, _, oauthRepo, _ := setupTestAuthService()
	ctx := context.Background()

	req := &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGitHub,
		ProviderID: "github_456",
		Email:      "githubuser@example.com",
		Name:       "GitHub User",
		Username:   "githubuser",
	}

	resp, err := authService.OAuthLogin(ctx, req)
	if err != nil {
		t.Fatalf("GitHub OAuth登录失败: %v", err)
	}

	// GitHub用户应该使用提供的用户名
	if resp.User.Username != "githubuser" {
		t.Errorf("用户名错误: %s", resp.User.Username)
	}

	// 验证OAuth账号已创建
	key := fmt.Sprintf("%s_%s", req.Provider, req.ProviderID)
	if oauthRepo.accounts[key] == nil {
		t.Error("OAuth账号未创建")
	}

	t.Logf("GitHub OAuth登录成功: %+v", resp.User)
}

// TestLogout 测试用户登出
func TestLogout(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 先注册获取token
	registerReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}
	resp, _ := authService.Register(ctx, registerReq)

	// 登出
	err := authService.Logout(ctx, resp.Token)
	if err != nil {
		t.Fatalf("登出失败: %v", err)
	}

	// 验证token已被吊销
	revoked, _ := authService.jwtService.(*JWTServiceImpl).IsTokenRevoked(ctx, resp.Token)
	if !revoked {
		t.Error("Token应该已被吊销")
	}

	t.Logf("登出成功，Token已吊销")
}

// TestRefreshToken 测试刷新Token
func TestRefreshToken_AuthService(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 先注册获取token
	registerReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}
	resp, _ := authService.Register(ctx, registerReq)

	oldToken := resp.Token

	// 等待一小段时间确保token不同
	time.Sleep(1100 * time.Millisecond)

	// 刷新token
	newToken, err := authService.RefreshToken(ctx, oldToken)
	if err != nil {
		t.Fatalf("刷新Token失败: %v", err)
	}

	if newToken == oldToken {
		t.Error("新Token应该与旧Token不同")
	}

	t.Logf("刷新Token成功: old=%s..., new=%s...", oldToken[:20], newToken[:20])
}

// TestValidateToken 测试验证Token
func TestValidateToken_AuthService(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 先注册获取token
	registerReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}
	resp, _ := authService.Register(ctx, registerReq)

	// 验证token
	claims, err := authService.ValidateToken(ctx, resp.Token)
	if err != nil {
		t.Fatalf("验证Token失败: %v", err)
	}

	if claims.UserID != resp.User.ID {
		t.Errorf("UserID错误: %s", claims.UserID)
	}

	t.Logf("验证Token成功: %+v", claims)
}

// TestCheckPermission 测试权限检查
func TestCheckPermission_AuthService(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建带权限的角色
	editorRole := &authModel.Role{
		ID:          "role_editor",
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.read", "book.write"},
	}
	authRepo.roles["role_editor"] = editorRole

	// 创建用户并分配角色
	userID := "user_1"
	authRepo.userRoles[userID] = []string{"role_editor"}

	// 测试有权限
	hasRead, err := authService.CheckPermission(ctx, userID, "book.read")
	if err != nil {
		t.Fatalf("检查权限失败: %v", err)
	}
	if !hasRead {
		t.Error("用户应该有book.read权限")
	}

	// 测试无权限
	hasDelete, _ := authService.CheckPermission(ctx, userID, "book.delete")
	if hasDelete {
		t.Error("用户不应该有book.delete权限")
	}

	t.Logf("权限检查测试通过")
}

// TestGetUserPermissions 测试获取用户权限
func TestGetUserPermissions_AuthService(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建多个角色
	role1 := &authModel.Role{
		ID:          "role_1",
		Name:        "reader",
		Permissions: []string{"book.read"},
	}
	role2 := &authModel.Role{
		ID:          "role_2",
		Name:        "author",
		Permissions: []string{"book.write"},
	}
	authRepo.roles["role_1"] = role1
	authRepo.roles["role_2"] = role2

	// 分配多个角色
	userID := "user_1"
	authRepo.userRoles[userID] = []string{"role_1", "role_2"}

	// 获取权限
	perms, err := authService.GetUserPermissions(ctx, userID)
	if err != nil {
		t.Fatalf("获取权限失败: %v", err)
	}

	if len(perms) != 2 {
		t.Errorf("权限数量错误: %d", len(perms))
	}

	t.Logf("用户权限: %v", perms)
}

// TestHasRole 测试检查角色
func TestHasRole_AuthService(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_admin",
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"*"},
	}
	authRepo.roles["role_admin"] = role

	// 分配角色
	userID := "user_1"
	authRepo.userRoles[userID] = []string{"role_admin"}

	// 测试有角色
	hasAdmin, err := authService.HasRole(ctx, userID, "admin")
	if err != nil {
		t.Fatalf("检查角色失败: %v", err)
	}
	if !hasAdmin {
		t.Error("用户应该有admin角色")
	}

	// 测试无角色
	hasEditor, _ := authService.HasRole(ctx, userID, "editor")
	if hasEditor {
		t.Error("用户不应该有editor角色")
	}

	t.Logf("角色检查测试通过")
}

// TestGetUserRoles 测试获取用户角色
func TestGetUserRoles(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建多个角色
	role1 := &authModel.Role{
		ID:          "role_1",
		Name:        "reader",
		Description: "读者",
		Permissions: []string{"book.read"},
	}
	role2 := &authModel.Role{
		ID:          "role_2",
		Name:        "author",
		Description: "作者",
		Permissions: []string{"book.write"},
	}
	authRepo.roles["role_1"] = role1
	authRepo.roles["role_2"] = role2

	// 分配多个角色
	userID := "user_1"
	authRepo.userRoles[userID] = []string{"role_1", "role_2"}

	// 获取角色
	roles, err := authService.GetUserRoles(ctx, userID)
	if err != nil {
		t.Fatalf("获取角色失败: %v", err)
	}

	if len(roles) != 2 {
		t.Errorf("角色数量错误: %d", len(roles))
	}

	t.Logf("用户角色: %v", roles)
}

// TestAssignRole 测试分配角色
func TestAssignRole(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_editor",
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.write"},
	}
	authRepo.roles["role_editor"] = role

	// 分配角色
	userID := "user_1"
	roleID := "role_editor"

	err := authService.AssignRole(ctx, userID, roleID)
	if err != nil {
		t.Fatalf("分配角色失败: %v", err)
	}

	// 验证角色已分配
	roles := authRepo.userRoles[userID]
	if len(roles) != 1 || roles[0] != roleID {
		t.Error("角色未正确分配")
	}

	t.Logf("分配角色成功: %s -> %s", userID, roleID)
}

// TestRemoveRole 测试移除角色
func TestRemoveRole(t *testing.T) {
	authService, authRepo, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_editor",
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.write"},
	}
	authRepo.roles["role_editor"] = role

	// 先分配角色
	userID := "user_1"
	roleID := "role_editor"
	authRepo.userRoles[userID] = []string{roleID}

	// 移除角色
	err := authService.RemoveRole(ctx, userID, roleID)
	if err != nil {
		t.Fatalf("移除角色失败: %v", err)
	}

	// 验证角色已移除
	roles := authRepo.userRoles[userID]
	if len(roles) != 0 {
		t.Error("角色未正确移除")
	}

	t.Logf("移除角色成功: %s -> %s", userID, roleID)
}

// TestHealth 测试健康检查
func TestHealth(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	err := authService.Health(ctx)
	if err != nil {
		t.Fatalf("健康检查失败: %v", err)
	}

	t.Logf("健康检查通过")
}

// TestServiceName 测试服务名称
func TestServiceName(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()

	name := authService.GetServiceName()
	if name != "AuthService" {
		t.Errorf("服务名称错误: %s", name)
	}

	t.Logf("服务名称: %s", name)
}

// TestGetVersion 测试服务版本
func TestGetVersion(t *testing.T) {
	authService, _, _, _, _ := setupTestAuthService()

	version := authService.GetVersion()
	if version == "" {
		t.Error("版本不应为空")
	}

	t.Logf("服务版本: %s", version)
}
