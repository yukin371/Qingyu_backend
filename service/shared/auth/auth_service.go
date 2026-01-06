package auth

import (
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"context"
	"fmt"

	usersModel "Qingyu_backend/models/users"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"

	"go.uber.org/zap"
)

// AuthServiceImpl Auth服务实现（整合JWT、角色、权限、会话）
type AuthServiceImpl struct {
	jwtService        JWTService
	roleService       RoleService
	permissionService PermissionService
	authRepo          sharedRepo.AuthRepository
	userService       userServiceInterface.UserService // 依赖User服务
	sessionService    SessionService                   // MVP: 会话管理（多端登录限制）
	passwordValidator *PasswordValidator               // MVP: 密码强度验证
	initialized       bool                             // 初始化标志
}

// NewAuthService 创建Auth服务
func NewAuthService(
	jwtService JWTService,
	roleService RoleService,
	permissionService PermissionService,
	authRepo sharedRepo.AuthRepository,
	userService userServiceInterface.UserService,
	sessionService SessionService,
) AuthService {
	return &AuthServiceImpl{
		jwtService:        jwtService,
		roleService:       roleService,
		permissionService: permissionService,
		authRepo:          authRepo,
		userService:       userService,
		sessionService:    sessionService,
		passwordValidator: NewPasswordValidator(), // MVP: 使用默认密码验证规则
	}
}

// ============ 用户认证 ============

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 0. MVP: 验证密码强度
	if err := s.passwordValidator.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("密码不符合要求: %w", err)
	}

	// 1. 调用User服务创建用户
	createUserReq := &userServiceInterface.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	userResp, err := s.userService.CreateUser(ctx, createUserReq)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 2. 分配默认角色
	defaultRole := req.Role
	if defaultRole == "" {
		defaultRole = "reader" // 默认为reader角色
	}

	// 查找角色
	role, err := s.authRepo.GetRoleByName(ctx, defaultRole)
	if err != nil {
		// 如果角色不存在，使用默认角色
		role, _ = s.authRepo.GetRoleByName(ctx, "reader")
	}

	if role != nil {
		_ = s.authRepo.AssignUserRole(ctx, userResp.User.ID, role.ID)
	}

	// 3. 生成JWT Token
	roles := []string{defaultRole}
	token, err := s.jwtService.GenerateToken(ctx, userResp.User.ID, roles)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	// 4. 返回响应
	return &RegisterResponse{
		User: &UserInfo{
			ID:       userResp.User.ID,
			Username: userResp.User.Username,
			Email:    userResp.User.Email,
			Roles:    roles,
		},
		Token: token,
	}, nil
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 调用User服务登录
	loginReq := &userServiceInterface.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	}

	loginResp, err := s.userService.LoginUser(ctx, loginReq)
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	// 2. 获取用户角色
	userRoles, err := s.authRepo.GetUserRoles(ctx, loginResp.User.ID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	roleNames := make([]string, len(userRoles))
	for i, role := range userRoles {
		roleNames[i] = role.Name
	}

	// 如果没有角色，分配默认角色
	if len(roleNames) == 0 {
		roleNames = []string{"reader"}
	}

	// 2.5. MVP: 强制执行多端登录限制（最多5台设备，超限自动踢出最老设备）
	if err := s.sessionService.EnforceDeviceLimit(ctx, loginResp.User.ID, 5); err != nil {
		// 记录错误但不中断登录（宽松策略）
		zap.L().Warn("设备限制执行失败，允许登录",
			zap.String("user_id", loginResp.User.ID),
			zap.Error(err),
		)
	}

	// 3. 生成JWT Token
	token, err := s.jwtService.GenerateToken(ctx, loginResp.User.ID, roleNames)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	// 3.5. MVP: 创建会话
	session, err := s.sessionService.CreateSession(ctx, loginResp.User.ID)
	if err != nil {
		// 会话创建失败不影响登录（降级处理）
		zap.L().Warn("创建会话失败",
			zap.String("user_id", loginResp.User.ID),
			zap.Error(err),
		)
	}
	_ = session // 暂时不使用，后续可添加到响应中

	// 4. 返回响应
	return &LoginResponse{
		User: &UserInfo{
			ID:       loginResp.User.ID,
			Username: loginResp.User.Username,
			Email:    loginResp.User.Email,
			Roles:    roleNames,
		},
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
	if s.userService == nil {
		return fmt.Errorf("userService is nil")
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

// convertUserToUserInfo 转换User为UserInfo
func convertUserToUserInfo(user *usersModel.User, roles []string) *UserInfo {
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}
}
