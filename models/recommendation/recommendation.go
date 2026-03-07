package recommendation

import (
	"time"

	sharedtypes "Qingyu_backend/models/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserBehaviorRecord 用户行为记录（旧版，保留用于兼容）
// 注意：新的推荐系统应使用 Behavior 模型
type UserBehaviorRecord struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     string                 `json:"user_id" bson:"user_id"`
	ItemID     string                 `json:"item_id" bson:"item_id"`
	ItemType   string                 `json:"item_type" bson:"item_type"`                   // book, chapter, article
	ActionType string                 `json:"action_type" bson:"action_type"`               // view, click, favorite, read, purchase
	Duration   int64                  `json:"duration" bson:"duration"`                     // 停留时长（秒）
	Progress   float64                `json:"progress,omitempty" bson:"progress,omitempty"` // 阅读进度（0-1）
	Score      float64                `json:"score,omitempty" bson:"score,omitempty"`       // 评分（1-5）
	Metadata   map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"` // 额外数据
	IP         string                 `json:"ip,omitempty" bson:"ip,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
}

// UserBehavior 是 UserBehaviorRecord 的别名，用于向后兼容
// Deprecated: 使用 UserBehaviorRecord 或 Behavior 代替
type UserBehavior = UserBehaviorRecord

// RecommendedItem 推荐项（主要存储在缓存中）
type RecommendedItem struct {
	ItemID      string                 `json:"item_id"`
	ItemType    string                 `json:"item_type"`
	Score       float64                `json:"score"`                  // 推荐分数
	Reason      string                 `json:"reason"`                 // 推荐理由
	Rank        int                    `json:"rank"`                   // 排名
	Algorithm   string                 `json:"algorithm,omitempty"`    // 推荐算法
	Metadata    map[string]interface{} `json:"metadata,omitempty"`     // 额外信息
	GeneratedAt time.Time              `json:"generated_at,omitempty"` // 生成时间
}

// 行为类型（保留旧常量名用于兼容，规范值统一走 shared/types）
const (
	ActionTypeView     = string(sharedtypes.RecommendationBehaviorView)    // 浏览
	ActionTypeClick    = string(sharedtypes.RecommendationBehaviorClick)   // 点击
	ActionTypeFavorite = "favorite"                                        // 旧别名，规范值为 collect
	ActionTypeCollect  = string(sharedtypes.RecommendationBehaviorCollect) // 收藏
	ActionTypeRead     = string(sharedtypes.RecommendationBehaviorRead)    // 阅读
	ActionTypeFinish   = string(sharedtypes.RecommendationBehaviorFinish)  // 完读
	ActionTypePurchase = string(sharedtypes.RecommendationBehaviorPurchase)
	ActionTypeComment  = string(sharedtypes.RecommendationBehaviorComment)
	ActionTypeShare    = string(sharedtypes.RecommendationBehaviorShare)
	ActionTypeRate     = string(sharedtypes.RecommendationBehaviorRate)
)

// 内容类型
const (
	ItemTypeBook    = "book"    // 书籍
	ItemTypeChapter = "chapter" // 章节
	ItemTypeArticle = "article" // 文章
)
