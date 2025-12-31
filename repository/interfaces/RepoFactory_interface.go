package interfaces

import (
	"context"

	AIInterfaces "Qingyu_backend/repository/interfaces/ai"
	AuditInterfaces "Qingyu_backend/repository/interfaces/audit"
	BookstoreInterfaces "Qingyu_backend/repository/interfaces/bookstore"
	messagingInterfaces "Qingyu_backend/repository/interfaces/messaging"
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reading"
	RecommendationInterfaces "Qingyu_backend/repository/interfaces/recommendation"
	SharedInterfaces "Qingyu_backend/repository/interfaces/shared"
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
	CreateReadingProgressRepository() ReadingInterfaces.ReadingProgressRepository
	CreateAnnotationRepository() ReadingInterfaces.AnnotationRepository
	CreateCommentRepository() ReadingInterfaces.CommentRepository
	CreateLikeRepository() ReadingInterfaces.LikeRepository
	CreateCollectionRepository() ReadingInterfaces.CollectionRepository
	CreateReadingHistoryRepository() ReadingInterfaces.ReadingHistoryRepository

	// 书城相关Repository
	CreateBookRepository() BookstoreInterfaces.BookRepository
	CreateBookDetailRepository() BookstoreInterfaces.BookDetailRepository
	CreateCategoryRepository() BookstoreInterfaces.CategoryRepository
	CreateBookStatisticsRepository() BookstoreInterfaces.BookStatisticsRepository
	CreateBookRatingRepository() BookstoreInterfaces.BookRatingRepository
	CreateBookstoreChapterRepository() BookstoreInterfaces.ChapterRepository
	CreateChapterContentRepository() BookstoreInterfaces.ChapterContentRepository
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
	CreateAuthRepository() SharedInterfaces.AuthRepository
	CreateWalletRepository() SharedInterfaces.WalletRepository
	CreateRecommendationRepository() SharedInterfaces.RecommendationRepository
	CreateStorageRepository() SharedInterfaces.StorageRepository

	// Messaging相关Repository
	CreateAnnouncementRepository() messagingInterfaces.AnnouncementRepository

	// 审核相关Repository
	CreateSensitiveWordRepository() AuditInterfaces.SensitiveWordRepository

	// Health 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}
