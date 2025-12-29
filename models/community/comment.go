package community

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment 评论模型
type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID    string             `bson:"book_id" json:"book_id" binding:"required"`
	ChapterID string             `bson:"chapter_id,omitempty" json:"chapter_id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id" binding:"required"`

	// 评论内容
	Content string `bson:"content" json:"content" binding:"required,min=10,max=500"`
	Rating  int    `bson:"rating" json:"rating" binding:"min=0,max=5"` // 0表示没有评分，1-5星评分

	// 回复相关
	ParentID    string `bson:"parent_id,omitempty" json:"parent_id,omitempty"`         // 父评论ID（一级回复）
	RootID      string `bson:"root_id,omitempty" json:"root_id,omitempty"`             // 根评论ID（用于嵌套回复）
	ReplyToUser string `bson:"reply_to_user,omitempty" json:"reply_to_user,omitempty"` // 回复的目标用户ID

	// 统计数据
	LikeCount  int `bson:"like_count" json:"like_count"`   // 点赞数
	ReplyCount int `bson:"reply_count" json:"reply_count"` // 回复数

	// 审核状态
	Status       string `bson:"status" json:"status"`                                   // approved, rejected, pending
	RejectReason string `bson:"reject_reason,omitempty" json:"reject_reason,omitempty"` // 拒绝原因

	// 时间戳
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// CommentStatus 评论状态常量
const (
	CommentStatusPending  = "pending"  // 待审核
	CommentStatusApproved = "approved" // 已通过
	CommentStatusRejected = "rejected" // 已拒绝
)

// CommentSortBy 评论排序方式
const (
	CommentSortByLatest = "latest" // 最新
	CommentSortByHot    = "hot"    // 最热（点赞数）
)

// IsEditable 检查评论是否可编辑（30分钟内）
func (c *Comment) IsEditable() bool {
	return time.Since(c.CreatedAt) <= 30*time.Minute
}

// IsReply 判断是否为回复评论
func (c *Comment) IsReply() bool {
	return c.ParentID != ""
}

// HasRating 判断是否包含评分
func (c *Comment) HasRating() bool {
	return c.Rating > 0
}

// CommentThread 评论线程（树状结构）
type CommentThread struct {
	Comment *Comment         `json:"comment"`  // 主评论
	Replies []*CommentThread `json:"replies"`  // 回复列表（支持嵌套）
	Total   int64            `json:"total"`    // 总回复数
	HasMore bool             `json:"has_more"` // 是否还有更多回复
}
