package recommendation

import (
	"context"
	"time"
)

// RecommendationService 推荐服务接口（对外暴露）
type RecommendationService interface {
	// 获取推荐
	GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]*RecommendedItem, error)
	GetSimilarItems(ctx context.Context, itemID string, limit int) ([]*RecommendedItem, error)
	GetHotItems(ctx context.Context, itemType string, limit int) ([]*RecommendedItem, error)

	// 行为记录
	RecordUserBehavior(ctx context.Context, req *RecordBehaviorRequest) error
	GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*UserBehavior, error)

	// 刷新推荐（管理后台使用）
	RefreshRecommendations(ctx context.Context, userID string) error
	RefreshHotItems(ctx context.Context, itemType string) error

	// 健康检查
	Health(ctx context.Context) error
}

// ============ 请求结构 ============

// RecordBehaviorRequest 记录用户行为请求
type RecordBehaviorRequest struct {
	UserID     string                 `json:"user_id" binding:"required"`
	ItemID     string                 `json:"item_id" binding:"required"`
	ItemType   string                 `json:"item_type" binding:"required"`   // book, article, etc.
	ActionType string                 `json:"action_type" binding:"required"` // view, click, favorite, read
	Duration   int64                  `json:"duration,omitempty"`             // 阅读时长（秒）
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 数据结构 ============

// RecommendedItem 推荐项
type RecommendedItem struct {
	ItemID   string  `json:"item_id"`
	ItemType string  `json:"item_type"` // book, article, etc.
	Score    float64 `json:"score"`     // 推荐分数
	Reason   string  `json:"reason"`    // 推荐理由
	Rank     int     `json:"rank"`      // 排名
}

// UserBehavior 用户行为
type UserBehavior struct {
	ID         string                 `json:"id" bson:"_id,omitempty"`
	UserID     string                 `json:"user_id" bson:"user_id"`
	ItemID     string                 `json:"item_id" bson:"item_id"`
	ItemType   string                 `json:"item_type" bson:"item_type"`
	ActionType string                 `json:"action_type" bson:"action_type"`
	Duration   int64                  `json:"duration" bson:"duration"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
}
