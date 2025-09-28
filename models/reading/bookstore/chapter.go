package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chapter 章节模型
type Chapter struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID      primitive.ObjectID `bson:"book_id" json:"book_id"`
	Title       string             `bson:"title" json:"title"`
	ChapterNum  int                `bson:"chapter_num" json:"chapter_num"`
	Content     string             `bson:"content" json:"content"`
	WordCount   int                `bson:"word_count" json:"word_count"`
	IsFree      bool               `bson:"is_free" json:"is_free"`
	Price       float64            `bson:"price" json:"price"`
	PublishTime time.Time          `bson:"publish_time" json:"publish_time"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
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

// GetReadingPrice 获取阅读价格
func (c *Chapter) GetReadingPrice() float64 {
	if c.IsFree {
		return 0
	}
	return c.Price
}