package interfaces

import (
	"context"

	AIInterfaces "Qingyu_backend/repository/interfaces/ai"
	authInterface "Qingyu_backend/repository/interfaces/auth"
	AuditInterfaces "Qingyu_backend/repository/interfaces/audit"
	BookstoreInterfaces "Qingyu_backend/repository/interfaces/bookstore"
	FinanceInterfaces "Qingyu_backend/repository/interfaces/finance"
	messagingInterfaces "Qingyu_backend/repository/interfaces/messaging"
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reader"
	RecommendationInterfaces "Qingyu_backend/repository/interfaces/recommendation"
	SharedInterfaces "Qingyu_backend/repository/interfaces/shared"
	SocialInterfaces "Qingyu_backend/repository/interfaces/social"
	StatsInterfaces "Qingyu_backend/repository/interfaces/stats"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/repository/interfaces/writer"
)

// ProjectInterface 项目仓储接口
type ProjectInterface interface {
	GetProjectRepository() writer.ProjectRepository
}

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 用户相关Repository
	CreateUserRepository() UserInterface.UserRepository
	CreateRoleRepository() UserInterface.RoleRepository

	// 写作相关Repository
	CreateProjectRepository() writer.ProjectRepository
	CreateDocumentRepository() writer.DocumentRepository
	CreateDocumentContentRepository() writer.DocumentContentRepository

	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository
	CreateReadingProgressRepository() ReadingInterfaces.ReadingProgressRepository
	CreateAnnotationRepository() ReadingInterfaces.AnnotationRepository
	CreateCommentRepository() ReadingInterfaces.CommentRepository
	CreateLikeRepository() ReadingInterfaces.LikeRepository
	CreateCollectionRepository() ReadingInterfaces.CollectionRepository
	CreateReadingHistoryRepository() ReadingInterfaces.ReadingHistoryRepository
	CreateReaderThemeRepository() ReadingInterfaces.ReaderThemeRepository

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
	CreateOAuthRepository() authInterface.OAuthRepository
	CreateWalletRepository() SharedInterfaces.WalletRepository
	CreateRecommendationRepository() SharedInterfaces.RecommendationRepository
	CreateStorageRepository() SharedInterfaces.StorageRepository

	// Messaging相关Repository
	CreateAnnouncementRepository() messagingInterfaces.AnnouncementRepository

	// Social相关Repository
	CreateBookListRepository() SocialInterfaces.BookListRepository

	// Stats相关Repository
	CreateChapterStatsRepository() StatsInterfaces.ChapterStatsRepository
	CreateReaderBehaviorRepository() StatsInterfaces.ReaderBehaviorRepository
	CreateBookStatsRepository() StatsInterfaces.BookStatsRepository

	// 审核相关Repository
	CreateSensitiveWordRepository() AuditInterfaces.SensitiveWordRepository

	// 财务相关Repository
	CreateMembershipRepository() FinanceInterfaces.MembershipRepository
	CreateAuthorRevenueRepository() FinanceInterfaces.AuthorRevenueRepository

	// Health 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
}
