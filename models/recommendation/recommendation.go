package recommendation

import "time"

// UserBehavior 用户行为记录
type UserBehavior struct {
	ID         string                 `json:"id" bson:"_id,omitempty"`
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

// UserProfile 用户画像（可选，用于推荐算法）
type UserProfile struct {
	UserID          string         `json:"user_id" bson:"user_id"`
	FavoriteGenres  []string       `json:"favorite_genres,omitempty" bson:"favorite_genres,omitempty"` // 偏好类型
	FavoriteAuthors []string       `json:"favorite_authors,omitempty" bson:"favorite_authors,omitempty"`
	ReadingHabits   map[string]int `json:"reading_habits,omitempty" bson:"reading_habits,omitempty"`       // 阅读习惯统计
	AverageReadTime int64          `json:"average_read_time,omitempty" bson:"average_read_time,omitempty"` // 平均阅读时长
	PreferredLength string         `json:"preferred_length,omitempty" bson:"preferred_length,omitempty"`   // 偏好长度：short, medium, long
	Tags            []string       `json:"tags,omitempty" bson:"tags,omitempty"`                           // 用户标签
	LastUpdated     time.Time      `json:"last_updated,omitempty" bson:"last_updated,omitempty"`
}

// 行为类型
const (
	ActionTypeView     = "view"     // 浏览
	ActionTypeClick    = "click"    // 点击
	ActionTypeFavorite = "favorite" // 收藏
	ActionTypeRead     = "read"     // 阅读
	ActionTypePurchase = "purchase" // 购买
	ActionTypeComment  = "comment"  // 评论
	ActionTypeShare    = "share"    // 分享
	ActionTypeRate     = "rate"     // 评分
)

// 内容类型
const (
	ItemTypeBook    = "book"    // 书籍
	ItemTypeChapter = "chapter" // 章节
	ItemTypeArticle = "article" // 文章
)
