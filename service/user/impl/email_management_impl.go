package impl

import (
	"context"

	portInterfaces "Qingyu_backend/service/interfaces/user"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/user"
)

// EmailManagementImpl 邮箱管理端口实现
type EmailManagementImpl struct {
	userRepo repoInterfaces.UserRepository
}

// NewEmailManagementImpl 创建实例
func NewEmailManagementImpl(
	userRepo repoInterfaces.UserRepository,
) portInterfaces.EmailManagementPort {
	return &EmailManagementImpl{
		userRepo: userRepo,
	}
}

// SendEmailVerification 发送邮箱验证
func (s *EmailManagementImpl) SendEmailVerification(ctx context.Context, req *portInterfaces.SendEmailVerificationRequest) (*portInterfaces.SendEmailVerificationResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.GetByID 接口", nil)
}

// VerifyEmail 验证邮箱
func (s *EmailManagementImpl) VerifyEmail(ctx context.Context, req *portInterfaces.VerifyEmailRequest) (*portInterfaces.VerifyEmailResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.SetEmailVerified 接口", nil)
}

// UnbindEmail 解绑邮箱
func (s *EmailManagementImpl) UnbindEmail(ctx context.Context, userID string) error {
	// TODO: 需要适配现有 Repository 接口
	return user.InternalError("功能待实现: 需要适配 userRepo.UnbindEmail 接口", nil)
}

// EmailExists 检查邮箱是否已存在
func (s *EmailManagementImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	// TODO: 需要适配现有 Repository 接口
	return false, user.InternalError("功能待实现: 需要适配 userRepo.ExistsByEmail 接口", nil)
}

// Initialize 初始化服务
func (s *EmailManagementImpl) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *EmailManagementImpl) Health(ctx context.Context) error {
	return s.userRepo.Health(ctx)
}

// Close 关闭服务
func (s *EmailManagementImpl) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *EmailManagementImpl) GetServiceName() string {
	return "EmailManagementPort"
}

// GetVersion 获取服务版本
func (s *EmailManagementImpl) GetVersion() string {
	return "1.0.0"
}
