package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookStatus 书籍状态枚举
type BookStatus string

const (
	BookStatusDraft     BookStatus = "draft"     // 草稿
	BookStatusPublished BookStatus = "published" // 已发布
	BookStatusOngoing   BookStatus = "ongoing"   // 连载中
	BookStatusCompleted BookStatus = "completed" // 已完结
	BookStatusPaused    BookStatus = "paused"    // 暂停更新
)

// Book 书籍模型 - 用于书城列表展示的简略信息
type Book struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title         string               `bson:"title" json:"title" validate:"required,min=1,max=100"`   // 书名
	Author        string               `bson:"author" json:"author" validate:"required,min=1,max=50"`  // 作者
	AuthorID      primitive.ObjectID   `bson:"author_id,omitempty" json:"authorId,omitempty"`          // 作者ID
	Introduction  string               `bson:"introduction" json:"introduction" validate:"max=500"`    // 简介（列表展示用，较短）
	Cover         string               `bson:"cover" json:"cover" validate:"url"`                      // 封面URL
	CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`                        // 分类ID列表
	Categories    []string             `bson:"categories" json:"categories"`                           // 分类名称（冗余字段，便于展示）
	Tags          []string             `bson:"tags" json:"tags"`                                       // 标签
	Status        BookStatus           `bson:"status" json:"status" validate:"required"`               // 状态
	Rating        float64              `bson:"rating" json:"rating"`                                   // 评分 (0-10)
	RatingCount   int64                `bson:"rating_count" json:"ratingCount"`                        // 评分人数
	ViewCount     int64                `bson:"view_count" json:"viewCount"`                            // 浏览量
	WordCount     int64                `bson:"word_count" json:"wordCount"`                            // 字数
	ChapterCount  int                  `bson:"chapter_count" json:"chapterCount"`                      // 章节数
	Price         float64              `bson:"price" json:"price"`                                     // 价格
	IsFree        bool                 `bson:"is_free" json:"isFree"`                                  // 是否免费
	IsRecommended bool                 `bson:"is_recommended" json:"isRecommended"`                    // 是否推荐
	IsFeatured    bool                 `bson:"is_featured" json:"isFeatured"`                          // 是否精选
	IsHot         bool                 `bson:"is_hot" json:"isHot"`                                    // 是否热门
	PublishedAt   *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"`    // 发布时间
	LastUpdateAt  *time.Time           `bson:"last_update_at,omitempty" json:"lastUpdateAt,omitempty"` // 最后更新时间
	CreatedAt     time.Time            `bson:"created_at" json:"createdAt"`                            // 创建时间
	UpdatedAt     time.Time            `bson:"updated_at" json:"updatedAt"`                            // 更新时间
}

// BookFilter 书籍查询过滤器
type BookFilter struct {
	CategoryID    *primitive.ObjectID `json:"categoryId,omitempty"`
	Author        *string             `json:"author,omitempty"`
	AuthorID      *primitive.ObjectID `json:"authorId,omitempty"`
	Status        *BookStatus         `json:"status,omitempty"`
	IsRecommended *bool               `json:"isRecommended,omitempty"`
	IsFeatured    *bool               `json:"isFeatured,omitempty"`
	IsHot         *bool               `json:"isHot,omitempty"`
	IsFree        *bool               `json:"isFree,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
	MinPrice      *float64            `json:"minPrice,omitempty"`
	MaxPrice      *float64            `json:"maxPrice,omitempty"`
	Keyword       *string             `json:"keyword,omitempty"`
	SortBy        string              `json:"sortBy,omitempty"`    // created_at, updated_at, published_at, word_count, chapter_count
	SortOrder     string              `json:"sortOrder,omitempty"` // asc, desc
	Limit         int                 `json:"limit,omitempty"`
	Offset        int                 `json:"offset,omitempty"`
}

// BookStats 书籍统计信息
type BookStats struct {
	TotalBooks       int64 `json:"totalBooks"`
	PublishedBooks   int64 `json:"publishedBooks"`
	DraftBooks       int64 `json:"draftBooks"`
	RecommendedBooks int64 `json:"recommendedBooks"`
	FeaturedBooks    int64 `json:"featuredBooks"`
}
