package interfaces

import (
	"context"

	adminInterfaces "Qingyu_backend/repository/interfaces/admin"
	AIInterfaces "Qingyu_backend/repository/interfaces/ai"
	AuditInterfaces "Qingyu_backend/repository/interfaces/audit"
	authInterface "Qingyu_backend/repository/interfaces/auth"
	BookstoreInterfaces "Qingyu_backend/repository/interfaces/bookstore"
	FinanceInterfaces "Qingyu_backend/repository/interfaces/finance"
	messagingInterfaces "Qingyu_backend/repository/interfaces/messaging"
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reader"
	RecommendationInterfaces "Qingyu_backend/repository/interfaces/recommendation"
	SharedInterfaces "Qingyu_backend/repository/interfaces/shared"
	SocialInterfaces "Qingyu_backend/repository/interfaces/social"
	StatsInterfaces "Qingyu_backend/repository/interfaces/stats"
	storageInterfaces "Qingyu_backend/repository/interfaces/storage"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/mongo"
)

// ProjectInterface 项目仓储接口
type ProjectInterface interface {
	GetProjectRepository() writer.ProjectRepository
}

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 用户相关Repository
	CreateUserRepository() UserInterface.UserRepository

	// 角色管理Repository (Auth模块)
	CreateRoleRepository() authInterface.RoleRepository

	// 写作相关Repository
	CreateProjectRepository() writer.ProjectRepository
	CreateDocumentRepository() writer.DocumentRepository
	CreateDocumentContentRepository() writer.DocumentContentRepository
	CreateTemplateRepository() writer.TemplateRepository
	// 设定百科相关Repository
	CreateCharacterRepository() writer.CharacterRepository
	CreateLocationRepository() writer.LocationRepository
	CreateTimelineRepository() writer.TimelineRepository
	CreateTimelineEventRepository() writer.TimelineEventRepository
	CreateOutlineRepository() writer.OutlineRepository

	// 阅读相关Repository
	CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository
	CreateReadingProgressRepository() ReadingInterfaces.ReadingProgressRepository
	CreateAnnotationRepository() ReadingInterfaces.AnnotationRepository
	CreateCommentRepository() ReadingInterfaces.CommentRepository
	CreateLikeRepository() ReadingInterfaces.LikeRepository
	CreateCollectionRepository() ReadingInterfaces.CollectionRepository
	CreateReadingHistoryRepository() ReadingInterfaces.ReadingHistoryRepository
	CreateReaderThemeRepository() ReadingInterfaces.ReaderThemeRepository
	CreateBookmarkRepository() ReadingInterfaces.BookmarkRepository

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

	// Auth相关Repository
	CreateOAuthRepository() authInterface.OAuthRepository
	CreateRoleAuthRepository() authInterface.RoleRepository

	// 财务相关Repository (包括Wallet)
	CreateWalletRepository() FinanceInterfaces.WalletRepository
	CreateMembershipRepository() FinanceInterfaces.MembershipRepository
	CreateAuthorRevenueRepository() FinanceInterfaces.AuthorRevenueRepository

	// Admin相关Repository
	CreateAuditRepository() adminInterfaces.AuditRepository
	CreateAdminLogRepository() adminInterfaces.AdminLogRepository

	// Messaging相关Repository
	CreateAnnouncementRepository() messagingInterfaces.AnnouncementRepository
	CreateMessageRepository() messagingInterfaces.MessageRepository

	// Storage相关Repository
	CreateStorageRepository() storageInterfaces.StorageRepository

	// ========== 向后兼容的方法 (使用 shared 接口) ==========
	// Deprecated: 这些方法为了向后兼容而保留，新代码应使用上面的新接口
	// todo: 清理掉这些方法

	// Auth相关Repository (向后兼容)
	CreateAuthRepository() SharedInterfaces.AuthRepository

	// Recommendation相关Repository (向后兼容)
	CreateRecommendationRepository() SharedInterfaces.RecommendationRepository

	// Social相关Repository
	CreateBookListRepository() SocialInterfaces.BookListRepository
	CreateReviewRepository() SocialInterfaces.ReviewRepository
	CreateFollowRepository() SocialInterfaces.FollowRepository

	// Stats相关Repository
	CreateChapterStatsRepository() StatsInterfaces.ChapterStatsRepository
	CreateReaderBehaviorRepository() StatsInterfaces.ReaderBehaviorRepository
	CreateBookStatsRepository() StatsInterfaces.BookStatsRepository

	// 审核相关Repository
	CreateSensitiveWordRepository() AuditInterfaces.SensitiveWordRepository

	// Health 基础设施方法
	Health(ctx context.Context) error
	Close() error
	GetDatabaseType() string
	GetDatabase() *mongo.Database
	GetClient() *mongo.Client
}
