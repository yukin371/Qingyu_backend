package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
	adminRepo "Qingyu_backend/repository/interfaces/admin"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	auditRepo "Qingyu_backend/repository/interfaces/audit"
	authRepo "Qingyu_backend/repository/interfaces/auth"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	financeRepo "Qingyu_backend/repository/interfaces/finance"
	messagingRepo "Qingyu_backend/repository/interfaces/messaging"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	statsRepo "Qingyu_backend/repository/interfaces/stats"
	storageRepo "Qingyu_backend/repository/interfaces/storage"
	userRepo "Qingyu_backend/repository/interfaces/user"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	// 导入各个子包的具体实现
	mongoAdmin "Qingyu_backend/repository/mongodb/admin"
	mongoAI "Qingyu_backend/repository/mongodb/ai"
	mongoAudit "Qingyu_backend/repository/mongodb/audit"
	mongoAuth "Qingyu_backend/repository/mongodb/auth"
	mongoBookstore "Qingyu_backend/repository/mongodb/bookstore"
	mongoFinance "Qingyu_backend/repository/mongodb/finance"
	mongoMessaging "Qingyu_backend/repository/mongodb/messaging"
	mongoReading "Qingyu_backend/repository/mongodb/reader"
	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	mongoSocial "Qingyu_backend/repository/mongodb/social"
	mongoStats "Qingyu_backend/repository/mongodb/stats"
	mongoStorage "Qingyu_backend/repository/mongodb/storage"
	mongoUser "Qingyu_backend/repository/mongodb/user"
	mongoWriter "Qingyu_backend/repository/mongodb/writer"
)

// MongoRepositoryFactory MongoDB仓储工厂实现
type MongoRepositoryFactory struct {
	client   *mongo.Client
	db       *mongo.Database // 也保留别名以兼容旧代码
	database *mongo.Database
	config   *config.MongoDBConfig
}

// NewMongoRepositoryFactory 创建MongoDB仓储工厂
func NewMongoRepositoryFactory(config *config.MongoDBConfig) (*MongoRepositoryFactory, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("MongoDB配置验证失败: %w", err)
	}

	// 创建客户端选项
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetConnectTimeout(config.ConnectTimeout).
		SetServerSelectionTimeout(config.ServerTimeout)

	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接MongoDB失败: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("MongoDB连接测试失败: %w", err)
	}

	database := client.Database(config.Database)

	return &MongoRepositoryFactory{
		client:   client,
		db:       database,
		database: database,
		config:   config,
	}, nil
}

// NewMongoRepositoryFactoryWithClient 使用已有MongoDB连接创建工厂
// 推荐使用此方法，避免重复创建连接
// 此方法从ServiceContainer获取已初始化的MongoDB连接
func NewMongoRepositoryFactoryWithClient(
	client *mongo.Client,
	db *mongo.Database,
) *MongoRepositoryFactory {
	return &MongoRepositoryFactory{
		client:   client,
		db:       db,
		database: db,
		config:   nil, // 不需要配置，因为连接已创建
	}
}

// ========== User Module Repositories ==========

// CreateUserRepository 创建用户Repository
func (f *MongoRepositoryFactory) CreateUserRepository() userRepo.UserRepository {
	return mongoUser.NewMongoUserRepository(f.database)
}

// CreateRoleRepository 创建角色Repository (使用新的 auth 模块)
func (f *MongoRepositoryFactory) CreateRoleRepository() authRepo.RoleRepository {
	return mongoAuth.NewRoleRepository(f.database)
}

// ========== Writing Module Repositories ==========

// CreateProjectRepository 创建项目Repository
func (f *MongoRepositoryFactory) CreateProjectRepository() writerRepo.ProjectRepository {
	return mongoWriter.NewMongoProjectRepository(f.database)
}

// CreateDocumentRepository 创建文档Repository
func (f *MongoRepositoryFactory) CreateDocumentRepository() writerRepo.DocumentRepository {
	return mongoWriter.NewMongoDocumentRepository(f.database)
}

// CreateDocumentContentRepository 创建文档内容Repository
func (f *MongoRepositoryFactory) CreateDocumentContentRepository() writerRepo.DocumentContentRepository {
	return mongoWriter.NewMongoDocumentContentRepository(f.database)
}

// CreateCharacterRepository 创建角色Repository
func (f *MongoRepositoryFactory) CreateCharacterRepository() writerRepo.CharacterRepository {
	return mongoWriter.NewCharacterRepository(f.database)
}

// CreateLocationRepository 创建地点Repository
func (f *MongoRepositoryFactory) CreateLocationRepository() writerRepo.LocationRepository {
	return mongoWriter.NewLocationRepository(f.database)
}

// CreateTimelineRepository 创建时间线Repository
func (f *MongoRepositoryFactory) CreateTimelineRepository() writerRepo.TimelineRepository {
	return mongoWriter.NewTimelineRepository(f.database)
}

// CreateTimelineEventRepository 创建时间线事件Repository
func (f *MongoRepositoryFactory) CreateTimelineEventRepository() writerRepo.TimelineEventRepository {
	return mongoWriter.NewTimelineEventRepository(f.database)
}

// ========== Reading Module Repositories ==========

// CreateReadingSettingsRepository 创建阅读设置Repository
func (f *MongoRepositoryFactory) CreateReadingSettingsRepository() readerRepo.ReadingSettingsRepository {
	return mongoReading.NewMongoReadingSettingsRepository(f.database)
}

// CreateReadingProgressRepository 创建阅读进度Repository
func (f *MongoRepositoryFactory) CreateReadingProgressRepository() readerRepo.ReadingProgressRepository {
	return mongoReading.NewMongoReadingProgressRepository(f.database)
}

// CreateAnnotationRepository 创建注记Repository
func (f *MongoRepositoryFactory) CreateAnnotationRepository() readerRepo.AnnotationRepository {
	return mongoReading.NewMongoAnnotationRepository(f.database)
}

// CreateCommentRepository 创建评论Repository
func (f *MongoRepositoryFactory) CreateCommentRepository() readerRepo.CommentRepository {
	return mongoReading.NewMongoCommentRepository(f.database)
}

// CreateLikeRepository 创建点赞Repository
func (f *MongoRepositoryFactory) CreateLikeRepository() readerRepo.LikeRepository {
	return mongoReading.NewMongoLikeRepository(f.database)
}

// CreateCollectionRepository 创建收藏Repository
func (f *MongoRepositoryFactory) CreateCollectionRepository() readerRepo.CollectionRepository {
	return mongoReading.NewMongoCollectionRepository(f.database)
}

// CreateBookListRepository 创建书单Repository
func (f *MongoRepositoryFactory) CreateBookListRepository() socialRepo.BookListRepository {
	return mongoSocial.NewMongoBookListRepository(f.database)
}

// CreateFollowRepository 创建关注Repository
func (f *MongoRepositoryFactory) CreateFollowRepository() socialRepo.FollowRepository {
	return mongoSocial.NewMongoFollowRepository(f.database)
}

// CreateReadingHistoryRepository 创建阅读历史Repository
func (f *MongoRepositoryFactory) CreateReadingHistoryRepository() readerRepo.ReadingHistoryRepository {
	return mongoReading.NewMongoReadingHistoryRepository(f.database)
}

// CreateReaderThemeRepository 创建阅读器主题Repository
func (f *MongoRepositoryFactory) CreateReaderThemeRepository() readerRepo.ReaderThemeRepository {
	return mongoReading.NewReaderThemeRepositoryMongo(f.database)
}

// ========== Bookstore Module Repositories ==========

// CreateBookRepository 创建书籍Repository
func (f *MongoRepositoryFactory) CreateBookRepository() bookstoreRepo.BookRepository {
	return mongoBookstore.NewMongoBookRepository(f.client, f.database.Name())
}

// CreateBookDetailRepository 创建书籍详情Repository
func (f *MongoRepositoryFactory) CreateBookDetailRepository() bookstoreRepo.BookDetailRepository {
	return mongoBookstore.NewMongoBookDetailRepository(f.client, f.database.Name())
}

// CreateCategoryRepository 创建分类Repository
func (f *MongoRepositoryFactory) CreateCategoryRepository() bookstoreRepo.CategoryRepository {
	return mongoBookstore.NewMongoCategoryRepository(f.client, f.database.Name())
}

// CreateBookStatisticsRepository 创建书籍统计Repository
func (f *MongoRepositoryFactory) CreateBookStatisticsRepository() bookstoreRepo.BookStatisticsRepository {
	return mongoBookstore.NewMongoBookStatisticsRepository(f.client, f.database.Name())
}

// CreateBookRatingRepository 创建书籍评分Repository
func (f *MongoRepositoryFactory) CreateBookRatingRepository() bookstoreRepo.BookRatingRepository {
	return mongoBookstore.NewMongoBookRatingRepository(f.client, f.database.Name())
}

// CreateBookstoreChapterRepository 创建书城章节Repository
func (f *MongoRepositoryFactory) CreateBookstoreChapterRepository() bookstoreRepo.ChapterRepository {
	return mongoBookstore.NewMongoChapterRepository(f.client, f.database.Name())
}

// CreateChapterContentRepository 创建章节内容Repository
func (f *MongoRepositoryFactory) CreateChapterContentRepository() bookstoreRepo.ChapterContentRepository {
	return mongoBookstore.NewMongoChapterContentRepository(f.database)
}

// CreateBannerRepository 创建横幅Repository
func (f *MongoRepositoryFactory) CreateBannerRepository() bookstoreRepo.BannerRepository {
	return mongoBookstore.NewMongoBannerRepository(f.client, f.database.Name())
}

// CreateRankingRepository 创建榜单Repository
func (f *MongoRepositoryFactory) CreateRankingRepository() bookstoreRepo.RankingRepository {
	return mongoBookstore.NewMongoRankingRepository(f.client, f.database.Name())
}

// ========== Recommendation Module Repositories ==========

// CreateBehaviorRepository 创建行为Repository
func (f *MongoRepositoryFactory) CreateBehaviorRepository() recoRepo.BehaviorRepository {
	return mongoReco.NewMongoBehaviorRepository(f.database)
}

// CreateProfileRepository 创建用户画像Repository
func (f *MongoRepositoryFactory) CreateProfileRepository() recoRepo.ProfileRepository {
	return mongoReco.NewMongoProfileRepository(f.database)
}

// CreateItemFeatureRepository 创建物品特征Repository
func (f *MongoRepositoryFactory) CreateItemFeatureRepository() recoRepo.ItemFeatureRepository {
	return mongoReco.NewMongoItemFeatureRepository(f.database)
}

// CreateHotRecommendationRepository 创建热门推荐Repository
func (f *MongoRepositoryFactory) CreateHotRecommendationRepository() recoRepo.HotRecommendationRepository {
	return mongoReco.NewMongoHotRecommendationRepository(f.database)
}

// ========== Auth Module Repositories ==========

// CreateOAuthRepository 创建OAuth Repository
func (f *MongoRepositoryFactory) CreateOAuthRepository() authRepo.OAuthRepository {
	return mongoAuth.NewMongoOAuthRepository(f.database)
}

// CreateRoleAuthRepository 创建角色认证Repository (使用新的 auth 模块)
func (f *MongoRepositoryFactory) CreateRoleAuthRepository() authRepo.RoleRepository {
	return mongoAuth.NewRoleRepository(f.database)
}

// ========== Finance Module Repositories (包括Wallet) ==========

// CreateWalletRepository 创建钱包Repository (使用新的 finance 模块)
func (f *MongoRepositoryFactory) CreateWalletRepository() financeRepo.WalletRepository {
	return mongoFinance.NewWalletRepository(f.database)
}

// CreateMembershipRepository 创建会员Repository
func (f *MongoRepositoryFactory) CreateMembershipRepository() financeRepo.MembershipRepository {
	return mongoFinance.NewMembershipRepository(f.database)
}

// CreateAuthorRevenueRepository 创建作者收入Repository
func (f *MongoRepositoryFactory) CreateAuthorRevenueRepository() financeRepo.AuthorRevenueRepository {
	return mongoFinance.NewAuthorRevenueRepository(f.database)
}

// ========== Admin Module Repositories ==========

// CreateAuditRepository 创建审核记录Repository (使用新的 admin 模块)
func (f *MongoRepositoryFactory) CreateAuditRepository() adminRepo.AuditRepository {
	return mongoAdmin.NewAuditRepository(f.database)
}

// CreateAdminLogRepository 创建管理员日志Repository (使用新的 admin 模块)
func (f *MongoRepositoryFactory) CreateAdminLogRepository() adminRepo.AdminLogRepository {
	return mongoAdmin.NewAdminLogRepository(f.database)
}

// ========== Messaging Module Repositories ==========

// CreateAnnouncementRepository 创建公告Repository
func (f *MongoRepositoryFactory) CreateAnnouncementRepository() messagingRepo.AnnouncementRepository {
	return mongoMessaging.NewMongoAnnouncementRepository(f.client, f.database.Name())
}

// CreateMessageRepository 创建消息Repository (使用新的 messaging 模块)
func (f *MongoRepositoryFactory) CreateMessageRepository() messagingRepo.MessageRepository {
	return mongoMessaging.NewMessageRepository(f.database)
}

// ========== Storage Module Repositories ==========

// CreateStorageRepository 创建存储Repository (使用新的 storage 模块)
func (f *MongoRepositoryFactory) CreateStorageRepository() storageRepo.StorageRepository {
	return mongoStorage.NewMongoStorageRepository(f.database)
}

// ========== 向后兼容的方法 ==========
// Deprecated: 这些方法为了向后兼容而保留，新代码应使用上面的新接口

// CreateAuthRepository 创建认证Repository (向后兼容，实际返回 RoleRepository)
func (f *MongoRepositoryFactory) CreateAuthRepository() sharedRepo.AuthRepository {
	// shared.AuthRepository 是 auth.RoleRepository 的别名
	return mongoAuth.NewRoleRepository(f.database)
}

// CreateRecommendationRepository 创建推荐Repository (向后兼容)
func (f *MongoRepositoryFactory) CreateRecommendationRepository() sharedRepo.RecommendationRepository {
	// shared.RecommendationRepository 是 recommendation.RecommendationRepository 的别名
	return mongoReco.NewRecommendationRepository(f.database)
}

// ========== Stats Module Repositories ==========

// CreateChapterStatsRepository 创建章节统计Repository
func (f *MongoRepositoryFactory) CreateChapterStatsRepository() statsRepo.ChapterStatsRepository {
	return mongoStats.NewMongoChapterStatsRepository(f.database)
}

// CreateReaderBehaviorRepository 创建读者行为Repository
func (f *MongoRepositoryFactory) CreateReaderBehaviorRepository() statsRepo.ReaderBehaviorRepository {
	return mongoStats.NewMongoReaderBehaviorRepository(f.database)
}

// CreateBookStatsRepository 创建作品统计Repository
func (f *MongoRepositoryFactory) CreateBookStatsRepository() statsRepo.BookStatsRepository {
	return mongoStats.NewMongoBookStatsRepository(f.database)
}

// ========== AI Module Repositories ==========

// CreateQuotaRepository 创建配额Repository
func (f *MongoRepositoryFactory) CreateQuotaRepository() aiRepo.QuotaRepository {
	return mongoAI.NewMongoQuotaRepository(f.database)
}

// ========== Audit Module Repositories ==========

// CreateSensitiveWordRepository 创建敏感词Repository
func (f *MongoRepositoryFactory) CreateSensitiveWordRepository() auditRepo.SensitiveWordRepository {
	return mongoAudit.NewSensitiveWordRepository(f.database)
}

// ========== Factory Management Methods ==========

// Close 关闭数据库连接
func (f *MongoRepositoryFactory) Close() error {
	if f.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return f.client.Disconnect(ctx)
	}
	return nil
}

// Health 健康检查
func (f *MongoRepositoryFactory) Health(ctx context.Context) error {
	if f.client == nil {
		return fmt.Errorf("MongoDB客户端未初始化")
	}

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return f.client.Ping(pingCtx, nil)
}

// GetDatabaseType 获取数据库类型
func (f *MongoRepositoryFactory) GetDatabaseType() string {
	return "mongodb"
}

// GetDatabase 获取数据库实例（用于事务等高级操作）
func (f *MongoRepositoryFactory) GetDatabase() *mongo.Database {
	return f.database
}

// GetClient 获取客户端实例（用于事务等高级操作）
func (f *MongoRepositoryFactory) GetClient() *mongo.Client {
	return f.client
}

// GetDatabaseName 获取数据库名称
func (f *MongoRepositoryFactory) GetDatabaseName() string {
	if f.database != nil {
		return f.database.Name()
	}
	if f.config != nil {
		return f.config.Database
	}
	return ""
}
