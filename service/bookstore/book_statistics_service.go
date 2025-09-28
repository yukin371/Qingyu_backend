package bookstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// BookStatisticsService 书籍统计服务接口
type BookStatisticsService interface {
	// 统计基础操作
	CreateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error
	GetStatisticsByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookStatistics, error)
	GetStatisticsByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstore.BookStatistics, error)
	UpdateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error
	DeleteStatistics(ctx context.Context, id primitive.ObjectID) error

	// 统计数据查询
	GetTopViewedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error)
	GetTopFavoritedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error)
	GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error)
	GetHottestBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error)
	GetTrendingBooks(ctx context.Context, days int, limit int) ([]*bookstore.BookStatistics, error)

	// 统计数据更新
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error

	// 评分统计更新
	UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error
	RemoveRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error
	RecalculateRating(ctx context.Context, bookID primitive.ObjectID) error

	// 热度分数管理
	UpdateHotScore(ctx context.Context, bookID primitive.ObjectID) error
	BatchUpdateHotScore(ctx context.Context, bookIDs []primitive.ObjectID) error
	RecalculateAllHotScores(ctx context.Context) error

	// 聚合统计
	GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error)
	GetStatisticsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*bookstore.BookStatistics, error)
	GetBookPopularityLevel(ctx context.Context, bookID primitive.ObjectID) (string, error)

	// 批量操作
	BatchUpdateViewCount(ctx context.Context, bookIDs []primitive.ObjectID, increment int64) error
	BatchCreateStatistics(ctx context.Context, statsList []*bookstore.BookStatistics) error
	BatchDeleteStatistics(ctx context.Context, bookIDs []primitive.ObjectID) error

	// 统计报告
	GenerateDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error)
	GenerateWeeklyReport(ctx context.Context, startDate time.Time) (map[string]interface{}, error)
	GenerateMonthlyReport(ctx context.Context, year int, month int) (map[string]interface{}, error)

	// 搜索和过滤
	SearchStatistics(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookStatistics, int64, error)
	GetStatisticsByFilter(ctx context.Context, filter *BookstoreRepo.BookStatisticsFilter, page, pageSize int) ([]*bookstore.BookStatistics, int64, error)
}

// BookStatisticsServiceImpl 书籍统计服务实现
type BookStatisticsServiceImpl struct {
	statsRepo    BookstoreRepo.BookStatisticsRepository
	cacheService CacheService
}

// NewBookStatisticsService 创建书籍统计服务实例
func NewBookStatisticsService(statsRepo BookstoreRepo.BookStatisticsRepository, cacheService CacheService) BookStatisticsService {
	return &BookStatisticsServiceImpl{
		statsRepo:    statsRepo,
		cacheService: cacheService,
	}
}

// CreateStatistics 创建统计数据
func (s *BookStatisticsServiceImpl) CreateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error {
	if stats == nil {
		return errors.New("statistics cannot be nil")
	}

	// 验证必填字段
	if stats.BookID.IsZero() {
		return errors.New("book ID is required")
	}

	// 检查是否已存在
	existingStats, err := s.statsRepo.GetByBookID(ctx, stats.BookID)
	if err != nil {
		return fmt.Errorf("failed to check existing statistics: %w", err)
	}
	if existingStats != nil {
		return errors.New("statistics for this book already exists")
	}

	// 创建统计数据
	if err := s.statsRepo.Create(ctx, stats); err != nil {
		return fmt.Errorf("failed to create statistics: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, stats)

	return nil
}

// GetStatisticsByID 根据ID获取统计数据
func (s *BookStatisticsServiceImpl) GetStatisticsByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookStatistics, error) {
	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedStats, err := s.cacheService.GetBookStatistics(ctx, id.Hex()); err == nil && cachedStats != nil {
			return cachedStats, nil
		}
	}

	// 从数据库获取
	stats, err := s.statsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	if stats == nil {
		return nil, errors.New("statistics not found")
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetBookStatistics(ctx, id.Hex(), stats, 10*time.Minute)
	}

	return stats, nil
}

// GetStatisticsByBookID 根据书籍ID获取统计数据
func (s *BookStatisticsServiceImpl) GetStatisticsByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstore.BookStatistics, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedStats, err := s.cacheService.GetBookStatistics(ctx, bookID.Hex()); err == nil && cachedStats != nil {
			return cachedStats, nil
		}
	}

	// 从数据库获取
	stats, err := s.statsRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by book ID: %w", err)
	}

	// 如果不存在，创建默认统计数据
	if stats == nil {
		stats = &bookstore.BookStatistics{
			BookID:             bookID,
			ViewCount:          0,
			FavoriteCount:      0,
			CommentCount:       0,
			ShareCount:         0,
			AverageRating:      0,
			RatingCount:        0,
			RatingDistribution: make(map[int]int64),
			HotScore:           0,
		}
		if err := s.statsRepo.Create(ctx, stats); err != nil {
			return nil, fmt.Errorf("failed to create default statistics: %w", err)
		}
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetBookStatistics(ctx, bookID.Hex(), stats, 10*time.Minute)
	}

	return stats, nil
}

// UpdateStatistics 更新统计数据
func (s *BookStatisticsServiceImpl) UpdateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error {
	if stats == nil {
		return errors.New("statistics cannot be nil")
	}

	// 验证必填字段
	if stats.BookID.IsZero() {
		return errors.New("book ID is required")
	}

	// 更新统计数据
	if err := s.statsRepo.Update(ctx, stats.ID, map[string]interface{}{
		"view_count":     stats.ViewCount,
		"favorite_count": stats.FavoriteCount,
		"comment_count":  stats.CommentCount,
		"share_count":    stats.ShareCount,
		"average_rating": stats.AverageRating,
		"rating_count":   stats.RatingCount,
		"hot_score":      stats.HotScore,
		"updated_at":     time.Now(),
	}); err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, stats)

	return nil
}

// DeleteStatistics 删除统计数据
func (s *BookStatisticsServiceImpl) DeleteStatistics(ctx context.Context, id primitive.ObjectID) error {
	// 先获取统计数据用于清除缓存
	stats, err := s.statsRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get statistics for deletion: %w", err)
	}
	if stats == nil {
		return errors.New("statistics not found")
	}

	// 删除统计数据
	if err := s.statsRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete statistics: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, stats)

	return nil
}

// GetTopViewedBooks 获取浏览量最高的书籍
func (s *BookStatisticsServiceImpl) GetTopViewedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedBooks, err := s.cacheService.GetTopViewedBooks(ctx); err == nil && len(cachedBooks) > 0 {
			// 将 []*bookstore.Book 转换为 []*bookstore.BookStatistics
			var result []*bookstore.BookStatistics
			for _, book := range cachedBooks {
				if len(result) >= limit {
					break
				}
				// 这里需要根据实际情况转换，暂时返回空的统计数据
				stats := &bookstore.BookStatistics{
					BookID: book.ID,
				}
				result = append(result, stats)
			}
			return result, nil
		}
	}

	// 从数据库获取
	books, err := s.statsRepo.GetTopViewed(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get top viewed books: %w", err)
	}

	// 缓存结果 - 需要将统计数据转换为书籍数据进行缓存
	if s.cacheService != nil {
		var booksForCache []*bookstore.Book
		for _, stats := range books {
			// 这里需要根据实际情况构造Book对象，暂时使用基本信息
			book := &bookstore.Book{
				ID: stats.BookID,
			}
			booksForCache = append(booksForCache, book)
		}
		s.cacheService.SetTopViewedBooks(ctx, booksForCache, 5*time.Minute)
	}

	return books, nil
}

// GetTopFavoritedBooks 获取收藏量最高的书籍
func (s *BookStatisticsServiceImpl) GetTopFavoritedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedBooks, err := s.cacheService.GetTopFavoritedBooks(ctx); err == nil && len(cachedBooks) > 0 {
			// 将 []*bookstore.Book 转换为 []*bookstore.BookStatistics
			var result []*bookstore.BookStatistics
			for _, book := range cachedBooks {
				if len(result) >= limit {
					break
				}
				// 这里需要根据实际情况转换，暂时返回空的统计数据
				stats := &bookstore.BookStatistics{
					BookID: book.ID,
				}
				result = append(result, stats)
			}
			return result, nil
		}
	}

	// 从数据库获取
	books, err := s.statsRepo.GetTopFavorited(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get top favorited books: %w", err)
	}

	// 缓存结果 - 需要将统计数据转换为书籍数据进行缓存
	if s.cacheService != nil {
		var booksForCache []*bookstore.Book
		for _, stats := range books {
			// 这里需要根据实际情况构造Book对象，暂时使用基本信息
			book := &bookstore.Book{
				ID: stats.BookID,
			}
			booksForCache = append(booksForCache, book)
		}
		s.cacheService.SetTopFavoritedBooks(ctx, booksForCache, 5*time.Minute)
	}

	return books, nil
}

// GetTopRatedBooks 获取评分最高的书籍
func (s *BookStatisticsServiceImpl) GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedBooks, err := s.cacheService.GetTopRatedBooks(ctx); err == nil && len(cachedBooks) > 0 {
			// 将 []*bookstore.Book 转换为 []*bookstore.BookStatistics
			var result []*bookstore.BookStatistics
			for _, book := range cachedBooks {
				if len(result) >= limit {
					break
				}
				// 这里需要根据实际情况转换，暂时返回空的统计数据
				stats := &bookstore.BookStatistics{
					BookID: book.ID,
				}
				result = append(result, stats)
			}
			return result, nil
		}
	}

	// 从数据库获取
	books, err := s.statsRepo.GetTopRated(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated books: %w", err)
	}

	// 缓存结果 - 需要将统计数据转换为书籍数据进行缓存
	if s.cacheService != nil {
		var booksForCache []*bookstore.Book
		for _, stats := range books {
			// 这里需要根据实际情况构造Book对象，暂时使用基本信息
			book := &bookstore.Book{
				ID: stats.BookID,
			}
			booksForCache = append(booksForCache, book)
		}
		s.cacheService.SetTopRatedBooks(ctx, booksForCache, 5*time.Minute)
	}

	return books, nil
}

// GetHottestBooks 获取最热门的书籍
func (s *BookStatisticsServiceImpl) GetHottestBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedBooks, err := s.cacheService.GetHottestBooks(ctx); err == nil && len(cachedBooks) > 0 {
			// 将 []*bookstore.Book 转换为 []*bookstore.BookStatistics
			var result []*bookstore.BookStatistics
			for _, book := range cachedBooks {
				if len(result) >= limit {
					break
				}
				// 这里需要根据实际情况转换，暂时返回空的统计数据
				stats := &bookstore.BookStatistics{
					BookID: book.ID,
				}
				result = append(result, stats)
			}
			return result, nil
		}
	}

	// 从数据库获取
	books, err := s.statsRepo.GetHottest(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get hottest books: %w", err)
	}

	// 缓存结果 - 需要将统计数据转换为书籍数据进行缓存
	if s.cacheService != nil {
		var booksForCache []*bookstore.Book
		for _, stats := range books {
			// 这里需要根据实际情况构造Book对象，暂时使用基本信息
			book := &bookstore.Book{
				ID: stats.BookID,
			}
			booksForCache = append(booksForCache, book)
		}
		s.cacheService.SetHottestBooks(ctx, booksForCache, 5*time.Minute)
	}

	return books, nil
}

// GetTrendingBooks 获取趋势书籍
func (s *BookStatisticsServiceImpl) GetTrendingBooks(ctx context.Context, days int, limit int) ([]*bookstore.BookStatistics, error) {
	if days < 1 || days > 30 {
		days = 7
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 计算时间范围

	// 从数据库获取
	books, err := s.statsRepo.GetTrendingBooks(ctx, days, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending books: %w", err)
	}

	return books, nil
}

// IncrementViewCount 增加浏览量
func (s *BookStatisticsServiceImpl) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 增加浏览量
	if err := s.statsRepo.IncrementViewCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		// 热度分数更新失败不影响主要操作
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopViewedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// IncrementFavoriteCount 增加收藏量
func (s *BookStatisticsServiceImpl) IncrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 增加收藏量
	if err := s.statsRepo.IncrementFavoriteCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment favorite count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopFavoritedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// DecrementFavoriteCount 减少收藏量
func (s *BookStatisticsServiceImpl) DecrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 减少收藏量
	if err := s.statsRepo.DecrementFavoriteCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to decrement favorite count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopFavoritedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// IncrementCommentCount 增加评论量
func (s *BookStatisticsServiceImpl) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 增加评论量
	if err := s.statsRepo.IncrementCommentCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment comment count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// DecrementCommentCount 减少评论量
func (s *BookStatisticsServiceImpl) DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 减少评论量
	if err := s.statsRepo.DecrementCommentCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to decrement comment count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// IncrementShareCount 增加分享量
func (s *BookStatisticsServiceImpl) IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 增加分享量
	if err := s.statsRepo.IncrementShareCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment share count: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// UpdateRating 更新评分统计
func (s *BookStatisticsServiceImpl) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// 更新评分统计 - 将float64转换为int
	if err := s.statsRepo.UpdateRating(ctx, bookID, int(rating)); err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopRatedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// RemoveRating 移除评分统计
func (s *BookStatisticsServiceImpl) RemoveRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// 移除评分统计
	if err := s.statsRepo.RemoveRating(ctx, bookID, int(rating)); err != nil {
		return fmt.Errorf("failed to remove rating: %w", err)
	}

	// 更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		fmt.Printf("failed to update hot score for book %s: %v\n", bookID.Hex(), err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopRatedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// RecalculateRating 重新计算评分统计
func (s *BookStatisticsServiceImpl) RecalculateRating(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// TODO: 这里应该从评分表重新计算平均评分和分布
	// 简化处理，直接更新热度分数
	if err := s.UpdateHotScore(ctx, bookID); err != nil {
		return fmt.Errorf("failed to update hot score: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		s.cacheService.InvalidateTopRatedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// UpdateHotScore 更新热度分数
func (s *BookStatisticsServiceImpl) UpdateHotScore(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	// 获取当前统计数据
	stats, err := s.statsRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}
	if stats == nil {
		return errors.New("statistics not found")
	}

	// 计算热度分数
	stats.CalculateHotScore()

	// 更新热度分数
	if err := s.statsRepo.UpdateHotScore(ctx, bookID, stats.HotScore); err != nil {
		return fmt.Errorf("failed to update hot score: %w", err)
	}

	return nil
}

// BatchUpdateHotScore 批量更新热度分数
func (s *BookStatisticsServiceImpl) BatchUpdateHotScore(ctx context.Context, bookIDs []primitive.ObjectID) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}

	// 批量更新热度分数
	if err := s.statsRepo.BatchUpdateHotScore(ctx, bookIDs); err != nil {
		return fmt.Errorf("failed to batch update hot score: %w", err)
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		}
	}
	if s.cacheService != nil {
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// RecalculateAllHotScores 重新计算所有热度分数
func (s *BookStatisticsServiceImpl) RecalculateAllHotScores(ctx context.Context) error {
	// TODO: 这里应该分批处理所有书籍的热度分数
	// 简化处理，清除所有相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// GetAggregatedStatistics 获取聚合统计数据
func (s *BookStatisticsServiceImpl) GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error) {
	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedStats, err := s.cacheService.GetAggregatedStatistics(ctx); err == nil && cachedStats != nil {
			return cachedStats, nil
		}
	}

	// 从数据库获取
	stats, err := s.statsRepo.GetAggregatedStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregated statistics: %w", err)
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetAggregatedStatistics(ctx, stats, 5*time.Minute)
	}

	return stats, nil
}

// GetStatisticsByTimeRange 根据时间范围获取统计数据
func (s *BookStatisticsServiceImpl) GetStatisticsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*bookstore.BookStatistics, error) {
	if startTime.After(endTime) {
		return nil, errors.New("start time cannot be after end time")
	}

	stats, err := s.statsRepo.GetStatisticsByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by time range: %w", err)
	}

	return stats, nil
}

// GetBookPopularityLevel 获取书籍热度等级
func (s *BookStatisticsServiceImpl) GetBookPopularityLevel(ctx context.Context, bookID primitive.ObjectID) (string, error) {
	if bookID.IsZero() {
		return "", errors.New("book ID cannot be empty")
	}

	// 获取统计数据
	stats, err := s.statsRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return "", fmt.Errorf("failed to get statistics: %w", err)
	}
	if stats == nil {
		return "unknown", nil
	}

	// 获取热度等级
	level := stats.GetPopularityLevel()
	return level, nil
}

// BatchUpdateViewCount 批量更新浏览量
func (s *BookStatisticsServiceImpl) BatchUpdateViewCount(ctx context.Context, bookIDs []primitive.ObjectID, increment int64) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}
	if increment <= 0 {
		return errors.New("increment must be positive")
	}

	// 批量更新浏览量
	if err := s.statsRepo.BatchUpdateViewCount(ctx, bookIDs, increment); err != nil {
		return fmt.Errorf("failed to batch update view count: %w", err)
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		}
	}
	if s.cacheService != nil {
		s.cacheService.InvalidateTopViewedBooksCache(ctx)
		s.cacheService.InvalidateHottestBooksCache(ctx)
	}

	return nil
}

// BatchCreateStatistics 批量创建统计数据
func (s *BookStatisticsServiceImpl) BatchCreateStatistics(ctx context.Context, statsList []*bookstore.BookStatistics) error {
	if len(statsList) == 0 {
		return errors.New("statistics list cannot be empty")
	}

	// 验证数据
	for _, stats := range statsList {
		if stats.BookID.IsZero() {
			return errors.New("book ID is required for all statistics")
		}
	}

	// TODO: 实现批量创建逻辑
	// 这里简化处理，逐个创建
	for _, stats := range statsList {
		if err := s.statsRepo.Create(ctx, stats); err != nil {
			return fmt.Errorf("failed to create statistics for book %s: %w", stats.BookID.Hex(), err)
		}
	}

	return nil
}

// BatchDeleteStatistics 批量删除统计数据
func (s *BookStatisticsServiceImpl) BatchDeleteStatistics(ctx context.Context, bookIDs []primitive.ObjectID) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}

	// TODO: 实现批量删除逻辑
	// 这里简化处理，逐个删除
	for _, bookID := range bookIDs {
		stats, err := s.statsRepo.GetByBookID(ctx, bookID)
		if err != nil {
			continue // 忽略获取失败的情况
		}
		if stats != nil {
			if err := s.statsRepo.Delete(ctx, stats.ID); err != nil {
				return fmt.Errorf("failed to delete statistics for book %s: %w", bookID.Hex(), err)
			}
		}
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookStatisticsCache(ctx, bookID.Hex())
		}
	}

	return nil
}

// GenerateDailyReport 生成日报
func (s *BookStatisticsServiceImpl) GenerateDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.AddDate(0, 0, 1)

	// 获取时间范围内的统计数据
	stats, err := s.statsRepo.GetStatisticsByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for daily report: %w", err)
	}

	// 生成报告
	report := make(map[string]interface{})
	report["date"] = date.Format("2006-01-02")
	report["total_books"] = len(stats)

	var totalViews, totalFavorites, totalComments, totalShares int64
	for _, stat := range stats {
		totalViews += stat.ViewCount
		totalFavorites += stat.FavoriteCount
		totalComments += stat.CommentCount
		totalShares += stat.ShareCount
	}

	report["total_views"] = totalViews
	report["total_favorites"] = totalFavorites
	report["total_comments"] = totalComments
	report["total_shares"] = totalShares

	return report, nil
}

// GenerateWeeklyReport 生成周报
func (s *BookStatisticsServiceImpl) GenerateWeeklyReport(ctx context.Context, startDate time.Time) (map[string]interface{}, error) {
	endDate := startDate.AddDate(0, 0, 7)

	// 获取时间范围内的统计数据
	stats, err := s.statsRepo.GetStatisticsByTimeRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for weekly report: %w", err)
	}

	// 生成报告
	report := make(map[string]interface{})
	report["start_date"] = startDate.Format("2006-01-02")
	report["end_date"] = endDate.Format("2006-01-02")
	report["total_books"] = len(stats)

	var totalViews, totalFavorites, totalComments, totalShares int64
	for _, stat := range stats {
		totalViews += stat.ViewCount
		totalFavorites += stat.FavoriteCount
		totalComments += stat.CommentCount
		totalShares += stat.ShareCount
	}

	report["total_views"] = totalViews
	report["total_favorites"] = totalFavorites
	report["total_comments"] = totalComments
	report["total_shares"] = totalShares

	return report, nil
}

// GenerateMonthlyReport 生成月报
func (s *BookStatisticsServiceImpl) GenerateMonthlyReport(ctx context.Context, year int, month int) (map[string]interface{}, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// 获取时间范围内的统计数据
	stats, err := s.statsRepo.GetStatisticsByTimeRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for monthly report: %w", err)
	}

	// 生成报告
	report := make(map[string]interface{})
	report["year"] = year
	report["month"] = month
	report["total_books"] = len(stats)

	var totalViews, totalFavorites, totalComments, totalShares int64
	for _, stat := range stats {
		totalViews += stat.ViewCount
		totalFavorites += stat.FavoriteCount
		totalComments += stat.CommentCount
		totalShares += stat.ShareCount
	}

	report["total_views"] = totalViews
	report["total_favorites"] = totalFavorites
	report["total_comments"] = totalComments
	report["total_shares"] = totalShares

	return report, nil
}

// SearchStatistics 搜索统计数据
func (s *BookStatisticsServiceImpl) SearchStatistics(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("keyword cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 搜索统计数据
	stats, total, err := s.statsRepo.Search(ctx, keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search statistics: %w", err)
	}

	return stats, total, nil
}

// GetStatisticsByFilter 根据过滤条件获取统计数据
func (s *BookStatisticsServiceImpl) GetStatisticsByFilter(ctx context.Context, filter *BookstoreRepo.BookStatisticsFilter, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	if filter == nil {
		return nil, 0, errors.New("filter cannot be nil")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 根据过滤条件获取统计数据
	stats, total, err := s.statsRepo.SearchByFilter(ctx, filter, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get statistics by filter: %w", err)
	}

	return stats, total, nil
}

// invalidateRelatedCache 清除相关缓存
func (s *BookStatisticsServiceImpl) invalidateRelatedCache(ctx context.Context, stats *bookstore.BookStatistics) {
	if s.cacheService == nil {
		return
	}

	// 清除统计数据缓存
	s.cacheService.InvalidateBookStatisticsCache(ctx, stats.BookID.Hex())

	// 清除排行榜缓存
	s.cacheService.InvalidateTopViewedBooksCache(ctx)
	s.cacheService.InvalidateTopFavoritedBooksCache(ctx)
	s.cacheService.InvalidateTopRatedBooksCache(ctx)
	s.cacheService.InvalidateHottestBooksCache(ctx)

	// 清除聚合统计缓存
	s.cacheService.InvalidateAggregatedStatisticsCache(ctx)

	// 清除书籍详情缓存（因为统计数据可能影响显示）
	s.cacheService.InvalidateBookDetailCache(ctx, stats.BookID.Hex())
}
