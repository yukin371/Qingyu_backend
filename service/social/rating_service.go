package social

import (
    "context"

    "Qingyu_backend/models/social"
)

// RatingService 评分服务接口
type RatingService interface {
    // GetRatingStats 获取评分统计（带缓存）
    GetRatingStats(ctx context.Context, targetType, targetID string) (*social.RatingStats, error)

    // GetUserRating 获取用户对目标的评分
    GetUserRating(ctx context.Context, userID, targetType, targetID string) (int, error)

    // AggregateRatings 从现有数据源聚合评分
    AggregateRatings(ctx context.Context, targetType, targetID string) (*social.RatingStats, error)

    // InvalidateCache 使缓存失效
    InvalidateCache(ctx context.Context, targetType, targetID string) error
}
