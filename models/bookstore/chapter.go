package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chapter 章节模型（元数据，不含内容）
type Chapter struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID     primitive.ObjectID `bson:"book_id" json:"book_id"`
	Title      string             `bson:"title" json:"title"`
	ChapterNum int                `bson:"chapter_num" json:"chapter_num"`
	// Content 字段已移除，使用 ChapterContent 单独存储
	WordCount int     `bson:"word_count" json:"word_count"`
	IsFree    bool    `bson:"is_free" json:"is_free"`
	Price     int64   `bson:"price" json:"price"` // 价格 (分)

	// 内容引用信息
	ContentURL     string `bson:"content_url,omitempty" json:"contentUrl,omitempty"`         // 内容获取 URL
	ContentSize    int64  `bson:"content_size,omitempty" json:"contentSize,omitempty"`       // 内容大小（字节）
	ContentHash    string `bson:"content_hash,omitempty" json:"contentHash,omitempty"`       // 内容哈希（校验用）
	ContentVersion int    `bson:"content_version,omitempty" json:"contentVersion,omitempty"` // 内容版本

	PublishTime time.Time `bson:"publish_time" json:"publish_time"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate 在创建前设置时间戳
func (c *Chapter) BeforeCreate() {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.PublishTime.IsZero() {
		c.PublishTime = now
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (c *Chapter) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

// IsPublished 检查章节是否已发布
func (c *Chapter) IsPublished() bool {
	return !c.PublishTime.After(time.Now())
}

// GetReadingPrice 获取阅读价格 (分)
func (c *Chapter) GetReadingPrice() int64 {
	if c.IsFree {
		return 0
	}
	return c.Price
}

// HasContent 检查是否有关联的内容
func (c *Chapter) HasContent() bool {
	return c.ContentURL != "" || c.ContentHash != ""
}

// GetContentSizeKB 获取内容大小（KB）
func (c *Chapter) GetContentSizeKB() float64 {
	if c.ContentSize == 0 {
		return 0
	}
	return float64(c.ContentSize) / 1024
}

// UpdateContentInfo 更新内容信息
func (c *Chapter) UpdateContentInfo(contentURL string, contentSize int64, hash string, version int) {
	c.ContentURL = contentURL
	c.ContentSize = contentSize
	c.ContentHash = hash
	c.ContentVersion = version
	c.UpdatedAt = time.Now()
}
