package shared

import "time"

// BaseEntity 通用实体基类
// 提供所有实体共有的基础字段，减少重复代码
type BaseEntity struct {
	CreatedAt time.Time  `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
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

// SoftDelete 软删除
func (b *BaseEntity) SoftDelete() {
	now := time.Now()
	b.DeletedAt = &now
	b.Touch()
}

// IsDeleted 判断是否已删除
func (b *BaseEntity) IsDeleted() bool {
	return b.DeletedAt != nil
}

// IdentifiedEntity 包含ID字段的基础实体
type IdentifiedEntity struct {
	ID string `bson:"_id,omitempty" json:"id"`
}

// GetID 获取ID
func (i *IdentifiedEntity) GetID() string {
	return i.ID
}

// SetID 设置ID
func (i *IdentifiedEntity) SetID(id string) {
	i.ID = id
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
