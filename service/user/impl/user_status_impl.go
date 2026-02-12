package impl

import (
	"context"

	portInterfaces "Qingyu_backend/service/interfaces/user"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/user"
)

// UserStatusImpl 用户状态端口实现
type UserStatusImpl struct {
	userRepo repoInterfaces.UserRepository
}

// NewUserStatusImpl 创建实例
func NewUserStatusImpl(
	userRepo repoInterfaces.UserRepository,
) portInterfaces.UserStatusPort {
	return &UserStatusImpl{
		userRepo: userRepo,
	}
}

// UpdateLastLogin 更新最后登录时间
func (s *UserStatusImpl) UpdateLastLogin(ctx context.Context, req *portInterfaces.UpdateLastLoginRequest) (*portInterfaces.UpdateLastLoginResponse, error) {
	// TODO: 需要适配现有 Repository 接口
	return nil, user.InternalError("功能待实现: 需要适配 userRepo.UpdateLastLogin 接口", nil)
}

// DeleteDevice 删除设备
func (s *UserStatusImpl) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	// TODO: 需要适配现有 Repository 接口
	return user.InternalError("功能待实现: 需要适配 userRepo.DeleteDevice 接口", nil)
}

// UnbindPhone 解绑手机
func (s *UserStatusImpl) UnbindPhone(ctx context.Context, userID string) error {
	// TODO: 需要适配现有 Repository 接口
	return user.InternalError("功能待实现: 需要适配 userRepo.UnbindPhone 接口", nil)
}

// Initialize 初始化服务
func (s *UserStatusImpl) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *UserStatusImpl) Health(ctx context.Context) error {
	return s.userRepo.Health(ctx)
}

// Close 关闭服务
func (s *UserStatusImpl) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *UserStatusImpl) GetServiceName() string {
	return "UserStatusPort"
}

// GetVersion 获取服务版本
func (s *UserStatusImpl) GetVersion() string {
	return "1.0.0"
}
