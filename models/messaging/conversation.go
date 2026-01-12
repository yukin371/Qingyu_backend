package messaging

import (
	"time"

	"Qingyu_backend/models/messaging/base"
)

// Conversation 对话模型（用户之间的私信会话）
type Conversation struct {
	base.IdentifiedEntity `bson:",inline"`
	base.Timestamps       `bson:",inline"`

	// 对话参与者
	ParticipantIDs []string `bson:"participant_ids" json:"participantIds" validate:"required,min=2"` // 参与者ID列表
	CreatedBy      string   `bson:"created_by" json:"createdBy" validate:"required"`                 // 创建者ID

	// 对话信息
	LastMessageID      *string    `bson:"last_message_id,omitempty" json:"lastMessageId,omitempty"`           // 最后一条消息ID
	LastMessageAt      *time.Time `bson:"last_message_at,omitempty" json:"lastMessageAt,omitempty"`           // 最后消息时间
	LastMessagePreview string     `bson:"last_message_preview,omitempty" json:"lastMessagePreview,omitempty"` // 最后消息预览
	LastMessageBy      *string    `bson:"last_message_by,omitempty" json:"lastMessageBy,omitempty"`           // 最后消息发送者

	// 对话状态
	IsActive   bool            `bson:"is_active" json:"isActive"`                         // 是否活跃
	IsArchived bool            `bson:"is_archived" json:"isArchived"`                     // 是否已归档
	ArchivedBy map[string]bool `bson:"archived_by,omitempty" json:"archivedBy,omitempty"` // 归档状态（按用户）
	MutedBy    map[string]bool `bson:"muted_by,omitempty" json:"mutedBy,omitempty"`       // 免打扰状态（按用户）

	// 对话类型
	Type      ConversationType       `bson:"type" json:"type" validate:"required,oneof=direct group"` // 对话类型
	GroupInfo *ConversationGroupInfo `bson:"group_info,omitempty" json:"groupInfo,omitempty"`         // 群组信息（仅type=group时）

	// 参与者快照（防止改名后历史消息显示问题）
	ParticipantSnapshots map[string]*ConversationParticipantSnapshot `bson:"participant_snapshots,omitempty" json:"participantSnapshots,omitempty"`
}

// ConversationType 对话类型
type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct" // 单聊
	ConversationTypeGroup  ConversationType = "group"  // 群聊
)

// ConversationGroupInfo 群组信息
type ConversationGroupInfo struct {
	Name        string   `bson:"name" json:"name" validate:"required,min=1,max=50"`  // 群组名称
	Description string   `bson:"description,omitempty" json:"description,omitempty"` // 群组描述
	Avatar      string   `bson:"avatar,omitempty" json:"avatar,omitempty"`           // 群组头像
	OwnerID     string   `bson:"owner_id" json:"ownerId" validate:"required"`        // 群主ID
	AdminIDs    []string `bson:"admin_ids,omitempty" json:"adminIds,omitempty"`      // 管理员ID列表
	MaxMembers  int      `bson:"max_members" json:"maxMembers"`                      // 最大成员数
}

// ConversationParticipantSnapshot 对话参与者快照
type ConversationParticipantSnapshot struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}

// DirectMessage 私信消息
type DirectMessage struct {
	base.IdentifiedEntity     `bson:",inline"`
	base.Timestamps           `bson:",inline"`
	base.CommunicationBase    `bson:",inline"`
	base.ThreadedConversation `bson:",inline"`

	ConversationID string                 `bson:"conversation_id" json:"conversationId" validate:"required"`                // 所属对话ID
	Type           MessageType            `bson:"type" json:"type" validate:"required,oneof=text image file system notice"` // 消息类型
	Content        string                 `bson:"content" json:"content" validate:"required,min=1,max=5000"`                // 消息内容
	Status         MessageStatus          `bson:"status" json:"status" validate:"required,oneof=normal recalled deleted"`   // 消息状态
	Extra          map[string]interface{} `bson:"extra,omitempty" json:"extra,omitempty"`                                   // 额外数据
	SenderSnapshot *MessageSenderSnapshot `bson:"sender_snapshot,omitempty" json:"senderSnapshot,omitempty"`                // 发送者快照
}

// MessageType 消息类型
type MessageType string

const (
	MessageTypeText   MessageType = "text"   // 文本消息
	MessageTypeImage  MessageType = "image"  // 图片消息
	MessageTypeFile   MessageType = "file"   // 文件消息
	MessageTypeSystem MessageType = "system" // 系统消息
	MessageTypeNotice MessageType = "notice" // 通知消息
)

// MessageStatus 消息状态
type MessageStatus string

const (
	MessageStatusNormal   MessageStatus = "normal"   // 正常
	MessageStatusRecalled MessageStatus = "recalled" // 已撤回
	MessageStatusDeleted  MessageStatus = "deleted"  // 已删除
)

// MessageSenderSnapshot 发送者快照
type MessageSenderSnapshot struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}

// ConversationFilter 对话查询过滤器
type ConversationFilter struct {
	UserID     *string           `json:"userId,omitempty"`
	Type       *ConversationType `json:"type,omitempty"`
	IsActive   *bool             `json:"isActive,omitempty"`
	IsArchived *bool             `json:"isArchived,omitempty"`
	StartDate  *time.Time        `json:"startDate,omitempty"`
	EndDate    *time.Time        `json:"endDate,omitempty"`
	SortBy     string            `json:"sortBy,omitempty"`
	SortOrder  string            `json:"sortOrder,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f *ConversationFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	// 必须指定用户ID
	if f.UserID != nil {
		conditions["participant_ids"] = *f.UserID
	}

	if f.Type != nil {
		conditions["type"] = *f.Type
	}
	if f.IsActive != nil {
		conditions["is_active"] = *f.IsActive
	}
	if f.IsArchived != nil {
		conditions["is_archived"] = *f.IsArchived
	}

	// 时间范围
	if f.StartDate != nil || f.EndDate != nil {
		timeCondition := make(map[string]interface{})
		if f.StartDate != nil {
			timeCondition["$gte"] = *f.StartDate
		}
		if f.EndDate != nil {
			timeCondition["$lte"] = *f.EndDate
		}
		conditions["created_at"] = timeCondition
	}

	return conditions
}

// GetSort 获取排序
func (f *ConversationFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	sortValue := -1 // 默认降序
	if f.SortOrder == "asc" {
		sortValue = 1
	}

	switch f.SortBy {
	case "last_message_at":
		sort["last_message_at"] = sortValue
	case "created_at":
		sort["created_at"] = sortValue
	default:
		sort["last_message_at"] = -1 // 默认按最后消息时间降序
	}

	return sort
}

// GetLimit 获取限制
func (f *ConversationFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 获取偏移
func (f *ConversationFilter) GetOffset() int {
	return f.Offset
}

// GetFields 获取字段
func (f *ConversationFilter) GetFields() []string {
	return []string{}
}

// Validate 验证
func (f *ConversationFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	if f.UserID == nil {
		return &ValidationError{Message: "必须指定用户ID"}
	}
	return nil
}

// HasParticipant 检查对话是否包含指定参与者
func (c *Conversation) HasParticipant(userID string) bool {
	for _, pid := range c.ParticipantIDs {
		if pid == userID {
			return true
		}
	}
	return false
}

// IsGroup 判断是否为群聊
func (c *Conversation) IsGroup() bool {
	return c.Type == ConversationTypeGroup
}

// IsDirect 判断是否为单聊
func (c *Conversation) IsDirect() bool {
	return c.Type == ConversationTypeDirect
}

// GetOtherParticipantID 获取单聊中的对方ID
func (c *Conversation) GetOtherParticipantID(userID string) (string, bool) {
	if c.IsGroup() {
		return "", false
	}

	for _, pid := range c.ParticipantIDs {
		if pid != userID {
			return pid, true
		}
	}

	return "", false
}

// Archive 归档对话（针对指定用户）
func (c *Conversation) Archive(userID string) {
	if c.ArchivedBy == nil {
		c.ArchivedBy = make(map[string]bool)
	}
	c.ArchivedBy[userID] = true
	c.Touch()
}

// Unarchive 取消归档（针对指定用户）
func (c *Conversation) Unarchive(userID string) {
	if c.ArchivedBy != nil {
		delete(c.ArchivedBy, userID)
	}
	c.Touch()
}

// IsArchivedByUser 判断对话是否被指定用户归档
func (c *Conversation) IsArchivedByUser(userID string) bool {
	if c.ArchivedBy == nil {
		return false
	}
	return c.ArchivedBy[userID]
}

// Mute 设置免打扰（针对指定用户）
func (c *Conversation) Mute(userID string) {
	if c.MutedBy == nil {
		c.MutedBy = make(map[string]bool)
	}
	c.MutedBy[userID] = true
	c.Touch()
}

// Unmute 取消免打扰（针对指定用户）
func (c *Conversation) Unmute(userID string) {
	if c.MutedBy != nil {
		delete(c.MutedBy, userID)
	}
	c.Touch()
}

// IsMutedByUser 判断对话是否被指定用户免打扰
func (c *Conversation) IsMutedByUser(userID string) bool {
	if c.MutedBy == nil {
		return false
	}
	return c.MutedBy[userID]
}

// AddParticipant 添加参与者（仅群聊）
func (c *Conversation) AddParticipant(userID string, snapshot *ConversationParticipantSnapshot) bool {
	if !c.IsGroup() {
		return false
	}

	// 检查是否已存在
	if c.HasParticipant(userID) {
		return false
	}

	// 检查人数限制
	if c.GroupInfo != nil && c.GroupInfo.MaxMembers > 0 {
		if len(c.ParticipantIDs) >= c.GroupInfo.MaxMembers {
			return false
		}
	}

	c.ParticipantIDs = append(c.ParticipantIDs, userID)
	if c.ParticipantSnapshots == nil {
		c.ParticipantSnapshots = make(map[string]*ConversationParticipantSnapshot)
	}
	c.ParticipantSnapshots[userID] = snapshot
	c.Touch()

	return true
}

// RemoveParticipant 移除参与者（仅群聊）
func (c *Conversation) RemoveParticipant(userID string) bool {
	if !c.IsGroup() {
		return false
	}

	for i, pid := range c.ParticipantIDs {
		if pid == userID {
			c.ParticipantIDs = append(c.ParticipantIDs[:i], c.ParticipantIDs[i+1:]...)
			if c.ParticipantSnapshots != nil {
				delete(c.ParticipantSnapshots, userID)
			}
			c.Touch()
			return true
		}
	}

	return false
}

// UpdateLastMessage 更新最后消息信息
func (c *Conversation) UpdateLastMessage(messageID string, content string, senderID string, timestamp time.Time) {
	c.LastMessageID = &messageID
	c.LastMessageAt = &timestamp
	c.LastMessageBy = &senderID

	// 生成预览（最多50字符）
	if len(content) > 50 {
		content = content[:50] + "..."
	}
	c.LastMessagePreview = content

	c.Touch()
}

// DirectMessageFilter 消息查询过滤器
type DirectMessageFilter struct {
	ConversationID *string        `json:"conversationId,omitempty"`
	SenderID       *string        `json:"senderId,omitempty"`
	ReceiverID     *string        `json:"receiverId,omitempty"`
	Type           *MessageType   `json:"type,omitempty"`
	Status         *MessageStatus `json:"status,omitempty"`
	IsRead         *bool          `json:"isRead,omitempty"`
	StartTime      *time.Time     `json:"startTime,omitempty"`
	EndTime        *time.Time     `json:"endTime,omitempty"`
	SortBy         string         `json:"sortBy,omitempty"`
	SortOrder      string         `json:"sortOrder,omitempty"`
	Limit          int            `json:"limit,omitempty"`
	Offset         int            `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f *DirectMessageFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.ConversationID != nil {
		conditions["conversation_id"] = *f.ConversationID
	}
	if f.SenderID != nil {
		conditions["sender_id"] = *f.SenderID
	}
	if f.ReceiverID != nil {
		conditions["receiver_id"] = *f.ReceiverID
	}
	if f.Type != nil {
		conditions["type"] = *f.Type
	}
	if f.Status != nil {
		conditions["status"] = *f.Status
	}
	if f.IsRead != nil {
		conditions["is_read"] = *f.IsRead
	}

	// 时间范围
	if f.StartTime != nil || f.EndTime != nil {
		timeCondition := make(map[string]interface{})
		if f.StartTime != nil {
			timeCondition["$gte"] = *f.StartTime
		}
		if f.EndTime != nil {
			timeCondition["$lte"] = *f.EndTime
		}
		conditions["created_at"] = timeCondition
	}

	// 排除已删除的消息
	if f.Status == nil {
		conditions["status"] = map[string]interface{}{"$ne": MessageStatusDeleted}
	}

	return conditions
}

// GetSort 获取排序
func (f *DirectMessageFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	sortValue := -1 // 默认降序
	if f.SortOrder == "asc" {
		sortValue = 1
	}

	switch f.SortBy {
	case "created_at":
		sort["created_at"] = sortValue
	default:
		sort["created_at"] = -1
	}

	return sort
}

// GetLimit 获取限制
func (f *DirectMessageFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 获取偏移
func (f *DirectMessageFilter) GetOffset() int {
	return f.Offset
}

// GetFields 获取字段
func (f *DirectMessageFilter) GetFields() []string {
	return []string{}
}

// Validate 验证
func (f *DirectMessageFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	return nil
}

// IsRecalled 判断消息是否已撤回
func (m *DirectMessage) IsRecalled() bool {
	return m.Status == MessageStatusRecalled
}

// Recall 撤回消息
func (m *DirectMessage) Recall() bool {
	// 消息发送超过2分钟不能撤回
	if time.Since(m.CreatedAt) > 2*time.Minute {
		return false
	}
	m.Status = MessageStatusRecalled
	m.Touch()
	return true
}

// CanDelete 判断是否可以删除
func (m *DirectMessage) CanDelete(userID string) bool {
	return m.SenderID == userID
}

// SoftDelete 软删除消息
func (m *DirectMessage) SoftDelete(userID string) bool {
	if !m.CanDelete(userID) {
		return false
	}
	m.Status = MessageStatusDeleted
	m.Touch()
	return true
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
