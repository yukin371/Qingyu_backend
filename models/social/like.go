package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like 点赞模型
type Like struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"userId" binding:"required"`
	TargetType string             `bson:"target_type" json:"targetType" binding:"required"` // book, comment, chapter
	TargetID   string             `bson:"target_id" json:"targetId" binding:"required"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
}

// LikeTargetType 点赞目标类型常量
const (
	LikeTargetTypeBook    = "book"    // 书籍
	LikeTargetTypeComment = "comment" // 评论
	LikeTargetTypeChapter = "chapter" // 章节
)

// IsBook 判断是否为书籍点赞
func (l *Like) IsBook() bool {
	return l.TargetType == LikeTargetTypeBook
}

// IsComment 判断是否为评论点赞
func (l *Like) IsComment() bool {
	return l.TargetType == LikeTargetTypeComment
}

// IsChapter 判断是否为章节点赞
func (l *Like) IsChapter() bool {
	return l.TargetType == LikeTargetTypeChapter
}
