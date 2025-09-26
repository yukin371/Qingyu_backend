package bookstore

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book 书籍模型
type Book struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title        string              `bson:"title" json:"title" validate:"required,min=1,max=100"`               // 书名
	Author       string              `bson:"author" json:"author" validate:"required,min=1,max=50"`             // 作者
	Introduction string              `bson:"introduction" json:"introduction" validate:"max=1000"`               // 简介
	Cover        string              `bson:"cover" json:"cover" validate:"url"`                                 // 封面URL
	CategoryIDs  []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`                                  // 分类ID列表
	Tags         []string            `bson:"tags" json:"tags"`                                                  // 标签
	Status       string              `bson:"status" json:"status" validate:"required,oneof=draft published"`   // 状态：draft草稿 published已发布
	WordCount    int64               `bson:"word_count" json:"wordCount"`                                       // 字数
	ChapterCount int                 `bson:"chapter_count" json:"chapterCount"`                                 // 章节数
	ViewCount    int64               `bson:"view_count" json:"viewCount"`                                       // 浏览量
	LikeCount    int64               `bson:"like_count" json:"likeCount"`                                       // 点赞数
	CommentCount int64               `bson:"comment_count" json:"commentCount"`                                 // 评论数
	Rating       float64             `bson:"rating" json:"rating"`                                              // 评分 0-5
	RatingCount  int64               `bson:"rating_count" json:"ratingCount"`                                   // 评分人数
	IsRecommended bool               `bson:"is_recommended" json:"isRecommended"`                               // 是否推荐
	IsFeatured   bool                `bson:"is_featured" json:"isFeatured"`                                     // 是否精选
	PublishedAt  *time.Time          `bson:"published_at,omitempty" json:"publishedAt,omitempty"`              // 发布时间
	CreatedAt    time.Time           `bson:"created_at" json:"createdAt"`                                       // 创建时间
	UpdatedAt    time.Time           `bson:"updated_at" json:"updatedAt"`                                       // 更新时间
}

// BookFilter 书籍查询过滤器
type BookFilter struct {
	CategoryID    *primitive.ObjectID `json:"categoryId,omitempty"`
	Author        *string            `json:"author,omitempty"`
	Status        *string            `json:"status,omitempty"`
	IsRecommended *bool              `json:"isRecommended,omitempty"`
	IsFeatured    *bool              `json:"isFeatured,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
	MinRating     *float64           `json:"minRating,omitempty"`
	Keyword       *string            `json:"keyword,omitempty"`
	SortBy        string             `json:"sortBy,omitempty"` // created_at, updated_at, view_count, like_count, rating
	SortOrder     string             `json:"sortOrder,omitempty"` // asc, desc
	Limit         int                `json:"limit,omitempty"`
	Offset        int                `json:"offset,omitempty"`
}

// BookStats 书籍统计信息
type BookStats struct {
	TotalBooks      int64 `json:"totalBooks"`
	PublishedBooks  int64 `json:"publishedBooks"`
	DraftBooks      int64 `json:"draftBooks"`
	RecommendedBooks int64 `json:"recommendedBooks"`
	FeaturedBooks   int64 `json:"featuredBooks"`
}
