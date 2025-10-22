package interfaces

import (
	"context"

	ReadingInterfaces "Qingyu_backend/repository/interfaces/reading"
	RecommendationInterfaces "Qingyu_backend/repository/interfaces/recommendation"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/repository/interfaces/writing"
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

	// 文档相关Repository
	CreateDocumentRepository() writing.DocumentRepository
	CreateDocumentContentRepository() writing.DocumentContentRepository

	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository

	// 推荐系统相关Repository
	CreateBehaviorRepository() RecommendationInterfaces.BehaviorRepository
	CreateProfileRepository() RecommendationInterfaces.ProfileRepository
	CreateItemFeatureRepository() RecommendationInterfaces.ItemFeatureRepository
	CreateHotRecommendationRepository() RecommendationInterfaces.HotRecommendationRepository

	// 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}
