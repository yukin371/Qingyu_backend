package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
)

// BookStatus 书籍状态枚举
type BookStatus string

const (
	BookStatusDraft     BookStatus = "draft"     // 草稿
	BookStatusOngoing   BookStatus = "ongoing"   // 连载中 (已发布且正在更新)
	BookStatusCompleted BookStatus = "completed" // 已完结
	BookStatusPaused    BookStatus = "paused"    // 暂停更新
)

// Book 书籍模型 - 用于书城列表展示的简略信息
type Book struct {
	shared.IdentifiedEntity `bson:",inline"`
	shared.BaseEntity       `bson:",inline"`

	Title         string               `bson:"title" json:"title" validate:"required,min=1,max=200"`   // 书名
	Author        string               `bson:"author" json:"author" validate:"required,min=1,max=100"` // 作者
	AuthorID      primitive.ObjectID   `bson:"author_id,omitempty" json:"authorId,omitempty"`        // 作者ID
	Introduction  string               `bson:"introduction" json:"introduction" validate:"max=1000"` // 简介
	Cover         string               `bson:"cover" json:"cover" validate:"url"`                    // 封面URL
	CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`                     // 分类ID列表
	Categories    []string             `bson:"categories" json:"categories"`                        // 分类名称（冗余字段，便于展示）
	Tags          []string             `bson:"tags" json:"tags"`                                    // 标签
	Status        BookStatus           `bson:"status" json:"status" validate:"required"`            // 状态
	Rating        types.Rating         `bson:"rating" json:"rating" validate:"min=0,max=5"`         // 评分 (0-5星，平均分可为0)
	RatingCount   int64                `bson:"rating_count" json:"ratingCount" validate:"min=0"`    // 评分人数
	ViewCount     int64                `bson:"view_count" json:"viewCount" validate:"min=0"`        // 浏览量
	WordCount     int64                `bson:"word_count" json:"wordCount" validate:"min=0"`        // 字数
	ChapterCount  int                  `bson:"chapter_count" json:"chapterCount"`                   // 章节数
	Price         float64              `bson:"price" json:"price" validate:"min=0"`                  // 价格 (分，使用float64以兼容MongoDB默认数字类型)
	IsFree        bool                 `bson:"is_free" json:"isFree"`                               // 是否免费
	IsRecommended bool                 `bson:"is_recommended" json:"isRecommended"`                 // 是否推荐
	IsFeatured    bool                 `bson:"is_featured" json:"isFeatured"`                       // 是否精选
	IsHot         bool                 `bson:"is_hot" json:"isHot"`                                 // 是否热门
	PublishedAt   *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"` // 发布时间
	LastUpdateAt  *time.Time           `bson:"last_update_at,omitempty" json:"lastUpdateAt,omitempty"` // 最后更新时间
}

// BookFilter 书籍查询过滤器
type BookFilter struct {
	CategoryID    *string     `json:"categoryId,omitempty"`
	Author        *string     `json:"author,omitempty"`
	AuthorID      *string     `json:"authorId,omitempty"`
	Status        *BookStatus `json:"status,omitempty"`
	IsRecommended *bool       `json:"isRecommended,omitempty"`
	IsFeatured    *bool       `json:"isFeatured,omitempty"`
	IsHot         *bool       `json:"isHot,omitempty"`
	IsFree        *bool       `json:"isFree,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	MinPrice      *float64    `json:"minPrice,omitempty"`
	MaxPrice      *float64    `json:"maxPrice,omitempty"`
	Keyword       *string     `json:"keyword,omitempty"`
	SortBy        string      `json:"sortBy,omitempty"`    // created_at, updated_at, published_at, word_count, chapter_count
	SortOrder     string      `json:"sortOrder,omitempty"` // asc, desc
	Limit         int         `json:"limit,omitempty"`
	Offset        int         `json:"offset,omitempty"`
}

// BookStats 书籍统计信息
type BookStats struct {
	TotalBooks       int64 `json:"totalBooks"`
	PublishedBooks   int64 `json:"publishedBooks"`
	DraftBooks       int64 `json:"draftBooks"`
	RecommendedBooks int64 `json:"recommendedBooks"`
	FeaturedBooks    int64 `json:"featuredBooks"`
}

// IsValid 检查状态是否有效
func (s BookStatus) IsValid() bool {
	switch s {
	case BookStatusDraft, BookStatusOngoing,
		BookStatusCompleted, BookStatusPaused:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (s BookStatus) String() string {
	return string(s)
}

// IsPublic 是否公开
func (s BookStatus) IsPublic() bool {
	return s == BookStatusOngoing || s == BookStatusCompleted
}

// CanEdit 是否可编辑
func (s BookStatus) CanEdit() bool {
	return s == BookStatusDraft || s == BookStatusPaused
}
