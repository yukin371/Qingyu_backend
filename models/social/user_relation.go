package social

import (
	"time"
)

// RelationStatus 关系状态
type RelationStatus string

const (
	RelationStatusActive   RelationStatus = "active"   // 有效关注
	RelationStatusInactive RelationStatus = "inactive" // 已取消关注
)

// UserRelation 用户关系模型
type UserRelation struct {
	IdentifiedEntity `bson:",inline"`
	Timestamps       `bson:",inline"`

	// 关注关系
	FollowerID string `bson:"follower_id" json:"followerId" validate:"required"` // 关注者ID
	FolloweeID string `bson:"followee_id" json:"followeeId" validate:"required"` // 被关注者ID

	// 状态
	Status RelationStatus `bson:"status" json:"status" validate:"required,oneof=active inactive"` // 关系状态
}

// TableName 集合名称
func (UserRelation) TableName() string {
	return "user_relations"
}

// IsActive 判断关系是否有效
func (r *UserRelation) IsActive() bool {
	return r.Status == RelationStatusActive
}

// Activate 激活关系
func (r *UserRelation) Activate() {
	r.Status = RelationStatusActive
	r.UpdatedAt = time.Now()
}

// Deactivate 停用关系
func (r *UserRelation) Deactivate() {
	r.Status = RelationStatusInactive
	r.UpdatedAt = time.Now()
}
