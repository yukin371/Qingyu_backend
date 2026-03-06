package recommendation

import (
	sharedtypes "Qingyu_backend/models/shared/types"
	"time"
)

// Behavior 用户行为事件
// 用于记录用户在书籍/章节上的交互，用于离线/实时推荐特征计算
type Behavior struct {
	ID           string                 `bson:"_id,omitempty" json:"id"`
	UserID       string                 `bson:"user_id" json:"userId"`
	ItemID       string                 `bson:"item_id" json:"itemId"` // 物品ID（书籍ID）
	ChapterID    string                 `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`
	BehaviorType string                 `bson:"behavior_type" json:"behaviorType"` // view/click/collect/read/finish/like/share
	Value        float64                `bson:"value" json:"value"`                // 行为强度或得分（如阅读时长、权重）
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	OccurredAt   time.Time              `bson:"occurred_at" json:"occurredAt"`
	CreatedAt    time.Time              `bson:"created_at" json:"createdAt"`
}

const (
	BehaviorTypeView     = string(sharedtypes.RecommendationBehaviorView)
	BehaviorTypeClick    = string(sharedtypes.RecommendationBehaviorClick)
	BehaviorTypeCollect  = string(sharedtypes.RecommendationBehaviorCollect)
	BehaviorTypeRead     = string(sharedtypes.RecommendationBehaviorRead)
	BehaviorTypeFinish   = string(sharedtypes.RecommendationBehaviorFinish)
	BehaviorTypeLike     = string(sharedtypes.RecommendationBehaviorLike)
	BehaviorTypeShare    = string(sharedtypes.RecommendationBehaviorShare)
	BehaviorTypePurchase = string(sharedtypes.RecommendationBehaviorPurchase)
	BehaviorTypeComment  = string(sharedtypes.RecommendationBehaviorComment)
	BehaviorTypeRate     = string(sharedtypes.RecommendationBehaviorRate)
)
