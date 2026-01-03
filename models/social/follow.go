package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Follow 用户关注关系
type Follow struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID   string             `bson:"follower_id" json:"follower_id"`     // 关注者ID
	FollowingID  string             `bson:"following_id" json:"following_id"`   // 被关注者ID
	FollowType   string             `bson:"follow_type" json:"follow_type"`     // 关注类型: user, author
	IsMutual     bool               `bson:"is_mutual" json:"is_mutual"`         // 是否互相关注
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// AuthorFollow 作者关注（扩展信息）
type AuthorFollow struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	AuthorID     string             `bson:"author_id" json:"author_id"`
	AuthorName   string             `bson:"author_name" json:"author_name"`
	AuthorAvatar string             `bson:"author_avatar,omitempty" json:"author_avatar,omitempty"`
	NotifyNewBook bool              `bson:"notify_new_book" json:"notify_new_book"` // 新书通知
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// FollowStats 关注统计
type FollowStats struct {
	UserID        string    `bson:"user_id" json:"user_id"`
	FollowerCount int       `bson:"follower_count" json:"follower_count"` // 粉丝数
	FollowingCount int      `bson:"following_count" json:"following_count"` // 关注数
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

// FollowInfo 关注信息
type FollowInfo struct {
	FollowerID  string    `json:"follower_id"`
	FollowerName string   `json:"follower_name"`
	FollowerAvatar string `json:"follower_avatar,omitempty"`
	IsMutual    bool      `json:"is_mutual"`
	CreatedAt   time.Time `json:"created_at"`
}

// FollowingInfo 关注信息
type FollowingInfo struct {
	FollowingID  string    `json:"following_id"`
	FollowingName string   `json:"following_name"`
	FollowingAvatar string `json:"following_avatar,omitempty"`
	IsMutual      bool     `json:"is_mutual"`
	CreatedAt     time.Time `json:"created_at"`
}
