package user

import (
	"Qingyu_backend/service/interfaces/user"
	usermigration "Qingyu_backend/service/user/_migration"
)

// UserServiceFactory 用户服务工厂
// 提供创建和组装用户服务的方法
type UserServiceFactory struct{}

// NewUserServiceFactory 创建工厂实例
func NewUserServiceFactory() *UserServiceFactory {
	return &UserServiceFactory{}
}

// CreateWithPorts 使用 Port 接口创建服务（推荐方式）
//
// 新架构推荐的使用方式：
// 1. 实现 6 个 Port 接口的具体实现
// 2. 使用 UserServiceAdapter 组装它们
// 3. 返回 UserService 接口供 API 层使用
func (f *UserServiceFactory) CreateWithPorts(
	managementPort user.UserManagementPort,
	authPort user.UserAuthPort,
	passwordPort user.PasswordManagementPort,
	emailPort user.EmailManagementPort,
	permissionPort user.UserPermissionPort,
	statusPort user.UserStatusPort,
) user.UserService {
	return usermigration.NewUserServiceAdapter(
		managementPort,
		authPort,
		passwordPort,
		emailPort,
		permissionPort,
		statusPort,
	)
}

// PortImplementations Port 接口实现集合
type PortImplementations struct {
	ManagementPort user.UserManagementPort
	AuthPort       user.UserAuthPort
	PasswordPort   user.PasswordManagementPort
	EmailPort      user.EmailManagementPort
	PermissionPort user.UserPermissionPort
	StatusPort     user.UserStatusPort
}

// CreateFromImplementations 从结构体创建服务
func (f *UserServiceFactory) CreateFromImplementations(ports PortImplementations) user.UserService {
	return usermigration.NewUserServiceAdapter(
		ports.ManagementPort,
		ports.AuthPort,
		ports.PasswordPort,
		ports.EmailPort,
		ports.PermissionPort,
		ports.StatusPort,
	)
}
