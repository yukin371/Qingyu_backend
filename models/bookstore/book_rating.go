package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRating 书籍评分模型
type BookRating struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID    primitive.ObjectID `bson:"book_id" json:"book_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Rating    int                `bson:"rating" json:"rating"` // 1-5星
	Comment   string             `bson:"comment" json:"comment"`
	Tags      []string           `bson:"tags" json:"tags"`
	Likes     int                `bson:"likes" json:"likes"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate 在创建前设置时间戳
func (br *BookRating) BeforeCreate() {
	now := time.Now()
	br.CreatedAt = now
	br.UpdatedAt = now
}

// BeforeUpdate 在更新前刷新更新时间戳
func (br *BookRating) BeforeUpdate() {
	br.UpdatedAt = time.Now()
}

// IsValidRating 验证评分是否有效
func (br *BookRating) IsValidRating() bool {
	return br.Rating >= 1 && br.Rating <= 5
}

// HasComment 检查是否有评论
func (br *BookRating) HasComment() bool {
	return len(br.Comment) > 0
}

// HasTags 检查是否有标签
func (br *BookRating) HasTags() bool {
	return len(br.Tags) > 0
}

// IncrementLikes 增加点赞数
func (br *BookRating) IncrementLikes() {
	br.Likes++
	br.BeforeUpdate()
}

// DecrementLikes 减少点赞数
func (br *BookRating) DecrementLikes() {
	if br.Likes > 0 {
		br.Likes--
		br.BeforeUpdate()
	}
}
