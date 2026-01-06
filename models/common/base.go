package common

import "time"

// BaseEntity 通用实体基类
// 提供所有实体共有的基础字段，减少重复代码
type BaseEntity struct {
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

// Touch 更新 UpdatedAt 时间戳
func (b *BaseEntity) Touch() {
	b.UpdatedAt = time.Now()
}

// TouchForCreate 创建时设置时间戳
func (b *BaseEntity) TouchForCreate() {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}
}

// ReadStatus 已读状态混入
// 用于需要追踪阅读状态的消息、通知等实体
type ReadStatus struct {
	IsRead bool       `json:"isRead" bson:"is_read"`
	ReadAt *time.Time `json:"readAt,omitempty" bson:"read_at,omitempty"`
}

// MarkAsRead 标记为已读
func (r *ReadStatus) MarkAsRead() {
	r.IsRead = true
	now := time.Now()
	r.ReadAt = &now
}

// MarkAsUnread 标记为未读
func (r *ReadStatus) MarkAsUnread() {
	r.IsRead = false
	r.ReadAt = nil
}

// IsRecentlyRead 检查是否在最近N分钟内已读
func (r *ReadStatus) IsRecentlyRead(minutes int) bool {
	if !r.IsRead || r.ReadAt == nil {
		return false
	}
	duration := time.Since(*r.ReadAt)
	return duration <= time.Duration(minutes)*time.Minute
}
