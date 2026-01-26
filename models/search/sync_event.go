package search

import "time"

// SyncEventType 同步事件类型
type SyncEventType string

const (
	SyncEventInsert SyncEventType = "insert"
	SyncEventUpdate SyncEventType = "update"
	SyncEventDelete SyncEventType = "delete"
)

// SyncEvent 同步事件
type SyncEvent struct {
	ID           string                 `json:"id"`                      // 文档 ID
	Type         SyncEventType          `json:"type"`                    // 事件类型
	Index        string                 `json:"index"`                   // 目标索引
	OpType       SyncEventType          `json:"op_type"`                 // 操作类型（insert/update/delete）
	ChangedFields []string              `json:"changed_fields,omitempty"` // 变更字段列表
	FullDocument map[string]interface{} `json:"doc,omitempty"`           // 完整文档（可选）
	Timestamp    time.Time              `json:"timestamp"`               // 时间戳
}

// BookSyncEvent 书籍同步事件
type BookSyncEvent struct {
	Event SyncEvent
	Book  interface{} // Book 数据
}

// ProjectSyncEvent 项目同步事件
type ProjectSyncEvent struct {
	Event   SyncEvent
	Project interface{} // Project 数据
}

// DocumentSyncEvent 文档同步事件
type DocumentSyncEvent struct {
	Event    SyncEvent
	Document interface{} // Document 数据
}

// UserSyncEvent 用户同步事件
type UserSyncEvent struct {
	Event SyncEvent
	User  interface{} // User 数据
}
