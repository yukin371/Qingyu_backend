package interfaces

import (
	"context"

	AIInterfaces "Qingyu_backend/repository/interfaces/ai"
	BookstoreInterfaces "Qingyu_backend/repository/interfaces/bookstore"
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
	CreateRoleRepository() UserInterface.RoleRepository

	// 写作相关Repository
	CreateProjectRepository() writing.ProjectRepository
	CreateDocumentRepository() writing.DocumentRepository
	CreateDocumentContentRepository() writing.DocumentContentRepository

	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository
	CreateChapterRepository() ReadingInterfaces.ChapterRepository
	CreateReadingProgressRepository() ReadingInterfaces.ReadingProgressRepository
	CreateAnnotationRepository() ReadingInterfaces.AnnotationRepository

	// 书城相关Repository
	CreateBookRepository() BookstoreInterfaces.BookRepository
	CreateBookDetailRepository() BookstoreInterfaces.BookDetailRepository
	CreateCategoryRepository() BookstoreInterfaces.CategoryRepository
	CreateBookStatisticsRepository() BookstoreInterfaces.BookStatisticsRepository
	CreateBookRatingRepository() BookstoreInterfaces.BookRatingRepository
	CreateBookstoreChapterRepository() BookstoreInterfaces.ChapterRepository
	CreateBannerRepository() BookstoreInterfaces.BannerRepository
	CreateRankingRepository() BookstoreInterfaces.RankingRepository

	// AI相关Repository
	CreateQuotaRepository() AIInterfaces.QuotaRepository

	// 推荐系统相关Repository
	CreateBehaviorRepository() RecommendationInterfaces.BehaviorRepository
	CreateProfileRepository() RecommendationInterfaces.ProfileRepository
	CreateItemFeatureRepository() RecommendationInterfaces.ItemFeatureRepository
	CreateHotRecommendationRepository() RecommendationInterfaces.HotRecommendationRepository

	// 共享服务相关Repository
	// 注意：这些方法暂时不在标准接口中，因为SharedService还未完全集成
	// CreateAuthRepository() interface{}
	// CreateWalletRepository() interface{}
	// CreateRecommendationRepository() interface{}

	// Health 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}
