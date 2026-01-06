package stats

import (
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	userRepo "Qingyu_backend/repository/interfaces/user"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// PlatformStatsService 平台统计服务接口
// 职责：跨业务域的聚合统计、平台级数据分析
// 领域：Platform/Shared
type PlatformStatsService interface {
	// 用户统计
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	GetPlatformUserStats(ctx context.Context, startDate, endDate time.Time) (*PlatformUserStats, error)

	// 内容统计
	GetContentStats(ctx context.Context, userID string) (*ContentStats, error)
	GetPlatformContentStats(ctx context.Context, startDate, endDate time.Time) (*PlatformContentStats, error)

	// 活跃度统计
	GetUserActivityStats(ctx context.Context, userID string, days int) (*ActivityStats, error)

	// 收益统计
	GetRevenueStats(ctx context.Context, userID string, startDate, endDate time.Time) (*RevenueStats, error)

	// 健康检查
	Health(ctx context.Context) error
}

// PlatformStatsServiceImpl 平台统计服务实现
type PlatformStatsServiceImpl struct {
	userRepo    userRepo.UserRepository
	bookRepo    bookstoreRepo.BookRepository
	projectRepo writerRepo.ProjectRepository
	chapterRepo bookstoreRepo.ChapterRepository

	initialized bool
}

// NewPlatformStatsService 创建平台统计服务
func NewPlatformStatsService(
	userRepository userRepo.UserRepository,
	bookRepository bookstoreRepo.BookRepository,
	projectRepository writerRepo.ProjectRepository,
	chapterRepository bookstoreRepo.ChapterRepository,
) PlatformStatsService {
	return &PlatformStatsServiceImpl{
		userRepo:    userRepository,
		bookRepo:    bookRepository,
		projectRepo: projectRepository,
		chapterRepo: chapterRepository,
		initialized: false,
	}
}

// ============ 用户统计 ============

// UserStats 用户统计数据
type UserStats struct {
	UserID        string    `json:"user_id"`
	TotalProjects int64     `json:"total_projects"` // 总项目数
	TotalBooks    int64     `json:"total_books"`    // 总书籍数
	TotalWords    int64     `json:"total_words"`    // 总字数
	TotalReading  int64     `json:"total_reading"`  // 总阅读数
	TotalLikes    int64     `json:"total_likes"`    // 总点赞数
	TotalComments int64     `json:"total_comments"` // 总评论数
	TotalRevenue  float64   `json:"total_revenue"`  // 总收益
	MemberLevel   string    `json:"member_level"`   // 会员等级
	RegisteredAt  time.Time `json:"registered_at"`  // 注册时间
	LastActiveAt  time.Time `json:"last_active_at"` // 最后活跃时间
	ActiveDays    int       `json:"active_days"`    // 活跃天数
}

// GetUserStats 获取用户统计
func (s *PlatformStatsServiceImpl) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 1. 获取用户基本信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 2. 统计项目数
	projectCount, err := s.projectRepo.CountByOwner(ctx, userID)
	if err != nil {
		zap.L().Warn("统计项目数失败，使用默认值0",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		projectCount = 0
	}

	// 3. 统计书籍数和总字数
	// 注意：当前GetByAuthorID需要ObjectID，暂时跳过书籍统计
	// TODO(优化): 扩展BookRepository支持string author_id查询，或建立author_id索引
	bookCount := int64(0)
	totalWords := int64(0)

	zap.L().Debug("书籍统计暂时跳过",
		zap.String("user_id", userID),
		zap.String("reason", "需要ObjectID类型author_id，当前用户ID为string类型"),
	)

	// 4. 构建统计结果
	// VIPLevel字段在User模型中不存在，暂时使用Role
	memberLevel := "普通用户"
	if user.Role == "author" {
		memberLevel = "作者"
	} else if user.Role == "admin" {
		memberLevel = "管理员"
	}

	stats := &UserStats{
		UserID:        userID,
		TotalProjects: projectCount,
		TotalBooks:    bookCount,   // TODO(Task3): 需要支持string author_id查询
		TotalWords:    totalWords,  // TODO(Task3): 需要支持string author_id查询
		TotalReading:  0,           // TODO(Task3): 需要阅读行为统计
		TotalLikes:    0,           // TODO(Task3): 需要点赞统计
		TotalComments: 0,           // TODO(Task3): 需要评论统计
		TotalRevenue:  0,           // TODO(Task3): 需要钱包统计
		MemberLevel:   memberLevel, // 暂时使用Role代替VIPLevel
		RegisteredAt:  user.CreatedAt,
		LastActiveAt:  user.UpdatedAt,
		ActiveDays:    0, // TODO(Task3): 需要活跃度统计
	}

	return stats, nil
}

// PlatformUserStats 平台用户统计
type PlatformUserStats struct {
	TotalUsers       int64   `json:"total_users"`        // 总用户数
	NewUsers         int64   `json:"new_users"`          // 新增用户
	ActiveUsers      int64   `json:"active_users"`       // 活跃用户
	VIPUsers         int64   `json:"vip_users"`          // VIP用户
	RetentionRate    float64 `json:"retention_rate"`     // 留存率
	AverageActiveDay float64 `json:"average_active_day"` // 平均活跃天数
}

// GetPlatformUserStats 获取平台用户统计
func (s *PlatformStatsServiceImpl) GetPlatformUserStats(ctx context.Context, startDate, endDate time.Time) (*PlatformUserStats, error) {
	// TODO(Task3-聚合查询): 实现MongoDB聚合管道统计
	// 需要实现：
	// 1. db.users.countDocuments({}) - 总用户数
	// 2. db.users.countDocuments({created_at: {$gte: startDate, $lte: endDate}}) - 新增用户
	// 3. db.users.countDocuments({last_active_at: {$gte: startDate}}) - 活跃用户
	// 4. db.users.countDocuments({vip_level: {$ne: "none"}}) - VIP用户
	// 5. 留存率计算：需要活跃度记录表
	// 6. 平均活跃天数：需要聚合管道计算
	//
	// 实施策略：Task 3 实现Repository层聚合方法
	// 预计工期：4小时

	stats := &PlatformUserStats{
		TotalUsers:       0,
		NewUsers:         0,
		ActiveUsers:      0,
		VIPUsers:         0,
		RetentionRate:    0,
		AverageActiveDay: 0,
	}

	zap.L().Warn("GetPlatformUserStats: 当前返回空数据，等待Task3实现")

	return stats, nil
}

// ============ 内容统计 ============

// ContentStats 内容统计
type ContentStats struct {
	UserID             string  `json:"user_id"`
	TotalProjects      int64   `json:"total_projects"`        // 总项目数
	PublishedBooks     int64   `json:"published_books"`       // 已发布书籍
	DraftBooks         int64   `json:"draft_books"`           // 草稿书籍
	TotalChapters      int64   `json:"total_chapters"`        // 总章节数
	TotalWords         int64   `json:"total_words"`           // 总字数
	AverageWordsPerDay float64 `json:"average_words_per_day"` // 日均字数
	TotalViews         int64   `json:"total_views"`           // 总浏览量
	TotalCollections   int64   `json:"total_collections"`     // 总收藏数
	AverageRating      float64 `json:"average_rating"`        // 平均评分
}

// GetContentStats 获取内容统计
func (s *PlatformStatsServiceImpl) GetContentStats(ctx context.Context, userID string) (*ContentStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 1. 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 2. 统计项目数
	projectCount, err := s.projectRepo.CountByOwner(ctx, userID)
	if err != nil {
		zap.L().Warn("统计项目数失败，使用默认值0",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		projectCount = 0
	}

	// 3. 获取用户的所有书籍
	// 注意：当前GetByAuthorID需要ObjectID，暂时跳过书籍统计
	// TODO(优化): 扩展BookRepository支持string author_id查询
	publishedBooks := int64(0)
	draftBooks := int64(0)
	totalChapters := int64(0)
	totalWords := int64(0)
	averageWordsPerDay := float64(0)

	zap.L().Debug("书籍统计暂时跳过",
		zap.String("user_id", userID),
		zap.String("reason", "需要ObjectID类型author_id，当前用户ID为string类型"),
	)

	// 6. 构建统计结果
	stats := &ContentStats{
		UserID:             userID,
		TotalProjects:      projectCount,
		PublishedBooks:     publishedBooks,
		DraftBooks:         draftBooks,
		TotalChapters:      totalChapters,
		TotalWords:         totalWords,
		AverageWordsPerDay: averageWordsPerDay,
		TotalViews:         0, // TODO(Task3): 需要阅读统计
		TotalCollections:   0, // TODO(Task3): 需要收藏统计
		AverageRating:      0, // TODO(Task3): 需要评分统计
	}

	return stats, nil
}

// PlatformContentStats 平台内容统计
type PlatformContentStats struct {
	TotalBooks        int64    `json:"total_books"`        // 总书籍数
	NewBooks          int64    `json:"new_books"`          // 新增书籍
	TotalChapters     int64    `json:"total_chapters"`     // 总章节数
	TotalWords        int64    `json:"total_words"`        // 总字数
	TotalViews        int64    `json:"total_views"`        // 总浏览量
	AverageRating     float64  `json:"average_rating"`     // 平均评分
	PopularCategories []string `json:"popular_categories"` // 热门分类
}

// GetPlatformContentStats 获取平台内容统计
func (s *PlatformStatsServiceImpl) GetPlatformContentStats(ctx context.Context, startDate, endDate time.Time) (*PlatformContentStats, error) {
	// TODO(Task3-聚合查询): 实现MongoDB聚合管道统计
	// 需要实现：
	// 1. db.books.countDocuments({}) - 总书籍数
	// 2. db.books.countDocuments({created_at: {$gte: startDate, $lte: endDate}}) - 新增书籍
	// 3. db.chapters.aggregate([{$group: {_id: null, total: {$sum: 1}}}]) - 总章节数
	// 4. db.chapters.aggregate([{$group: {_id: null, total: {$sum: "$word_count"}}}]) - 总字数
	// 5. db.book_stats.aggregate([{$group: {_id: null, total: {$sum: "$view_count"}}}]) - 总浏览量
	// 6. db.book_ratings.aggregate([{$group: {_id: null, avg: {$avg: "$rating"}}}]) - 平均评分
	// 7. db.books.aggregate([{$group: {_id: "$category", count: {$sum: 1}}}, {$sort: {count: -1}}, {$limit: 5}]) - 热门分类
	//
	// 实施策略：Task 3 实现Repository层聚合方法
	// 预计工期：6小时

	stats := &PlatformContentStats{
		TotalBooks:        0,
		NewBooks:          0,
		TotalChapters:     0,
		TotalWords:        0,
		TotalViews:        0,
		AverageRating:     0,
		PopularCategories: []string{},
	}

	zap.L().Warn("GetPlatformContentStats: 当前返回空数据，等待Task3实现")

	return stats, nil
}

// ============ 活跃度统计 ============

// ActivityStats 活跃度统计
type ActivityStats struct {
	UserID       string           `json:"user_id"`
	Days         int              `json:"days"`          // 统计天数
	TotalActions int64            `json:"total_actions"` // 总操作数
	DailyActions []DailyAction    `json:"daily_actions"` // 每日操作
	ActionTypes  map[string]int64 `json:"action_types"`  // 操作类型分布
	ActiveHours  []int            `json:"active_hours"`  // 活跃时段
}

// DailyAction 每日活跃数据
type DailyAction struct {
	Date    string `json:"date"`    // 日期
	Actions int64  `json:"actions"` // 操作数
	Words   int64  `json:"words"`   // 字数
}

// GetUserActivityStats 获取用户活跃度统计
func (s *PlatformStatsServiceImpl) GetUserActivityStats(ctx context.Context, userID string, days int) (*ActivityStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if days <= 0 {
		days = 7 // 默认7天
	}

	// TODO(Task3-活跃度统计): 需要实现活跃度记录表和统计逻辑
	// 需要设计：
	// 1. user_activity_logs表：记录用户操作（写作、阅读、评论、点赞）
	// 2. 聚合查询：按天统计操作数
	// 3. 聚合查询：按操作类型分组统计
	// 4. 聚合查询：按小时统计活跃时段
	//
	// 实施策略：Task 3 设计活跃度记录表，实现聚合查询
	// 预计工期：4小时

	stats := &ActivityStats{
		UserID:       userID,
		Days:         days,
		TotalActions: 0,
		DailyActions: []DailyAction{},
		ActionTypes:  map[string]int64{},
		ActiveHours:  []int{},
	}

	zap.L().Warn("GetUserActivityStats: 当前返回空数据，等待Task3实现",
		zap.String("user_id", userID),
		zap.Int("days", days),
	)

	return stats, nil
}

// ============ 收益统计 ============

// RevenueStats 收益统计
type RevenueStats struct {
	UserID        string             `json:"user_id"`
	TotalRevenue  float64            `json:"total_revenue"`   // 总收益
	PeriodRevenue float64            `json:"period_revenue"`  // 期间收益
	DailyRevenue  []DailyRevenue     `json:"daily_revenue"`   // 每日收益
	RevenueByBook map[string]float64 `json:"revenue_by_book"` // 按书籍收益
	RevenueByType map[string]float64 `json:"revenue_by_type"` // 按类型收益
}

// DailyRevenue 每日收益
type DailyRevenue struct {
	Date    string  `json:"date"`    // 日期
	Revenue float64 `json:"revenue"` // 收益
}

// GetRevenueStats 获取收益统计
func (s *PlatformStatsServiceImpl) GetRevenueStats(ctx context.Context, userID string, startDate, endDate time.Time) (*RevenueStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// TODO(Task3-收益统计): 需要实现钱包交易记录和聚合查询
	// 需要设计：
	// 1. wallet_transactions表：记录所有收入交易
	// 2. 聚合查询：总收益（sum(amount)）
	// 3. 聚合查询：期间收益（sum(amount) WHERE created_at BETWEEN startDate AND endDate）
	// 4. 聚合查询：按日统计收益
	// 5. 聚合查询：按书籍分组统计收益
	// 6. 聚合查询：按收益类型分组统计
	//
	// 实施策略：Task 3 设计交易记录表，实现聚合查询
	// 预计工期：4小时

	stats := &RevenueStats{
		UserID:        userID,
		TotalRevenue:  0,
		PeriodRevenue: 0,
		DailyRevenue:  []DailyRevenue{},
		RevenueByBook: map[string]float64{},
		RevenueByType: map[string]float64{},
	}

	zap.L().Warn("GetRevenueStats: 当前返回空数据，等待Task3实现",
		zap.String("user_id", userID),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	return stats, nil
}

// ============ BaseService接口实现 ============

// Initialize 初始化服务
func (s *PlatformStatsServiceImpl) Initialize(ctx context.Context) error {
	s.initialized = true
	return nil
}

// Health 健康检查
func (s *PlatformStatsServiceImpl) Health(ctx context.Context) error {
	// StatsService没有复杂的初始化逻辑，直接返回成功
	return nil
}

// Close 关闭服务
func (s *PlatformStatsServiceImpl) Close(ctx context.Context) error {
	s.initialized = false
	return nil
}

// GetServiceName 获取服务名称
func (s *PlatformStatsServiceImpl) GetServiceName() string {
	return "PlatformStatsService"
}

// GetVersion 获取服务版本
func (s *PlatformStatsServiceImpl) GetVersion() string {
	return "v1.0.0"
}

// ============ TODO: 高级统计功能（Task 3+） ============
//
// **P0任务2完成情况**：
// ✅ GetUserStats - 实际Repository查询
// ✅ GetContentStats - 实际Repository查询
// ⏸️ GetPlatformUserStats - 延后到Task 3（需要聚合查询）
// ⏸️ GetPlatformContentStats - 延后到Task 3（需要聚合查询）
// ⏸️ GetUserActivityStats - 延后到Task 3（需要活跃度记录表）
// ⏸️ GetRevenueStats - 延后到Task 3（需要钱包交易记录）
//
// **Task 3实施计划**（MongoDB聚合查询）：
//
// 1. **平台级统计** (6小时):
//    - 实现UserRepository.GetPlatformStats() - 聚合查询
//    - 实现BookRepository.GetPlatformStats() - 聚合查询
//    - 实现ChapterRepository.GetPlatformStats() - 聚合查询
//
// 2. **活跃度统计** (4小时):
//    - 设计user_activity_logs表
//    - 实现ActivityRepository及聚合查询
//
// 3. **收益统计** (4小时):
//    - 设计wallet_transactions表
//    - 实现TransactionRepository及聚合查询
//
// 4. **性能优化** (2小时):
//    - Redis缓存层（热门统计缓存1小时）
//    - 异步统计更新（EventBus触发）
//
// **Phase 4+高级功能**（延后）：
// - [ ] 趋势分析（增长率、环比等）
// - [ ] 用户画像分析
// - [ ] 内容质量分析
// - [ ] 收益预测
// - [ ] 数据导出（Excel/PDF）
// - [ ] 自定义统计报表
// - [ ] 数据可视化（图表数据）
