package interfaces

import (
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reading"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/repository/interfaces/writing"
	"context"
)

// ProjectInterface 项目仓储接口
type ProjectInterface interface {
	GetProjectRepository() writing.ProjectRepository
}

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 用户相关Repository
	CreateUserRepository() UserInterface.UserRepository
	CreateProjectRepository() writing.ProjectRepository
	CreateRoleRepository() UserInterface.RoleRepository

	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository

	// 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}
