package sync

import (
	"context"
	"time"

	"Qingyu_backend/models/search"
)

// SyncWorker 同步工作器接口
type SyncWorker interface {
	// Start 启动同步工作器
	Start(ctx context.Context) error

	// Stop 停止同步工作器
	Stop() error

	// ProcessEvent 处理同步事件
	ProcessEvent(ctx context.Context, event *search.SyncEvent) error

	// Status 获取同步状态
	Status() (*SyncStatus, error)
}

// ChangeStreamListener 变更流监听器接口
type ChangeStreamListener interface {
	// Watch 监听指定集合的变更
	Watch(ctx context.Context, db, collection, index string) error

	// ProcessChange 处理变更事件
	ProcessChange(ctx context.Context, change *ChangeData) error

	// Stop 停止监听
	Stop() error
}

// ChangeData 变更数据
type ChangeData struct {
	ID           string                 `json:"id"`           // 文档 ID
	Operation    string                 `json:"operation"`    // 操作类型：insert, update, delete
	Database     string                 `json:"database"`     // 数据库名
	Collection   string                 `json:"collection"`   // 集合名
	DocumentKey  map[string]interface{} `json:"document_key"` // 文档键
	FullDocument map[string]interface{} `json:"doc"`          // 完整文档
	UpdateTime   time.Time              `json:"update_time"`  // 更新时间
	// update 操作时的变更字段
	UpdatedFields map[string]interface{} `json:"updated_fields,omitempty"`
}

// SyncStatus 同步状态
type SyncStatus struct {
	Running       bool      `json:"running"`        // 是否运行中
	LastSyncTime  time.Time `json:"last_sync_time"` // 最后同步时间
	TotalEvents   int64     `json:"total_events"`   // 总事件数
	SuccessEvents int64     `json:"success_events"` // 成功事件数
	FailedEvents  int64     `json:"failed_events"`  // 失败事件数
	QueuedEvents  int64     `json:"queued_events"`  // 队列中事件数
}

// SyncConfig 同步配置
type SyncConfig struct {
	// 批量处理大小
	BatchSize int
	// 批量处理间隔
	BatchInterval time.Duration
	// 重试次数
	RetryCount int
	// 重试间隔
	RetryInterval time.Duration
	// 队列大小
	QueueSize int
	// 是否启用同步
	Enabled bool
	// 监听的数据库和集合
	Collections []CollectionConfig
}

// CollectionConfig 集合配置
type CollectionConfig struct {
	Database   string `json:"database"`   // 数据库名
	Collection string `json:"collection"` // 集合名
	Index      string `json:"index"`      // 目标索引
	Enabled    bool   `json:"enabled"`    // 是否启用
}

// SyncEventProcessor 同步事件处理器接口
type SyncEventProcessor interface {
	// ProcessInsert 处理插入事件
	ProcessInsert(ctx context.Context, event *search.SyncEvent) error

	// ProcessUpdate 处理更新事件
	ProcessUpdate(ctx context.Context, event *search.SyncEvent) error

	// ProcessDelete 处理删除事件
	ProcessDelete(ctx context.Context, event *search.SyncEvent) error

	// ProcessBatch 批量处理事件
	ProcessBatch(ctx context.Context, events []*search.SyncEvent) error
}
