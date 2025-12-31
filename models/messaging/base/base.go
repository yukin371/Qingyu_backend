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

// SoftDelete 软删除
func (t *Timestamps) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
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

// CommunicationBase 通信基础实体混入
type CommunicationBase struct {
	SenderID   string     `bson:"sender_id" json:"senderId" validate:"required"` // 发送者ID
	ReceiverID string     `bson:"receiver_id" json:"receiverId" validate:"required"` // 接收者ID
	IsRead     bool       `bson:"is_read" json:"isRead"`    // 是否已读
	ReadAt     *time.Time `bson:"read_at,omitempty" json:"readAt,omitempty"` // 已读时间
}

// MarkAsRead 标记为已读
func (c *CommunicationBase) MarkAsRead() {
	if !c.IsRead {
		c.IsRead = true
		now := time.Now()
		c.ReadAt = &now
	}
}

// MarkAsUnread 标记为未读
func (c *CommunicationBase) MarkAsUnread() {
	c.IsRead = false
	c.ReadAt = nil
}

// ThreadedConversation 支持回复的对话混入
type ThreadedConversation struct {
	ParentID   *string `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父消息ID（用于回复）
	RootID     *string `bson:"root_id,omitempty" json:"rootId,omitempty"`     // 根消息ID（用于对话树）
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

// Expirable 有效期混入
type Expirable struct {
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expiresAt,omitempty"` // 过期时间
}

// IsExpired 判断是否已过期
func (e *Expirable) IsExpired() bool {
	if e.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*e.ExpiresAt)
}

// SetExpiration 设置过期时间
func (e *Expirable) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().Add(duration)
	e.ExpiresAt = &expiresAt
}

// TargetEntity 关联实体混入
type TargetEntity struct {
	TargetType string `bson:"target_type" json:"targetType"` // 目标类型
	TargetID   string `bson:"target_id" json:"targetId"`     // 目标ID
}

// Pinned 置顶状态混入
type Pinned struct {
	IsPinned bool       `bson:"is_pinned" json:"isPinned"`                 // 是否置顶
	PinnedAt *time.Time `bson:"pinned_at,omitempty" json:"pinnedAt,omitempty"` // 置顶时间
	PinnedBy *string    `bson:"pinned_by,omitempty" json:"pinnedBy,omitempty"`   // 置顶操作者ID
}

// Pin 置顶
func (p *Pinned) Pin(operatorID string) {
	p.IsPinned = true
	now := time.Now()
	p.PinnedAt = &now
	p.PinnedBy = &operatorID
}

// Unpin 取消置顶
func (p *Pinned) Unpin() {
	p.IsPinned = false
	p.PinnedAt = nil
	p.PinnedBy = nil
}

// TitledEntity 标题实体混入
type TitledEntity struct {
	Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`
}
