package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"fmt"
	"time"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// RankingConfig 榜单计算配置
// 将业务规则（权重配置）从 Repository 层移到 Service 层
type RankingConfig struct {
	// 实时榜权重
	RealtimeViewWeight float64 // 浏览量权重
	RealtimeLikeWeight float64 // 点赞数权重

	// 周榜权重
	WeeklyViewWeight float64 // 浏览量权重

	// 月榜权重
	MonthlyViewWeight float64 // 浏览量权重
	MonthlyLikeWeight float64 // 点赞数权重

	// 新人榜权重
	NewbieViewWeight   float64 // 浏览量权重
	NewbieLikeWeight   float64 // 点赞数权重
	NewbieRankingLimit int     // 榜单数量限制

	// 通用配置
	MaxRankingItems int // 榜单最大数量
}

// DefaultRankingConfig 返回默认榜单配置
func DefaultRankingConfig() *RankingConfig {
	return &RankingConfig{
		// 实时榜：浏览量70%，点赞30%
		RealtimeViewWeight: 0.7,
		RealtimeLikeWeight: 0.3,

		// 周榜：浏览量60%
		WeeklyViewWeight: 0.6,

		// 月榜：浏览量50%，点赞30%
		MonthlyViewWeight: 0.5,
		MonthlyLikeWeight: 0.3,

		// 新人榜：浏览量60%，点赞40%
		NewbieViewWeight:   0.6,
		NewbieLikeWeight:   0.4,
		NewbieRankingLimit: 100,

		// 通用
		MaxRankingItems: 100,
	}
}

// RankingService 榜单服务接口
type RankingService interface {
	// 榜单计算方法（业务逻辑）
	CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)

	// 榜单更新（编排计算和持久化）
	UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error

	// 榜单查询（代理到 Repository）
	GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error)
	GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error)

	// 配置管理
	GetConfig() *RankingConfig
	SetConfig(config *RankingConfig)
}

// RankingServiceImpl 榜单服务实现
type RankingServiceImpl struct {
	rankingRepo BookstoreRepo.RankingRepository
	bookRepo    BookstoreRepo.BookRepository
	config      *RankingConfig
}

// NewRankingService 创建榜单服务
func NewRankingService(
	rankingRepo BookstoreRepo.RankingRepository,
	bookRepo BookstoreRepo.BookRepository,
	config *RankingConfig,
) RankingService {
	if config == nil {
		config = DefaultRankingConfig()
	}
	return &RankingServiceImpl{
		rankingRepo: rankingRepo,
		bookRepo:    bookRepo,
		config:      config,
	}
}

// GetConfig 获取当前配置
func (s *RankingServiceImpl) GetConfig() *RankingConfig {
	return s.config
}

// SetConfig 设置配置
func (s *RankingServiceImpl) SetConfig(config *RankingConfig) {
	if config != nil {
		s.config = config
	}
}

// CalculateRealtimeRanking 计算实时榜单
// 业务逻辑：基于浏览量和点赞数，使用配置的权重计算热度分数
func (s *RankingServiceImpl) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	// 调用 Repository 的数据查询方法，但使用 Service 层的权重配置
	items, err := s.rankingRepo.CalculateRealtimeRanking(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate realtime ranking: %w", err)
	}

	// 应用 Service 层的权重配置重新计算分数
	for _, item := range items {
		item.Score = float64(item.ViewCount)*s.config.RealtimeViewWeight +
			float64(item.LikeCount)*s.config.RealtimeLikeWeight
	}

	// 限制数量
	if len(items) > s.config.MaxRankingItems {
		items = items[:s.config.MaxRankingItems]
	}

	return items, nil
}

// CalculateWeeklyRanking 计算周榜
// 业务逻辑：基于浏览量，使用配置的权重计算周榜分数
// 注：当前实现仅使用浏览量，章节数权重保留配置以便未来扩展
func (s *RankingServiceImpl) CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	items, err := s.rankingRepo.CalculateWeeklyRanking(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate weekly ranking: %w", err)
	}

	// 应用 Service 层的权重配置重新计算分数
	// 当前仅使用浏览量权重
	for _, item := range items {
		item.Score = float64(item.ViewCount) * s.config.WeeklyViewWeight
	}

	if len(items) > s.config.MaxRankingItems {
		items = items[:s.config.MaxRankingItems]
	}

	return items, nil
}

// CalculateMonthlyRanking 计算月榜
// 业务逻辑：基于浏览量和点赞数，使用配置的权重计算月榜分数
// 注：当前实现使用浏览量和点赞数，字数权重保留配置以便未来扩展
func (s *RankingServiceImpl) CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	items, err := s.rankingRepo.CalculateMonthlyRanking(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate monthly ranking: %w", err)
	}

	// 应用 Service 层的权重配置重新计算分数
	for _, item := range items {
		item.Score = float64(item.ViewCount)*s.config.MonthlyViewWeight +
			float64(item.LikeCount)*s.config.MonthlyLikeWeight
	}

	if len(items) > s.config.MaxRankingItems {
		items = items[:s.config.MaxRankingItems]
	}

	return items, nil
}

// CalculateNewbieRanking 计算新人榜
// 业务逻辑：筛选新书（配置的最大年龄内），基于浏览量和点赞数计算分数
func (s *RankingServiceImpl) CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	items, err := s.rankingRepo.CalculateNewbieRanking(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate newbie ranking: %w", err)
	}

	// 应用 Service 层的权重配置重新计算分数
	for _, item := range items {
		item.Score = float64(item.ViewCount)*s.config.NewbieViewWeight +
			float64(item.LikeCount)*s.config.NewbieLikeWeight
	}

	// 限制新人榜数量
	if len(items) > s.config.NewbieRankingLimit {
		items = items[:s.config.NewbieRankingLimit]
	}

	return items, nil
}

// UpdateRankings 更新榜单数据
// 编排计算和持久化操作
func (s *RankingServiceImpl) UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error {
	var items []*bookstore2.RankingItem
	var err error

	switch rankingType {
	case bookstore2.RankingTypeRealtime:
		items, err = s.CalculateRealtimeRanking(ctx, period)
	case bookstore2.RankingTypeWeekly:
		items, err = s.CalculateWeeklyRanking(ctx, period)
	case bookstore2.RankingTypeMonthly:
		items, err = s.CalculateMonthlyRanking(ctx, period)
	case bookstore2.RankingTypeNewbie:
		items, err = s.CalculateNewbieRanking(ctx, period)
	default:
		return fmt.Errorf("unsupported ranking type: %s", rankingType)
	}

	if err != nil {
		return fmt.Errorf("failed to calculate ranking: %w", err)
	}

	// 使用事务更新榜单
	return s.rankingRepo.Transaction(ctx, func(txCtx context.Context) error {
		return s.rankingRepo.UpdateRankings(txCtx, rankingType, period, items)
	})
}

// GetRealtimeRanking 获取实时榜
func (s *RankingServiceImpl) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error) {
	period := bookstore2.GetPeriodString(bookstore2.RankingTypeRealtime, time.Now())
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeRealtime, period, limit, 0)
}

// GetWeeklyRanking 获取周榜
func (s *RankingServiceImpl) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeWeekly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeWeekly, period, limit, 0)
}

// GetMonthlyRanking 获取月榜
func (s *RankingServiceImpl) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeMonthly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeMonthly, period, limit, 0)
}

// GetNewbieRanking 获取新人榜
func (s *RankingServiceImpl) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeNewbie, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeNewbie, period, limit, 0)
}

// GetRankingByType 根据类型获取榜单
func (s *RankingServiceImpl) GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(rankingType, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, rankingType, period, limit, 0)
}
