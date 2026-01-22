package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review 书评
type Review struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID       string             `bson:"book_id" json:"bookId"`
	UserID       string             `bson:"user_id" json:"userId"`
	UserName     string             `bson:"user_name" json:"userName"`
	UserAvatar   string             `bson:"user_avatar,omitempty" json:"userAvatar,omitempty"`
	Title        string             `bson:"title" json:"title"`
	Content      string             `bson:"content" json:"content"`
	Rating       int                `bson:"rating" json:"rating"` // 1-5星评分
	LikeCount    int                `bson:"like_count" json:"likeCount"`
	CommentCount int                `bson:"comment_count" json:"commentCount"`
	IsSpoiler    bool               `bson:"is_spoiler" json:"isSpoiler"` // 是否包含剧透
	IsPublic     bool               `bson:"is_public" json:"isPublic"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ReviewLike 书评点赞
type ReviewLike struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReviewID  string             `bson:"review_id" json:"reviewId"`
	UserID    string             `bson:"user_id" json:"userId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// ReviewInfo 书评信息
type ReviewInfo struct {
	ID           string    `json:"id"`
	BookID       string    `json:"bookId"`
	BookTitle    string    `json:"book_title,omitempty"`
	BookCover    string    `json:"book_cover,omitempty"`
	UserID       string    `json:"userId"`
	UserName     string    `json:"userName"`
	UserAvatar   string    `json:"userAvatar,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Rating       int       `json:"rating"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
	IsLiked      bool      `json:"isLiked"`
	IsSpoiler    bool      `json:"isSpoiler"`
	CreatedAt    time.Time `json:"createdAt"`
}
