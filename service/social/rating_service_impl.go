package social

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"Qingyu_backend/models/social"
	"Qingyu_backend/pkg/cache"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"go.uber.org/zap"
)

const (
	// BaseCacheTTL 基础缓存TTL
	BaseCacheTTL = 5 * time.Minute
	// EmptyCacheValue 空值缓存标记
	EmptyCacheValue = "EMPTY"
	// EmptyCacheTTL 空值缓存TTL（防止缓存穿透）
	EmptyCacheTTL = 1 * time.Minute
	// TTLJitterPercent TTL抖动百分比（防止缓存雪崩）
	TTLJitterPercent = 0.1 // 10%
)

// init 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RatingServiceImplementation 评分服务实现
type RatingServiceImplementation struct {
	commentRepo socialRepo.CommentRepository
	reviewRepo  socialRepo.ReviewRepository
	redisClient cache.RedisClient
	logger      *zap.Logger
}

// NewRatingService 创建评分服务实例
func NewRatingService(
	commentRepo socialRepo.CommentRepository,
	reviewRepo socialRepo.ReviewRepository,
	redisClient cache.RedisClient,
	logger *zap.Logger,
) RatingService {
	return &RatingServiceImplementation{
		commentRepo: commentRepo,
		reviewRepo:  reviewRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetRatingStats 获取评分统计（带高级缓存策略）
func (s *RatingServiceImplementation) GetRatingStats(ctx context.Context, targetType, targetID string) (*social.RatingStats, error) {
	// 1. 尝试从缓存获取
	cacheKey := fmt.Sprintf("rating:stats:%s:%s", targetType, targetID)
	cached, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		// 检查是否为空值缓存标记
		if cached == EmptyCacheValue {
			s.logger.Debug("命中空值缓存")
			return nil, fmt.Errorf("评分数据不存在")
		}

		stats, err := s.deserializeStats(cached)
		if err == nil {
			s.logger.Debug("缓存命中")
			return stats, nil
		}
		s.logger.Warn("缓存反序列化失败，将重新计算",
			zap.Error(err))
	}

	// 2. 缓存未命中，聚合数据库数据
	stats, err := s.AggregateRatings(ctx, targetType, targetID)
	if err != nil {
		// 3. 缓存空值（防止缓存穿透）
		if s.redisClient != nil {
			if setErr := s.redisClient.Set(ctx, cacheKey, EmptyCacheValue, EmptyCacheTTL); setErr != nil {
				s.logger.Warn("写入空值缓存失败", zap.Error(setErr))
			} else {
				s.logger.Debug("已缓存空值", zap.Duration("ttl", EmptyCacheTTL))
			}
		}
		return nil, fmt.Errorf("聚合评分失败: %w", err)
	}

	// 4. 写入缓存（TTL带随机抖动，防止缓存雪崩）
	serialized, _ := s.serializeStats(stats)
	if s.redisClient != nil {
		ttl := s.calculateTTLWithJitter(BaseCacheTTL)
		if err := s.redisClient.Set(ctx, cacheKey, serialized, ttl); err != nil {
			s.logger.Warn("写入缓存失败", zap.Error(err))
		} else {
			s.logger.Debug("已写入缓存", zap.Duration("ttl", ttl))
		}
	}

	return stats, nil
}

// AggregateRatings 从现有数据源聚合评分
func (s *RatingServiceImplementation) AggregateRatings(ctx context.Context, targetType, targetID string) (*social.RatingStats, error) {
	switch targetType {
	case "comment":
		return s.aggregateCommentRatings(ctx, targetID)
	case "review":
		return s.aggregateReviewRatings(ctx, targetID)
	case "book":
		// 可以扩展支持书籍级别的评分聚合
		return s.aggregateBookRatings(ctx, targetID)
	default:
		return nil, fmt.Errorf("不支持的目标类型: %s", targetType)
	}
}

// aggregateCommentRatings 聚合评论评分
func (s *RatingServiceImplementation) aggregateCommentRatings(ctx context.Context, commentID string) (*social.RatingStats, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("获取评论失败: %w", err)
	}

	return &social.RatingStats{
		TargetID:      commentID,
		TargetType:    "comment",
		AverageRating: float64(comment.Rating),
		TotalRatings:  1,
		Distribution:  map[int]int64{comment.Rating: 1},
		UpdatedAt:     time.Now(),
	}, nil
}

// aggregateReviewRatings 聚合书评评分
func (s *RatingServiceImplementation) aggregateReviewRatings(ctx context.Context, reviewID string) (*social.RatingStats, error) {
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("获取书评失败: %w", err)
	}

	return &social.RatingStats{
		TargetID:      reviewID,
		TargetType:    "review",
		AverageRating: float64(review.Rating),
		TotalRatings:  1,
		Distribution:  map[int]int64{review.Rating: 1},
		UpdatedAt:     time.Now(),
	}, nil
}

// aggregateBookRatings 聚合书籍评分（从评论中聚合）
func (s *RatingServiceImplementation) aggregateBookRatings(ctx context.Context, bookID string) (*social.RatingStats, error) {
	// 获取书籍的评分统计
	statsMap, err := s.commentRepo.GetBookRatingStats(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取书籍评分统计失败: %w", err)
	}

	// 解析评分统计
	averageRating, _ := statsMap["average_rating"].(float64)
	totalRatings, _ := statsMap["total_ratings"].(int64)

	// 解析评分分布
	distribution := make(map[int]int64)
	if distMap, ok := statsMap["distribution"].(map[string]int64); ok {
		for ratingStr, count := range distMap {
			var rating int
			fmt.Sscanf(ratingStr, "%d", &rating)
			distribution[rating] = count
		}
	}

	return &social.RatingStats{
		TargetID:      bookID,
		TargetType:    "book",
		AverageRating: averageRating,
		TotalRatings:  totalRatings,
		Distribution:  distribution,
		UpdatedAt:     time.Now(),
	}, nil
}

// GetUserRating 获取用户对目标的评分
func (s *RatingServiceImplementation) GetUserRating(ctx context.Context, userID, targetType, targetID string) (int, error) {
	switch targetType {
	case "book":
		// 获取用户对书籍的评论评分
		comments, _, err := s.commentRepo.GetCommentsByBookID(ctx, targetID, 1, 100)
		if err != nil {
			return 0, fmt.Errorf("获取评论失败: %w", err)
		}

		// 查找该用户的评论
		for _, comment := range comments {
			if comment.AuthorID == userID && comment.Rating > 0 {
				return comment.Rating, nil
			}
		}
		return 0, nil // 用户未评分
	case "review":
		// 获取书评并返回评分
		review, err := s.reviewRepo.GetReviewByID(ctx, targetID)
		if err != nil {
			return 0, fmt.Errorf("获取书评失败: %w", err)
		}
		if review.UserID == userID {
			return review.Rating, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的目标类型: %s", targetType)
	}
}

// InvalidateCache 使缓存失效
func (s *RatingServiceImplementation) InvalidateCache(ctx context.Context, targetType, targetID string) error {
	if s.redisClient == nil {
		s.logger.Debug("Redis客户端未初始化，跳过缓存失效")
		return nil
	}

	cacheKey := fmt.Sprintf("rating:stats:%s:%s", targetType, targetID)
	err := s.redisClient.Delete(ctx, cacheKey)
	if err != nil {
		s.logger.Warn("缓存失效失败",
			zap.String("cacheKey", cacheKey),
			zap.Error(err))
		return fmt.Errorf("缓存失效失败: %w", err)
	}

	s.logger.Debug("缓存已失效",
		zap.String("cacheKey", cacheKey))
	return nil
}

// serializeStats 序列化评分统计
func (s *RatingServiceImplementation) serializeStats(stats *social.RatingStats) (string, error) {
	data, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}
	return string(data), nil
}

// deserializeStats 反序列化评分统计
func (s *RatingServiceImplementation) deserializeStats(data string) (*social.RatingStats, error) {
	var stats social.RatingStats
	err := json.Unmarshal([]byte(data), &stats)
	if err != nil {
		return nil, fmt.Errorf("反序列化失败: %w", err)
	}
	return &stats, nil
}

// calculateTTLWithJitter 计算带随机抖动的TTL（防止缓存雪崩）
func (s *RatingServiceImplementation) calculateTTLWithJitter(baseTTL time.Duration) time.Duration {
	// 计算抖动范围：±10%
	jitterRange := float64(baseTTL) * TTLJitterPercent
	// 生成随机抖动值：[-jitterRange, +jitterRange]
	jitter := rand.Float64()*2*jitterRange - jitterRange
	// 应用抖动
	ttl := float64(baseTTL) + jitter
	if ttl < 0 {
		ttl = 0
	}
	return time.Duration(ttl)
}
