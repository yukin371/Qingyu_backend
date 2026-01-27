package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book 书籍实体
// 代表书店领域中的书籍概念,包含书籍的核心业务属性和行为
type Book struct {
	// ID 唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Title 书籍标题
	Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`

	// Author 作者
	Author string `bson:"author" json:"author" validate:"required,min=1,max=100"`

	// ISBN 国际标准书号
	ISBN string `bson:"isbn,omitempty" json:"isbn,omitempty" validate:"omitempty,isbn"`

	// Description 书籍简介
	Description string `bson:"description,omitempty" json:"description,omitempty" validate:"omitempty,max=2000"`

	// CategoryID 分类ID (外键,关联到分类表)
	CategoryID *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`

	// Tags 标签列表 (用于分类和搜索)
	Tags []string `bson:"tags,omitempty" json:"tags,omitempty" validate:"omitempty,max=10"`

	// CoverImage 封面图片URL
	CoverImage string `bson:"cover_image,omitempty" json:"cover_image,omitempty" validate:"omitempty,url"`

	// Price 价格 (单位:分)
	Price int64 `bson:"price" json:"price" validate:"required,min=0"`

	// Stock 库存数量
	Stock int `bson:"stock" json:"stock" validate:"required,min=0"`

	// Status 书籍状态
	Status BookStatus `bson:"status" json:"status" validate:"required"`

	// ViewCount 浏览次数
	ViewCount int64 `bson:"view_count" json:"view_count" validate:"min=0"`

	// LikeCount 点赞数
	LikeCount int64 `bson:"like_count" json:"like_count" validate:"min=0"`

	// CollectCount 收藏数
	CollectCount int64 `bson:"collect_count" json:"collect_count" validate:"min=0"`

	// WordCount 字数 (用于小说类书籍)
	WordCount int64 `bson:"word_count,omitempty" json:"word_count,omitempty" validate:"omitempty,min=0"`

	// PublishDate 出版日期
	PublishDate *time.Time `bson:"publish_date,omitempty" json:"publish_date,omitempty"`

	// Publisher 出版社
	Publisher string `bson:"publisher,omitempty" json:"publisher,omitempty" validate:"omitempty,max=100"`

	// PublishedAt 发布时间
	PublishedAt time.Time `bson:"published_at" json:"published_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// CreatedAt 创建时间
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// BookStatus 书籍状态枚举
type BookStatus string

const (
	// BookStatusDraft 草稿
	BookStatusDraft BookStatus = "draft"

	// BookStatusOngoing 连载中
	BookStatusOngoing BookStatus = "ongoing"

	// BookStatusCompleted 已完结
	BookStatusCompleted BookStatus = "completed"

	// BookStatusPaused 已暂停
	BookStatusPaused BookStatus = "paused"

	// BookStatusDeleted 已删除
	BookStatusDeleted BookStatus = "deleted"
)

// IsValid 验证书籍状态是否有效
func (s BookStatus) IsValid() bool {
	switch s {
	case BookStatusDraft, BookStatusOngoing, BookStatusCompleted, BookStatusPaused, BookStatusDeleted:
		return true
	default:
		return false
	}
}

// AddStock 增加库存
// 业务规则: 库存增加必须为正数
func (b *Book) AddStock(delta int) error {
	if delta <= 0 {
		return &BusinessError{
			Field:   "stock",
			Message: "库存增加量必须为正数",
		}
	}
	b.Stock += delta
	b.UpdatedAt = time.Now()
	return nil
}

// DecreaseStock 减少库存
// 业务规则: 库存减少不能导致库存为负数
func (b *Book) DecreaseStock(delta int) error {
	if delta <= 0 {
		return &BusinessError{
			Field:   "stock",
			Message: "库存减少量必须为正数",
		}
	}
	if b.Stock < delta {
		return &BusinessError{
			Field:   "stock",
			Message: "库存不足",
		}
	}
	b.Stock -= delta
	b.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus 更新书籍状态
// 业务规则: 只能更新到有效的状态
func (b *Book) UpdateStatus(status BookStatus) error {
	if !status.IsValid() {
		return &BusinessError{
			Field:   "status",
			Message: "无效的书籍状态",
		}
	}
	b.Status = status
	b.UpdatedAt = time.Now()
	return nil
}

// IsInStock 检查是否有库存
func (b *Book) IsInStock() bool {
	return b.Stock > 0
}

// IsPublished 检查是否已发布
func (b *Book) IsPublished() bool {
	return b.Status == BookStatusOngoing || b.Status == BookStatusCompleted
}

// IncrementView 增加浏览次数
func (b *Book) IncrementView() {
	b.ViewCount++
	b.UpdatedAt = time.Now()
}

// IncrementLike 增加点赞数
func (b *Book) IncrementLike() {
	b.LikeCount++
	b.UpdatedAt = time.Now()
}

// DecrementLike 减少点赞数
func (b *Book) DecrementLike() {
	if b.LikeCount > 0 {
		b.LikeCount--
		b.UpdatedAt = time.Now()
	}
}

// IncrementCollect 增加收藏数
func (b *Book) IncrementCollect() {
	b.CollectCount++
	b.UpdatedAt = time.Now()
}

// DecrementCollect 减少收藏数
func (b *Book) DecrementCollect() {
	if b.CollectCount > 0 {
		b.CollectCount--
		b.UpdatedAt = time.Now()
	}
}

// BusinessError 业务错误
type BusinessError struct {
	Field   string
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}
