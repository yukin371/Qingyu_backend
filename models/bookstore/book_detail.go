package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 使用 book.go 中定义的 BookStatus 枚举

// BookDetail 书籍详情模型 - 用于详情页面展示的完整信息
// 包含统计数据、交互数据等详细信息，适用于书籍详情页面
type BookDetail struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                      // 书籍ID
	Title        string             `bson:"title" json:"title" validate:"required,min=1,max=200"`   // 书名
	Subtitle     string             `bson:"subtitle" json:"subtitle" validate:"max=200"`            // 副标题
	Author       string             `bson:"author" json:"author" validate:"required,min=1,max=100"` // 作者名
	AuthorID     primitive.ObjectID `bson:"author_id" json:"author_id"`                             // 作者ID
	Description  string             `bson:"description" json:"description" validate:"max=5000"`     // 详细描述
	Introduction string             `bson:"introduction" json:"introduction" validate:"max=1000"`   // 简介
	CoverURL     string             `bson:"cover_url" json:"cover_url" validate:"url"`              // 封面图片URL

	// 网络小说特有字段
	SerializedAt time.Time  `bson:"serialized_at" json:"serialized_at"`                   // 开始连载时间
	CompletedAt  *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"` // 完结时间

	Categories   []string             `bson:"categories" json:"categories"`                        // 分类列表
	CategoryIDs  []primitive.ObjectID `bson:"category_ids" json:"category_ids"`                    // 分类ID列表
	Tags         []string             `bson:"tags" json:"tags"`                                    // 标签
	Status       BookStatus           `bson:"status" json:"status"`                                // 状态
	WordCount    int64                `bson:"word_count" json:"word_count" validate:"min=0"`       // 总字数
	ChapterCount int64                `bson:"chapter_count" json:"chapter_count" validate:"min=0"` // 章节数
	Price        int64                `bson:"price" json:"price" validate:"min=0"`                 // 价格 (分，按章节或全本)
	IsFree       bool                 `bson:"is_free" json:"is_free"`                              // 是否免费

	// 统计数据
	ViewCount    int64   `bson:"view_count" json:"view_count" validate:"min=0"`       // 浏览量
	LikeCount    int64   `bson:"like_count" json:"like_count" validate:"min=0"`       // 点赞数
	CommentCount int64   `bson:"comment_count" json:"comment_count" validate:"min=0"` // 评论数
	ShareCount   int64   `bson:"share_count" json:"share_count" validate:"min=0"`     // 分享数
	CollectCount int64   `bson:"collect_count" json:"collect_count" validate:"min=0"` // 收藏数
	Rating       float64 `bson:"rating" json:"rating" validate:"min=1,max=5"`         // 评分 (1-5星)
	RatingCount  int64   `bson:"rating_count" json:"rating_count" validate:"min=0"`   // 评分人数

	// 最新章节信息
	LastChapterTitle string    `bson:"last_chapter_title" json:"last_chapter_title"` // 最新章节标题
	LastChapterAt    time.Time `bson:"last_chapter_at" json:"last_chapter_at"`       // 最新章节更新时间

	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"` // 更新时间
}

// BeforeCreate 在创建前设置时间戳
func (bd *BookDetail) BeforeCreate() {
	now := time.Now()
	bd.CreatedAt = now
	bd.UpdatedAt = now
}

// BeforeUpdate 在更新前刷新更新时间戳
func (bd *BookDetail) BeforeUpdate() {
	bd.UpdatedAt = time.Now()
}

// IsCompleted 检查书籍是否已完结
func (bd *BookDetail) IsCompleted() bool {
	return bd.Status == BookStatusCompleted
}

// IsOngoing 检查书籍是否连载中
func (bd *BookDetail) IsOngoing() bool {
	return bd.Status == BookStatusOngoing
}

// IsPaused 检查书籍是否暂停更新
func (bd *BookDetail) IsPaused() bool {
	return bd.Status == BookStatusPaused
}
