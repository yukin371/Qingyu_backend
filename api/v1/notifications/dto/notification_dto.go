package dto

import "time"

// MarkAsReadRequest 标记已读请求
type MarkAsReadRequest struct {
	ReadAt int64 `json:"read_at" binding:"required"` // 阅读时间戳
}

// BatchMarkReadRequest 批量标记已读请求
type BatchMarkReadRequest struct {
	NotificationIDs []string `json:"notification_ids" binding:"required,min=1"`
	ReadAt          int64    `json:"read_at" binding:"required"`
}

// ReadAllRequest 全部标记已读请求
type ReadAllRequest struct {
	ReadAt int64 `json:"read_at" binding:"required"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	NotificationIDs []string `json:"notification_ids" binding:"required,min=1"`
}

// ResendNotificationRequest 重新发送通知请求
type ResendNotificationRequest struct {
	Method string `json:"method" binding:"required,oneof=email push"` // 重新发送方式
}

// MarkAsReadResponse 标记已读响应
type MarkAsReadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	Success    bool     `json:"success"`
	Total      int      `json:"total"`
	Succeeded  int      `json:"succeeded"`
	Failed     int      `json:"failed"`
	Errors     []string `json:"errors,omitempty"`
}

// WSEndpointResponse WebSocket端点响应
type WSEndpointResponse struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

// ListNotificationsResponse 通知列表响应
type ListNotificationsResponse struct {
	Notifications []NotificationItem `json:"notifications"`
	Total         int                `json:"total"`
	UnreadCount   int                `json:"unread_count"`
}

// NotificationItem 通知项
type NotificationItem struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`       // system, comment, like, follow, etc.
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Data      interface{} `json:"data"`       // 额外数据
	Read      bool        `json:"read"`
	CreatedAt time.Time   `json:"created_at"`
}
