package interfaces

import (
	"context"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reading"
)

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 用户相关Repository
	CreateUserRepository() UserInterface.UserRepository
	CreateProjectRepository() UserInterface.ProjectRepository
	CreateRoleRepository() UserInterface.RoleRepository
	
	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository
	
	// 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}