package base

import "time"

// BaseEntity 基础实体接口
type BaseEntity interface {
	GetID() string
	SetID(string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	Touch()
}

// Timestamps 时间戳混入
type Timestamps struct {
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// Touch 更新时间戳
func (t *Timestamps) Touch() {
	t.UpdatedAt = time.Now()
}

// TouchForCreate 创建时设置时间戳
func (t *Timestamps) TouchForCreate() {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
}

// SoftDelete 软删除
func (t *Timestamps) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
	t.Touch()
}

// Restore 恢复
func (t *Timestamps) Restore() {
	t.DeletedAt = nil
	t.Touch()
}

// IsDeleted 判断是否已删除
func (t *Timestamps) IsDeleted() bool {
	return t.DeletedAt != nil
}

// NamedEntity 命名实体混入
type NamedEntity struct {
	Name string `bson:"name" json:"name" validate:"required,min=1,max=100"`
}

// DescriptedEntity 描述实体混入
type DescriptedEntity struct {
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

// TitledEntity 标题实体混入
type TitledEntity struct {
	Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`
}
