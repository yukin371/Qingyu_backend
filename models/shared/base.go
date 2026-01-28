package shared

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseEntity 通用实体基类
// 提供所有实体共有的基础字段，减少重复代码
type BaseEntity struct {
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// Touch 更新时间戳
func (b *BaseEntity) Touch(t ...time.Time) {
	if len(t) > 1 {
		// 记录错误但不处理，使用默认值
		_ = fmt.Errorf("Touch 方法只允许接受最多一个时间参数")
	}
	if len(t) > 0 {
		b.UpdatedAt = t[0]
	} else {
		b.UpdatedAt = time.Now()
	}
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
	b.Touch(now)
}

// IsDeleted 判断是否已删除
func (b *BaseEntity) IsDeleted() bool {
	if b.DeletedAt == nil {
		return false
	}
	return !b.DeletedAt.IsZero()
}

// IdentifiedEntity 包含ID字段的基础实体
type IdentifiedEntity struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

// GetID 获取ID
func (i *IdentifiedEntity) GetID() primitive.ObjectID {
	return i.ID
}

// SetID 设置ID
func (i *IdentifiedEntity) SetID(id primitive.ObjectID) {
	i.ID = id
}

// GenerateID 生成新的ObjectID
func (i *IdentifiedEntity) GenerateID() {
	if i.ID.IsZero() {
		i.ID = primitive.NewObjectID()
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

// Edited 编辑追踪混入
// 用于需要追踪最后编辑信息的实体（如文档内容、草稿等）
type Edited struct {
	LastSavedAt  time.Time `json:"lastSavedAt" bson:"last_saved_at"`
	LastEditedBy string    `json:"lastEditedBy" bson:"last_edited_by"`
}

// MarkEdited 标记为已编辑
func (e *Edited) MarkEdited(editorID string) {
	e.LastSavedAt = time.Now()
	e.LastEditedBy = editorID
}

// GetLastSavedAt 获取最后保存时间
func (e *Edited) GetLastSavedAt() time.Time {
	return e.LastSavedAt
}

// GetLastEditedBy 获取最后编辑人
func (e *Edited) GetLastEditedBy() string {
	return e.LastEditedBy
}
