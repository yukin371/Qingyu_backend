package social

import "time"

// RatingStats 评分统计（内存结构，不持久化）
type RatingStats struct {
    TargetID      string         `json:"targetId" bson:"target_id"`
    TargetType    string         `json:"targetType" bson:"target_type"`
    AverageRating float64        `json:"averageRating" bson:"average_rating"`
    TotalRatings  int64          `json:"totalRatings" bson:"total_ratings"`
    Distribution  map[int]int64  `json:"distribution" bson:"distribution"` // {1: count, 2: count, ...}
    UpdatedAt     time.Time      `json:"updatedAt" bson:"updated_at"`
}
