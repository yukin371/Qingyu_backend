package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review 书评
type Review struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID       string             `bson:"book_id" json:"book_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	UserName     string             `bson:"user_name" json:"user_name"`
	UserAvatar   string             `bson:"user_avatar,omitempty" json:"user_avatar,omitempty"`
	Title        string             `bson:"title" json:"title"`
	Content      string             `bson:"content" json:"content"`
	Rating       int                `bson:"rating" json:"rating"` // 1-5星评分
	LikeCount    int                `bson:"like_count" json:"like_count"`
	CommentCount int                `bson:"comment_count" json:"comment_count"`
	IsSpoiler    bool               `bson:"is_spoiler" json:"is_spoiler"` // 是否包含剧透
	IsPublic     bool               `bson:"is_public" json:"is_public"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// ReviewLike 书评点赞
type ReviewLike struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReviewID  string             `bson:"review_id" json:"review_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// ReviewInfo 书评信息
type ReviewInfo struct {
	ID           string    `json:"id"`
	BookID       string    `json:"book_id"`
	BookTitle    string    `json:"book_title,omitempty"`
	BookCover    string    `json:"book_cover,omitempty"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	UserAvatar   string    `json:"user_avatar,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Rating       int       `json:"rating"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	IsLiked      bool      `json:"is_liked"`
	IsSpoiler    bool      `json:"is_spoiler"`
	CreatedAt    time.Time `json:"created_at"`
}
