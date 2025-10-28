package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	auditRepo "Qingyu_backend/repository/interfaces/audit"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	readingRepo "Qingyu_backend/repository/interfaces/reading"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
	userRepo "Qingyu_backend/repository/interfaces/user"
	writingRepo "Qingyu_backend/repository/interfaces/writing"

	// 导入各个子包的具体实现
	mongoAI "Qingyu_backend/repository/mongodb/ai"
	mongoBookstore "Qingyu_backend/repository/mongodb/bookstore"
	mongoReading "Qingyu_backend/repository/mongodb/reading"
	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	mongoShared "Qingyu_backend/repository/mongodb/shared"
	mongoUser "Qingyu_backend/repository/mongodb/user"
	mongoWriting "Qingyu_backend/repository/mongodb/writing"
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

// CreateRoleRepository 创建角色Repository
// 注意：目前使用 AuthRepository 作为临时实现
// TODO: 创建专门的 RoleRepository 实现或适配器
func (f *MongoRepositoryFactory) CreateRoleRepository() userRepo.RoleRepository {
	// 暂时返回 nil，需要实现正确的 RoleRepository
	// return mongoShared.NewAuthRepository(f.database)
	return nil // TODO: 实现 RoleRepository
}

// ========== Writing Module Repositories ==========

// CreateProjectRepository 创建项目Repository
func (f *MongoRepositoryFactory) CreateProjectRepository() writingRepo.ProjectRepository {
	return mongoWriting.NewMongoProjectRepository(f.database)
}

// CreateDocumentRepository 创建文档Repository
func (f *MongoRepositoryFactory) CreateDocumentRepository() writingRepo.DocumentRepository {
	return mongoWriting.NewMongoDocumentRepository(f.database)
}

// CreateDocumentContentRepository 创建文档内容Repository
func (f *MongoRepositoryFactory) CreateDocumentContentRepository() writingRepo.DocumentContentRepository {
	return mongoWriting.NewMongoDocumentContentRepository(f.database)
}

// CreateCharacterRepository 创建角色Repository
func (f *MongoRepositoryFactory) CreateCharacterRepository() writingRepo.CharacterRepository {
	return NewCharacterRepository(f.database)
}

// CreateLocationRepository 创建地点Repository
func (f *MongoRepositoryFactory) CreateLocationRepository() writingRepo.LocationRepository {
	return NewLocationRepository(f.database)
}

// CreateTimelineRepository 创建时间线Repository
func (f *MongoRepositoryFactory) CreateTimelineRepository() writingRepo.TimelineRepository {
	return NewTimelineRepository(f.database)
}

// CreateTimelineEventRepository 创建时间线事件Repository
func (f *MongoRepositoryFactory) CreateTimelineEventRepository() writingRepo.TimelineEventRepository {
	return NewTimelineEventRepository(f.database)
}

// ========== Reading Module Repositories ==========

// CreateReadingSettingsRepository 创建阅读设置Repository
func (f *MongoRepositoryFactory) CreateReadingSettingsRepository() readingRepo.ReadingSettingsRepository {
	return mongoReading.NewMongoReadingSettingsRepository(f.database)
}

// CreateChapterRepository 创建章节Repository
func (f *MongoRepositoryFactory) CreateChapterRepository() readingRepo.ChapterRepository {
	return mongoReading.NewMongoChapterRepository(f.database)
}

// CreateReadingProgressRepository 创建阅读进度Repository
func (f *MongoRepositoryFactory) CreateReadingProgressRepository() readingRepo.ReadingProgressRepository {
	return mongoReading.NewMongoReadingProgressRepository(f.database)
}

// CreateAnnotationRepository 创建注记Repository
func (f *MongoRepositoryFactory) CreateAnnotationRepository() readingRepo.AnnotationRepository {
	return mongoReading.NewMongoAnnotationRepository(f.database)
}

// CreateCommentRepository 创建评论Repository
func (f *MongoRepositoryFactory) CreateCommentRepository() readingRepo.CommentRepository {
	return mongoReading.NewMongoCommentRepository(f.database)
}

// CreateLikeRepository 创建点赞Repository
func (f *MongoRepositoryFactory) CreateLikeRepository() readingRepo.LikeRepository {
	return mongoReading.NewMongoLikeRepository(f.database)
}

// CreateCollectionRepository 创建收藏Repository
func (f *MongoRepositoryFactory) CreateCollectionRepository() readingRepo.CollectionRepository {
	return mongoReading.NewMongoCollectionRepository(f.database)
}

// CreateReadingHistoryRepository 创建阅读历史Repository
func (f *MongoRepositoryFactory) CreateReadingHistoryRepository() readingRepo.ReadingHistoryRepository {
	return mongoReading.NewMongoReadingHistoryRepository(f.database)
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

// ========== Shared Module Repositories ==========

// CreateAuthRepository 创建认证Repository
func (f *MongoRepositoryFactory) CreateAuthRepository() sharedRepo.AuthRepository {
	return mongoShared.NewAuthRepository(f.database)
}

// CreateWalletRepository 创建钱包Repository
func (f *MongoRepositoryFactory) CreateWalletRepository() sharedRepo.WalletRepository {
	return mongoShared.NewWalletRepository(f.database)
}

// CreateRecommendationRepository 创建推荐Repository
func (f *MongoRepositoryFactory) CreateRecommendationRepository() sharedRepo.RecommendationRepository {
	return mongoShared.NewRecommendationRepository(f.database)
}

// CreateStorageRepository 创建存储Repository
func (f *MongoRepositoryFactory) CreateStorageRepository() sharedRepo.StorageRepository {
	return NewMongoStorageRepository(f.database)
}

// CreateAdminRepository 创建管理后台Repository
func (f *MongoRepositoryFactory) CreateAdminRepository() sharedRepo.AdminRepository {
	return mongoShared.NewMongoAdminRepository(f.database)
}

// CreateMessageRepository 创建消息Repository
func (f *MongoRepositoryFactory) CreateMessageRepository() sharedRepo.MessageRepository {
	return mongoShared.NewMongoMessageRepository(f.database)
}

// ========== AI Module Repositories ==========

// CreateQuotaRepository 创建配额Repository
func (f *MongoRepositoryFactory) CreateQuotaRepository() aiRepo.QuotaRepository {
	return mongoAI.NewMongoQuotaRepository(f.database)
}

// ========== Audit Module Repositories ==========

// CreateSensitiveWordRepository 创建敏感词Repository
// 注意：当前返回nil，表示敏感词检测未启用
// TODO: 实现敏感词Repository
func (f *MongoRepositoryFactory) CreateSensitiveWordRepository() auditRepo.SensitiveWordRepository {
	// 暂时返回nil，CommentService会处理nil的情况
	return nil
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
