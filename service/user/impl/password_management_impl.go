package impl

import (
	"context"

	portInterfaces "Qingyu_backend/service/interfaces/user"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/user"
)

// PasswordManagementImpl 密码管理端口实现
type PasswordManagementImpl struct {
	userRepo repoInterfaces.UserRepository
}

// NewPasswordManagementImpl 创建实例
func NewPasswordManagementImpl(
	userRepo repoInterfaces.UserRepository,
) portInterfaces.PasswordManagementPort {
	return &PasswordManagementImpl{
		userRepo: userRepo,
	}
}

// UpdatePassword 更新密码
func (s *PasswordManagementImpl) UpdatePassword(ctx context.Context, req *portInterfaces.UpdatePasswordRequest) (*portInterfaces.UpdatePasswordResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.GetByID/UpdatePassword 接口", nil)
}

// ResetPassword 重置密码（保持兼容性）
func (s *PasswordManagementImpl) ResetPassword(ctx context.Context, req *portInterfaces.ResetPasswordRequest) (*portInterfaces.ResetPasswordResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.GetByEmail/UpdatePasswordByEmail 接口", nil)
}

// RequestPasswordReset 请求密码重置
func (s *PasswordManagementImpl) RequestPasswordReset(ctx context.Context, req *portInterfaces.RequestPasswordResetRequest) (*portInterfaces.RequestPasswordResetResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.GetByEmail 接口", nil)
}

// ConfirmPasswordReset 确认密码重置
func (s *PasswordManagementImpl) ConfirmPasswordReset(ctx context.Context, req *portInterfaces.ConfirmPasswordResetRequest) (*portInterfaces.ConfirmPasswordResetResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.UpdatePasswordByEmail 接口", nil)
}

// VerifyPassword 验证密码
func (s *PasswordManagementImpl) VerifyPassword(ctx context.Context, userID string, password string) error {
	// TODO: 需要适配现有 Repository 接口
	return user.InternalError("功能待实现: 需要适配 userRepo.GetByID 接口", nil)
}

// Initialize 初始化服务
func (s *PasswordManagementImpl) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *PasswordManagementImpl) Health(ctx context.Context) error {
	return s.userRepo.Health(ctx)
}

// Close 关闭服务
func (s *PasswordManagementImpl) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *PasswordManagementImpl) GetServiceName() string {
	return "PasswordManagementPort"
}

// GetVersion 获取服务版本
func (s *PasswordManagementImpl) GetVersion() string {
	return "1.0.0"
}
