package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Follow 用户关注关系
type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID  string             `bson:"follower_id" json:"followerId"`   // 关注者ID
	FollowingID string             `bson:"following_id" json:"followingId"` // 被关注者ID
	FollowType  string             `bson:"follow_type" json:"followType"`   // 关注类型: user, author
	IsMutual    bool               `bson:"is_mutual" json:"isMutual"`       // 是否互相关注
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// AuthorFollow 作者关注（扩展信息）
type AuthorFollow struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"user_id"`
	AuthorID      string             `bson:"author_id" json:"authorId"`
	AuthorName    string             `bson:"author_name" json:"authorName"`
	AuthorAvatar  string             `bson:"author_avatar,omitempty" json:"author_avatar,omitempty"`
	NotifyNewBook bool               `bson:"notify_new_book" json:"notifyNewBook"` // 新书通知
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// FollowStats 关注统计
type FollowStats struct {
	UserID         string    `bson:"user_id" json:"user_id"`
	FollowerCount  int       `bson:"follower_count" json:"followerCount"`   // 粉丝数
	FollowingCount int       `bson:"following_count" json:"followingCount"` // 关注数
	UpdatedAt      time.Time `bson:"updated_at" json:"updatedAt"`
}

// FollowInfo 关注信息
type FollowInfo struct {
	FollowerID     string    `json:"followerId"`
	FollowerName   string    `json:"followerName"`
	FollowerAvatar string    `json:"follower_avatar,omitempty"`
	IsMutual       bool      `json:"isMutual"`
	CreatedAt      time.Time `json:"createdAt"`
}

// FollowingInfo 关注信息
type FollowingInfo struct {
	FollowingID     string    `json:"followingId"`
	FollowingName   string    `json:"followingName"`
	FollowingAvatar string    `json:"following_avatar,omitempty"`
	IsMutual        bool      `json:"isMutual"`
	CreatedAt       time.Time `json:"createdAt"`
}
