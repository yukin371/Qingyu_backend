package shared

import "time"

// CommunicationBase 通信基础实体混入
type CommunicationBase struct {
	SenderID   string     `bson:"sender_id" json:"senderId" validate:"required"`
	ReceiverID string     `bson:"receiver_id" json:"receiverId" validate:"required"`
	IsRead     bool       `bson:"is_read" json:"isRead"`
	ReadAt     *time.Time `bson:"read_at,omitempty" json:"readAt,omitempty"`
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
