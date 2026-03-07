package bookstore

import (
	"sort"
	"time"

	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"fmt"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// RankingConfig 榜单计算配置
type RankingConfig struct {
	RealtimeViewWeight float64
	RealtimeLikeWeight float64
	WeeklyViewWeight   float64
	MonthlyViewWeight  float64
	MonthlyLikeWeight  float64
	NewbieViewWeight   float64
	NewbieLikeWeight   float64
	NewbieMaxAge       time.Duration
	MaxRankingItems    int
}

func DefaultRankingConfig() *RankingConfig {
	return &RankingConfig{
		RealtimeViewWeight: 0.7,
		RealtimeLikeWeight: 0.3,
		WeeklyViewWeight:   0.6,
		MonthlyViewWeight:  0.5,
		MonthlyLikeWeight:  0.3,
		NewbieViewWeight:   0.6,
		NewbieLikeWeight:   0.4,
		NewbieMaxAge:       30 * 24 * time.Hour,
		MaxRankingItems:    100,
	}
}

type RankingService interface {
	CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error)
	UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error
	GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error)
	GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetConfig() *RankingConfig
	SetConfig(config *RankingConfig)
}

type RankingServiceImpl struct {
	rankingRepo BookstoreRepo.RankingRepository
	bookRepo    BookstoreRepo.BookRepository
	config      *RankingConfig
}

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

func (s *RankingServiceImpl) GetConfig() *RankingConfig {
	return s.config
}

func (s *RankingServiceImpl) SetConfig(config *RankingConfig) {
	if config != nil {
		s.config = config
	}
}

// calculateRanking 通用榜单计算逻辑
func (s *RankingServiceImpl) calculateRanking(ctx context.Context, rankingType bookstore2.RankingType, period string) ([]*bookstore2.RankingItem, error) {
	books, err := s.rankingRepo.GetBooksForRanking(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取书籍数据失败: %w", err)
	}

	now := time.Now()
	items := make([]*bookstore2.RankingItem, 0, len(books))
	for _, book := range books {
		if book == nil || !s.isEligible(book, rankingType, now) {
			continue
		}
		items = append(items, &bookstore2.RankingItem{
			BookID:    book.ID,
			Type:      rankingType,
			Score:     s.calculateScore(book, rankingType),
			ViewCount: book.ViewCount,
			LikeCount: book.RatingCount,
			Period:    period,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].BookID.Hex() < items[j].BookID.Hex()
		}
		return items[i].Score > items[j].Score
	})

	if len(items) > s.config.MaxRankingItems {
		items = items[:s.config.MaxRankingItems]
	}
	for i, item := range items {
		item.Rank = i + 1
	}

	return items, nil
}

// calculateScore 计算书籍分数
func (s *RankingServiceImpl) calculateScore(book *bookstore2.Book, rankingType bookstore2.RankingType) float64 {
	switch rankingType {
	case bookstore2.RankingTypeRealtime:
		return float64(book.ViewCount)*s.config.RealtimeViewWeight + float64(book.RatingCount)*s.config.RealtimeLikeWeight
	case bookstore2.RankingTypeWeekly:
		return float64(book.ViewCount) * s.config.WeeklyViewWeight
	case bookstore2.RankingTypeMonthly:
		return float64(book.ViewCount)*s.config.MonthlyViewWeight + float64(book.RatingCount)*s.config.MonthlyLikeWeight
	case bookstore2.RankingTypeNewbie:
		return float64(book.ViewCount)*s.config.NewbieViewWeight + float64(book.RatingCount)*s.config.NewbieLikeWeight
	default:
		return 0
	}
}

// isEligible 判断书籍是否符合榜单资格
func (s *RankingServiceImpl) isEligible(book *bookstore2.Book, rankingType bookstore2.RankingType, now time.Time) bool {
	if book.Status != bookstore2.BookStatusOngoing {
		return false
	}
	if rankingType != bookstore2.RankingTypeNewbie {
		return true
	}
	if book.PublishedAt == nil {
		return false
	}
	return now.Sub(*book.PublishedAt) <= s.config.NewbieMaxAge
}

func (s *RankingServiceImpl) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	return s.calculateRanking(ctx, bookstore2.RankingTypeRealtime, period)
}

func (s *RankingServiceImpl) CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	return s.calculateRanking(ctx, bookstore2.RankingTypeWeekly, period)
}

func (s *RankingServiceImpl) CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	return s.calculateRanking(ctx, bookstore2.RankingTypeMonthly, period)
}

func (s *RankingServiceImpl) CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	return s.calculateRanking(ctx, bookstore2.RankingTypeNewbie, period)
}

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

	return s.rankingRepo.Transaction(ctx, func(txCtx context.Context) error {
		return s.rankingRepo.UpdateRankings(txCtx, rankingType, period, items)
	})
}

func (s *RankingServiceImpl) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error) {
	period := bookstore2.GetPeriodString(bookstore2.RankingTypeRealtime, time.Now())
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeRealtime, period, limit, 0)
}

func (s *RankingServiceImpl) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeWeekly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeWeekly, period, limit, 0)
}

func (s *RankingServiceImpl) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeMonthly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeMonthly, period, limit, 0)
}

func (s *RankingServiceImpl) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeNewbie, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeNewbie, period, limit, 0)
}

func (s *RankingServiceImpl) GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(rankingType, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, rankingType, period, limit, 0)
}
