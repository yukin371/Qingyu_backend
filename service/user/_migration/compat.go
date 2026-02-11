package user

import (
	"context"

	useriface "Qingyu_backend/service/interfaces/user"
)

// ============================================================================
// 兼容层 - 向后兼容支持
// ============================================================================

// UserServiceAdapter 旧 UserService 接口的适配器
// 将旧的 UserService 方法调用委托给新的 Port 接口
type UserServiceAdapter struct {
	managementPort useriface.UserManagementPort
	authPort       useriface.UserAuthPort
	passwordPort   useriface.PasswordManagementPort
	emailPort      useriface.EmailManagementPort
	permissionPort useriface.UserPermissionPort
	statusPort     useriface.UserStatusPort
}

// ============================================================================
// BaseService 接口实现 - 委托给 managementPort
// ============================================================================

// Initialize 初始化服务
func (a *UserServiceAdapter) Initialize(ctx context.Context) error {
	return a.managementPort.Initialize(ctx)
}

// Health 健康检查
func (a *UserServiceAdapter) Health(ctx context.Context) error {
	return a.managementPort.Health(ctx)
}

// Close 关闭服务
func (a *UserServiceAdapter) Close(ctx context.Context) error {
	// 关闭所有 Port
	if err := a.managementPort.Close(ctx); err != nil {
		return err
	}
	if err := a.authPort.Close(ctx); err != nil {
		return err
	}
	if err := a.passwordPort.Close(ctx); err != nil {
		return err
	}
	if err := a.emailPort.Close(ctx); err != nil {
		return err
	}
	if err := a.permissionPort.Close(ctx); err != nil {
		return err
	}
	if err := a.statusPort.Close(ctx); err != nil {
		return err
	}
	return nil
}

// GetServiceName 获取服务名称
func (a *UserServiceAdapter) GetServiceName() string {
	return a.managementPort.GetServiceName()
}

// GetVersion 获取服务版本
func (a *UserServiceAdapter) GetVersion() string {
	return a.managementPort.GetVersion()
}

// NewUserServiceAdapter 创建新的适配器
func NewUserServiceAdapter(
	managementPort useriface.UserManagementPort,
	authPort useriface.UserAuthPort,
	passwordPort useriface.PasswordManagementPort,
	emailPort useriface.EmailManagementPort,
	permissionPort useriface.UserPermissionPort,
	statusPort useriface.UserStatusPort,
) *UserServiceAdapter {
	return &UserServiceAdapter{
		managementPort: managementPort,
		authPort:       authPort,
		passwordPort:   passwordPort,
		emailPort:      emailPort,
		permissionPort: permissionPort,
		statusPort:     statusPort,
	}
}

// ============================================================================
// 用户管理方法 - 委托给 UserManagementPort
// ============================================================================

// CreateUser 委托给 UserManagementPort
func (a *UserServiceAdapter) CreateUser(ctx context.Context, req *useriface.CreateUserRequest) (*useriface.CreateUserResponse, error) {
	return a.managementPort.CreateUser(ctx, req)
}

// GetUser 委托给 UserManagementPort
func (a *UserServiceAdapter) GetUser(ctx context.Context, req *useriface.GetUserRequest) (*useriface.GetUserResponse, error) {
	return a.managementPort.GetUser(ctx, req)
}

// UpdateUser 委托给 UserManagementPort
func (a *UserServiceAdapter) UpdateUser(ctx context.Context, req *useriface.UpdateUserRequest) (*useriface.UpdateUserResponse, error) {
	return a.managementPort.UpdateUser(ctx, req)
}

// DeleteUser 委托给 UserManagementPort
func (a *UserServiceAdapter) DeleteUser(ctx context.Context, req *useriface.DeleteUserRequest) (*useriface.DeleteUserResponse, error) {
	return a.managementPort.DeleteUser(ctx, req)
}

// ListUsers 委托给 UserManagementPort
func (a *UserServiceAdapter) ListUsers(ctx context.Context, req *useriface.ListUsersRequest) (*useriface.ListUsersResponse, error) {
	return a.managementPort.ListUsers(ctx, req)
}

// ============================================================================
// 用户认证方法 - 委托给 UserAuthPort
// ============================================================================

// RegisterUser 委托给 UserAuthPort
func (a *UserServiceAdapter) RegisterUser(ctx context.Context, req *useriface.RegisterUserRequest) (*useriface.RegisterUserResponse, error) {
	return a.authPort.RegisterUser(ctx, req)
}

// LoginUser 委托给 UserAuthPort
func (a *UserServiceAdapter) LoginUser(ctx context.Context, req *useriface.LoginUserRequest) (*useriface.LoginUserResponse, error) {
	return a.authPort.LoginUser(ctx, req)
}

// LogoutUser 委托给 UserAuthPort
func (a *UserServiceAdapter) LogoutUser(ctx context.Context, req *useriface.LogoutUserRequest) (*useriface.LogoutUserResponse, error) {
	return a.authPort.LogoutUser(ctx, req)
}

// ValidateToken 委托给 UserAuthPort
func (a *UserServiceAdapter) ValidateToken(ctx context.Context, req *useriface.ValidateTokenRequest) (*useriface.ValidateTokenResponse, error) {
	return a.authPort.ValidateToken(ctx, req)
}

// ============================================================================
// 密码管理方法 - 委托给 PasswordManagementPort
// ============================================================================

// UpdatePassword 委托给 PasswordManagementPort
func (a *UserServiceAdapter) UpdatePassword(ctx context.Context, req *useriface.UpdatePasswordRequest) (*useriface.UpdatePasswordResponse, error) {
	return a.passwordPort.UpdatePassword(ctx, req)
}

// ResetPassword 委托给 PasswordManagementPort
func (a *UserServiceAdapter) ResetPassword(ctx context.Context, req *useriface.ResetPasswordRequest) (*useriface.ResetPasswordResponse, error) {
	return a.passwordPort.ResetPassword(ctx, req)
}

// VerifyPassword 委托给 PasswordManagementPort
func (a *UserServiceAdapter) VerifyPassword(ctx context.Context, userID string, password string) error {
	return a.passwordPort.VerifyPassword(ctx, userID, password)
}

// RequestPasswordReset 委托给 PasswordManagementPort
func (a *UserServiceAdapter) RequestPasswordReset(ctx context.Context, req *useriface.RequestPasswordResetRequest) (*useriface.RequestPasswordResetResponse, error) {
	return a.passwordPort.RequestPasswordReset(ctx, req)
}

// ConfirmPasswordReset 委托给 PasswordManagementPort
func (a *UserServiceAdapter) ConfirmPasswordReset(ctx context.Context, req *useriface.ConfirmPasswordResetRequest) (*useriface.ConfirmPasswordResetResponse, error) {
	return a.passwordPort.ConfirmPasswordReset(ctx, req)
}

// ============================================================================
// 邮箱管理方法 - 委托给 EmailManagementPort
// ============================================================================

// SendEmailVerification 委托给 EmailManagementPort
func (a *UserServiceAdapter) SendEmailVerification(ctx context.Context, req *useriface.SendEmailVerificationRequest) (*useriface.SendEmailVerificationResponse, error) {
	return a.emailPort.SendEmailVerification(ctx, req)
}

// VerifyEmail 委托给 EmailManagementPort
func (a *UserServiceAdapter) VerifyEmail(ctx context.Context, req *useriface.VerifyEmailRequest) (*useriface.VerifyEmailResponse, error) {
	return a.emailPort.VerifyEmail(ctx, req)
}

// UnbindEmail 委托给 EmailManagementPort
func (a *UserServiceAdapter) UnbindEmail(ctx context.Context, userID string) error {
	return a.emailPort.UnbindEmail(ctx, userID)
}

// EmailExists 委托给 EmailManagementPort
func (a *UserServiceAdapter) EmailExists(ctx context.Context, email string) (bool, error) {
	return a.emailPort.EmailExists(ctx, email)
}

// ============================================================================
// 权限管理方法 - 委托给 UserPermissionPort
// ============================================================================

// AssignRole 委托给 UserPermissionPort
func (a *UserServiceAdapter) AssignRole(ctx context.Context, req *useriface.AssignRoleRequest) (*useriface.AssignRoleResponse, error) {
	return a.permissionPort.AssignRole(ctx, req)
}

// RemoveRole 委托给 UserPermissionPort
func (a *UserServiceAdapter) RemoveRole(ctx context.Context, req *useriface.RemoveRoleRequest) (*useriface.RemoveRoleResponse, error) {
	return a.permissionPort.RemoveRole(ctx, req)
}

// GetUserRoles 委托给 UserPermissionPort
func (a *UserServiceAdapter) GetUserRoles(ctx context.Context, req *useriface.GetUserRolesRequest) (*useriface.GetUserRolesResponse, error) {
	return a.permissionPort.GetUserRoles(ctx, req)
}

// GetUserPermissions 委托给 UserPermissionPort
func (a *UserServiceAdapter) GetUserPermissions(ctx context.Context, req *useriface.GetUserPermissionsRequest) (*useriface.GetUserPermissionsResponse, error) {
	return a.permissionPort.GetUserPermissions(ctx, req)
}

// DowngradeRole 委托给 UserPermissionPort
func (a *UserServiceAdapter) DowngradeRole(ctx context.Context, req *useriface.DowngradeRoleRequest) (*useriface.DowngradeRoleResponse, error) {
	return a.permissionPort.DowngradeRole(ctx, req)
}

// ============================================================================
// 用户状态管理方法 - 委托给 UserStatusPort
// ============================================================================

// UpdateLastLogin 委托给 UserStatusPort
func (a *UserServiceAdapter) UpdateLastLogin(ctx context.Context, req *useriface.UpdateLastLoginRequest) (*useriface.UpdateLastLoginResponse, error) {
	return a.statusPort.UpdateLastLogin(ctx, req)
}

// DeleteDevice 委托给 UserStatusPort
func (a *UserServiceAdapter) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	return a.statusPort.DeleteDevice(ctx, userID, deviceID)
}

// UnbindPhone 委托给 UserStatusPort
func (a *UserServiceAdapter) UnbindPhone(ctx context.Context, userID string) error {
	return a.statusPort.UnbindPhone(ctx, userID)
}

// ============================================================================
// 编译时检查
// ============================================================================

// 确保 UserServiceAdapter 实现了 UserService 接口
var _ useriface.UserService = (*UserServiceAdapter)(nil)
