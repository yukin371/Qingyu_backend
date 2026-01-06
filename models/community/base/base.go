package base

import "time"

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

// IdentifiedEntity ID实体混入
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

// ThreadedConversation 支持回复的对话混入
type ThreadedConversation struct {
	ParentID   *string `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父评论ID（用于回复）
	RootID     *string `bson:"root_id,omitempty" json:"rootId,omitempty"`     // 根评论ID（用于对话树）
	ThreadSize int     `bson:"thread_size" json:"threadSize"`                  // 回复数量
}

// IsReply 判断是否为回复
func (t *ThreadedConversation) IsReply() bool {
	return t.ParentID != nil
}

// AddReply 增加回复计数
func (t *ThreadedConversation) AddReply() {
	t.ThreadSize++
}

// Likable 可点赞实体混入
type Likable struct {
	LikeCount int64 `bson:"like_count" json:"likeCount"` // 点赞数
}

// IncrementLike 增加点赞
func (l *Likable) IncrementLike() {
	l.LikeCount++
}

// DecrementLike 减少点赞
func (l *Likable) DecrementLike() {
	if l.LikeCount > 0 {
		l.LikeCount--
	}
}
