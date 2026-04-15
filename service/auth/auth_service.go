package auth

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userPassword "Qingyu_backend/service/user"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	authModel "Qingyu_backend/models/auth"
	usersModel "Qingyu_backend/models/users"
	authRepo "Qingyu_backend/repository/interfaces/auth"
	userRepo "Qingyu_backend/repository/interfaces/user"

	"go.uber.org/zap"
)

// AuthServiceImpl Auth服务实现（整合JWT、角色、权限、会话）
type AuthServiceImpl struct {
	jwtService        JWTService
	roleService       RoleService
	permissionService PermissionService
	authRepo          authRepo.RoleRepository
	oauthRepo         authRepo.OAuthRepository        // OAuth仓储
	userRepo          userRepo.UserRepository         // 直接访问用户数据，避免 Auth ↔ UserService 循环依赖
	sessionService    SessionService                  // MVP: 会话管理（多端登录限制）
	passwordValidator *userPassword.PasswordValidator // MVP: 密码强度验证（使用 user 包统一实现）
	initialized       bool                            // 初始化标志
}

// NewAuthService 创建Auth服务
func NewAuthService(
	jwtService JWTService,
	roleService RoleService,
	permissionService PermissionService,
	authRepo authRepo.RoleRepository,
	oauthRepo authRepo.OAuthRepository,
	userRepo userRepo.UserRepository,
	sessionService SessionService,
) AuthService {
	return &AuthServiceImpl{
		jwtService:        jwtService,
		roleService:       roleService,
		permissionService: permissionService,
		authRepo:          authRepo,
		oauthRepo:         oauthRepo,
		userRepo:          userRepo,
		sessionService:    sessionService,
		passwordValidator: userPassword.NewPasswordValidator(), // 使用 user 包的统一密码验证器
	}
}

// ============ 用户认证 ============

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if req == nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeValidation, "注册请求不能为空", nil)
	}

	// 0. MVP: 验证密码强度
	if err := s.passwordValidator.ValidatePassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeValidation, "密码不符合要求: "+err.Error(), err)
	}

	// 1. 确定默认角色并直接创建用户，避免再经由 UserService 回流。
	defaultRole := req.Role
	if defaultRole == "" {
		defaultRole = "reader" // 默认为reader角色
	}

	user, err := s.createUser(ctx, req.Username, req.Email, req.Password, []string{defaultRole})
	if err != nil {
		return nil, err
	}

	// 2. 同步 auth 侧角色映射（历史兼容）。
	role, err := s.authRepo.GetRoleByName(ctx, defaultRole)
	if err != nil {
		// 如果角色不存在，使用默认角色
		role, _ = s.authRepo.GetRoleByName(ctx, "reader")
	}

	if role != nil {
		_ = s.authRepo.AssignUserRole(ctx, user.ID.Hex(), role.ID.Hex())
	}

	// 3. 生成JWT Token
	roles := user.Roles
	if len(roles) == 0 {
		roles = []string{defaultRole}
	}
	token, err := s.jwtService.GenerateToken(ctx, user.ID.Hex(), roles)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	// 4. 返回响应
	return &RegisterResponse{
		User:  convertUserToUserInfo(user, roles),
		Token: token,
	}, nil
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if req == nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeValidation, "登录请求不能为空", nil)
	}
	if req.Username == "" || req.Password == "" {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeValidation, "用户名和密码不能为空", nil)
	}

	user, err := s.authenticateUser(ctx, req.Username, req.Password, req.ClientIP)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户角色
	roleNames, err := s.resolveUserRoles(ctx, user)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "获取用户角色失败", err)
	}

	// 2.5. MVP: 强制执行多端登录限制（最多5台设备，超限自动踢出最老设备）
	if err := s.sessionService.EnforceDeviceLimit(ctx, user.ID.Hex(), 5); err != nil {
		// 记录错误但不中断登录（宽松策略）
		zap.L().Warn("设备限制执行失败，允许登录",
			zap.String("user_id", user.ID.Hex()),
			zap.Error(err),
		)
	}

	// 3. 生成JWT Token
	token, err := s.jwtService.GenerateToken(ctx, user.ID.Hex(), roleNames)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	// 3.5. MVP: 创建会话
	session, err := s.sessionService.CreateSession(ctx, user.ID.Hex())
	if err != nil {
		// 会话创建失败不影响登录（降级处理）
		zap.L().Warn("创建会话失败",
			zap.String("user_id", user.ID.Hex()),
			zap.Error(err),
		)
	}
	_ = session // 暂时不使用，后续可添加到响应中

	// 4. 返回响应
	return &LoginResponse{
		User:  convertUserToUserInfo(user, roleNames),
		Token: token,
	}, nil
}

// OAuthLogin OAuth登录
func (s *AuthServiceImpl) OAuthLogin(ctx context.Context, req *OAuthLoginRequest) (*LoginResponse, error) {
	// 1. 查找OAuth账号是否已存在
	oauthAccount, err := s.oauthRepo.FindByProviderAndProviderID(ctx, req.Provider, req.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("查询OAuth账号失败: %w", err)
	}

	// 2. 如果OAuth账号已存在，直接登录
	if oauthAccount != nil {
		user, err := s.userRepo.GetByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %w", err)
		}

		// 更新最后登录时间
		_ = s.oauthRepo.UpdateLastLogin(ctx, oauthAccount.ID.Hex())

		// 获取用户角色
		roleNames, err := s.resolveUserRoles(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("获取用户角色失败: %w", err)
		}

		// 生成JWT Token
		token, err := s.jwtService.GenerateToken(ctx, user.ID.Hex(), roleNames)
		if err != nil {
			return nil, fmt.Errorf("生成Token失败: %w", err)
		}

		return &LoginResponse{
			User:  convertUserToUserInfo(user, roleNames),
			Token: token,
		}, nil
	}

	// 3. OAuth账号不存在，创建新用户
	// 生成随机密码（OAuth用户不需要密码）
	randomPassword := generateRandomPassword(16)

	// 生成用户名（如果未提供）
	username := req.Username
	if username == "" {
		username = generateUsernameFromProvider(req.Provider, req.ProviderID)
	}

	// 创建用户
	user, err := s.createUser(ctx, username, req.Email, randomPassword, []string{"reader"})
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 4. 分配默认角色
	defaultRole := "reader"
	role, err := s.authRepo.GetRoleByName(ctx, defaultRole)
	if err == nil && role != nil {
		_ = s.authRepo.AssignUserRole(ctx, user.ID.Hex(), role.ID.Hex())
	}

	// 5. 创建OAuth账号记录
	oauthAccount = &authModel.OAuthAccount{
		UserID:         user.ID.Hex(),
		Provider:       req.Provider,
		ProviderUserID: req.ProviderID,
		Email:          req.Email,
		Username:       req.Username,
		Avatar:         req.Avatar,
		IsPrimary:      true, // 第一个OAuth账号设为主账号
		LastLoginAt:    time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := s.oauthRepo.Create(ctx, oauthAccount); err != nil {
		return nil, fmt.Errorf("创建OAuth账号失败: %w", err)
	}

	// 6. 生成JWT Token
	roles := []string{defaultRole}
	token, err := s.jwtService.GenerateToken(ctx, user.ID.Hex(), roles)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	// 7. 返回响应
	return &LoginResponse{
		User:  convertUserToUserInfo(user, roles),
		Token: token,
	}, nil
}

// Logout 用户登出
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// 将Token加入黑名单
	if err := s.jwtService.RevokeToken(ctx, token); err != nil {
		return fmt.Errorf("登出失败: %w", err)
	}

	return nil
}

// RefreshToken 刷新Token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	// 使用JWT服务刷新Token
	newToken, err := s.jwtService.RefreshToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("刷新Token失败: %w", err)
	}

	return newToken, nil
}

// ValidateToken 验证Token
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	// 使用JWT服务验证Token
	claims, err := s.jwtService.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("验证Token失败: %w", err)
	}

	return claims, nil
}

// ============ 权限管理 ============

// CheckPermission 检查权限
func (s *AuthServiceImpl) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return s.permissionService.CheckPermission(ctx, userID, permission)
}

// GetUserPermissions 获取用户权限
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return s.permissionService.GetUserPermissions(ctx, userID)
}

// HasRole 检查角色
func (s *AuthServiceImpl) HasRole(ctx context.Context, userID, role string) (bool, error) {
	return s.permissionService.HasRole(ctx, userID, role)
}

// GetUserRoles 获取用户角色
func (s *AuthServiceImpl) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := s.authRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}

// ============ 角色管理 ============

// CreateRole 创建角色
func (s *AuthServiceImpl) CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error) {
	return s.roleService.CreateRole(ctx, req)
}

// UpdateRole 更新角色
func (s *AuthServiceImpl) UpdateRole(ctx context.Context, roleID string, req *UpdateRoleRequest) error {
	return s.roleService.UpdateRole(ctx, roleID, req)
}

// DeleteRole 删除角色
func (s *AuthServiceImpl) DeleteRole(ctx context.Context, roleID string) error {
	return s.roleService.DeleteRole(ctx, roleID)
}

// AssignRole 分配角色
func (s *AuthServiceImpl) AssignRole(ctx context.Context, userID, roleID string) error {
	// 分配角色
	if err := s.authRepo.AssignUserRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	// 清除权限缓存
	if permSvc, ok := s.permissionService.(*PermissionServiceImpl); ok {
		_ = permSvc.InvalidateUserPermissionsCache(ctx, userID)
	}

	return nil
}

// RemoveRole 移除角色
func (s *AuthServiceImpl) RemoveRole(ctx context.Context, userID, roleID string) error {
	// 移除角色
	if err := s.authRepo.RemoveUserRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}

	// 清除权限缓存
	if permSvc, ok := s.permissionService.(*PermissionServiceImpl); ok {
		_ = permSvc.InvalidateUserPermissionsCache(ctx, userID)
	}

	return nil
}

// ============ 会话管理（预留） ============

// CreateSession 创建会话
func (s *AuthServiceImpl) CreateSession(ctx context.Context, userID string) (*Session, error) {
	// TODO: 实现会话管理
	return &Session{
		ID:     "session_placeholder",
		UserID: userID,
	}, nil
}

// GetSession 获取会话
func (s *AuthServiceImpl) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	// TODO: 实现会话管理
	return nil, fmt.Errorf("会话管理功能待实现")
}

// DestroySession 销毁会话
func (s *AuthServiceImpl) DestroySession(ctx context.Context, sessionID string) error {
	// TODO: 实现会话管理
	return fmt.Errorf("会话管理功能待实现")
}

// RefreshSession 刷新会话
func (s *AuthServiceImpl) RefreshSession(ctx context.Context, sessionID string) error {
	// TODO: 实现会话管理
	return fmt.Errorf("会话管理功能待实现")
}

// ============ BaseService 接口实现 ============

// Initialize 初始化认证服务
func (s *AuthServiceImpl) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// 验证依赖项
	if s.jwtService == nil {
		return fmt.Errorf("jwtService is nil")
	}
	if s.roleService == nil {
		return fmt.Errorf("roleService is nil")
	}
	if s.permissionService == nil {
		return fmt.Errorf("permissionService is nil")
	}
	if s.authRepo == nil {
		return fmt.Errorf("authRepo is nil")
	}
	if s.userRepo == nil {
		return fmt.Errorf("userRepo is nil")
	}
	if s.sessionService == nil {
		return fmt.Errorf("sessionService is nil")
	}

	// 检查Repository健康状态
	if err := s.authRepo.Health(ctx); err != nil {
		return fmt.Errorf("authRepo health check failed: %w", err)
	}

	s.initialized = true
	return nil
}

// Health 健康检查
func (s *AuthServiceImpl) Health(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("service not initialized")
	}
	return s.authRepo.Health(ctx)
}

// Close 关闭服务，清理资源
func (s *AuthServiceImpl) Close(ctx context.Context) error {
	// 认证服务暂无需要清理的资源
	// 未来如果有缓存等资源，在此处清理
	s.initialized = false
	return nil
}

// GetServiceName 获取服务名称
func (s *AuthServiceImpl) GetServiceName() string {
	return "AuthService"
}

// GetVersion 获取服务版本
func (s *AuthServiceImpl) GetVersion() string {
	return "v1.0.0"
}

// ============ 辅助函数 ============

func (s *AuthServiceImpl) createUser(ctx context.Context, username, email, password string, roles []string) (*usersModel.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeValidation, "用户名、邮箱和密码不能为空", nil)
	}

	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "检查用户名失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeBusiness, "用户名已存在", nil)
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "检查邮箱失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeBusiness, "邮箱已存在", nil)
	}

	user := &usersModel.User{
		Username: username,
		Email:    email,
		Password: password,
		Roles:    roles,
		Status:   usersModel.UserStatusActive,
	}
	if len(user.Roles) == 0 {
		user.Roles = []string{"reader"}
	}

	if err := user.SetPassword(password); err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "设置密码失败", err)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "创建用户失败", err)
	}

	return user, nil
}

func (s *AuthServiceImpl) authenticateUser(ctx context.Context, username, password, clientIP string) (*usersModel.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if userRepo.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	if !user.ValidatePassword(password) {
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeUnauthorized, "密码错误", nil)
	}

	switch user.Status {
	case usersModel.UserStatusInactive:
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeUnauthorized, "账号未激活，请先验证邮箱", nil)
	case usersModel.UserStatusBanned:
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeUnauthorized, "账号已被封禁，请联系管理员", nil)
	case usersModel.UserStatusDeleted:
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeUnauthorized, "账号已删除", nil)
	case usersModel.UserStatusActive:
	default:
		return nil, serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "未知的用户状态", nil)
	}

	if clientIP == "" {
		clientIP = "unknown"
	}
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID.Hex(), clientIP); err != nil {
		zap.L().Warn("更新最后登录时间失败", // codeql[go/log-injection]
			zap.String("user_id", user.ID.Hex()),
			zap.String("ip", clientIP),
			zap.Error(err),
		)
	}

	return user, nil
}

func (s *AuthServiceImpl) resolveUserRoles(ctx context.Context, user *usersModel.User) ([]string, error) {
	userRoles, err := s.authRepo.GetUserRoles(ctx, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		roleNames = append(roleNames, role.Name)
	}

	if len(roleNames) == 0 {
		if len(user.Roles) > 0 {
			return append([]string(nil), user.Roles...), nil
		}
		return []string{"reader"}, nil
	}

	return roleNames, nil
}

// convertUserToUserInfo 转换User为UserInfo
func convertUserToUserInfo(user *usersModel.User, roles []string) *UserInfo {
	return &UserInfo{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
		Status:   string(user.Status),
	}
}

// generateRandomPassword 生成随机密码（用于OAuth用户）
func generateRandomPassword(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// generateUsernameFromProvider 从OAuth提供商生成用户名
func generateUsernameFromProvider(provider authModel.OAuthProvider, providerID string) string {
	// 取providerID的前8位作为用户名
	shortID := providerID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	// 根据提供商生成不同的用户名格式
	switch provider {
	case authModel.OAuthProviderGoogle:
		return "google_" + shortID
	case authModel.OAuthProviderGitHub:
		return "github_" + shortID
	case authModel.OAuthProviderQQ:
		return "qq_" + shortID
	case authModel.OAuthProviderWeChat:
		return "wechat_" + shortID
	case authModel.OAuthProviderWeibo:
		return "weibo_" + shortID
	default:
		return "oauth_" + shortID
	}
}
