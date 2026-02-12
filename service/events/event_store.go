package events

import (
	"context"
	"time"

	"Qingyu_backend/service/base"
)

// EventStore 事件存储接口
// 负责持久化事件，支持事件查询和回放
type EventStore interface {
	// Store 存储事件
	Store(ctx context.Context, event base.Event) error

	// StoreBatch 批量存储事件
	StoreBatch(ctx context.Context, events []base.Event) error

	// GetByID 根据ID获取事件
	GetByID(ctx context.Context, eventID string) (*StoredEvent, error)

	// GetByType 根据类型获取事件
	GetByType(ctx context.Context, eventType string, limit, offset int64) ([]*StoredEvent, error)

	// GetBySource 根据来源获取事件
	GetBySource(ctx context.Context, source string, limit, offset int64) ([]*StoredEvent, error)

	// GetByTimeRange 根据时间范围获取事件
	GetByTimeRange(ctx context.Context, start, end time.Time, limit, offset int64) ([]*StoredEvent, error)

	// GetByTypeAndTimeRange 根据类型和时间范围获取事件
	GetByTypeAndTimeRange(ctx context.Context, eventType string, start, end time.Time, limit, offset int64) ([]*StoredEvent, error)

	// Replay 事件回放（重放事件到处理器）
	Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) (*ReplayResult, error)

	// Cleanup 清理过期事件
	Cleanup(ctx context.Context, before time.Time) (int64, error)

	// Count 统计事件数量
	Count(ctx context.Context, filter EventFilter) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// StoredEvent 存储的事件
type StoredEvent struct {
	ID        string      `bson:"_id" json:"id"`
	EventType string      `bson:"event_type" json:"event_type"`
	EventData interface{} `bson:"event_data" json:"event_data"`
	Timestamp time.Time   `bson:"timestamp" json:"timestamp"`
	Source    string      `bson:"source" json:"source"`
	Processed bool        `bson:"processed" json:"processed"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	ExpiresAt time.Time   `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// EventFilter 事件过滤器
type EventFilter struct {
	EventType string     `bson:"event_type,omitempty"`
	Source    string     `bson:"source,omitempty"`
	StartTime *time.Time `bson:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty"`
	Processed *bool      `bson:"processed,omitempty"`
	Limit     int64      `bson:"limit,omitempty"`
	Offset    int64      `bson:"offset,omitempty"`
	DryRun    bool       `bson:"dry_run,omitempty"` // DryRun模式：只统计不执行处理器
}

// ReplayResult 事件回放结果
type ReplayResult struct {
	// ReplayedCount 成功重放的事件数量
	ReplayedCount int64

	// FailedCount 失败的事件数量
	FailedCount int64

	// SkippedCount 跳过的事件数量（DryRun模式）
	SkippedCount int64

	// Duration 执行耗时
	Duration time.Duration
}

// EventStoreConfig 事件存储配置
type EventStoreConfig struct {
	// 是否启用持久化
	Enabled bool

	// 存储类型: mongodb/redis
	StorageType string

	// 事件保留时间（TTL）
	RetentionDuration time.Duration

	// 是否启用压缩
	CompressEnabled bool

	// 批量写入大小
	BatchSize int

	// 写入超时
	WriteTimeout time.Duration
}

// DefaultEventStoreConfig 默认配置
func DefaultEventStoreConfig() *EventStoreConfig {
	return &EventStoreConfig{
		Enabled:           true,
		StorageType:       "mongodb",
		RetentionDuration: 30 * 24 * time.Hour, // 30天
		CompressEnabled:   false,
		BatchSize:         100,
		WriteTimeout:      5 * time.Second,
	}
}

// EventReplayer 事件回放器接口
type EventReplayer interface {
	// ReplayFromTimestamp 从指定时间戳回放事件
	ReplayFromTimestamp(ctx context.Context, timestamp time.Time, handler base.EventHandler) error

	// ReplayFromEventID 从指定事件ID回放事件
	ReplayFromEventID(ctx context.Context, eventID string, handler base.EventHandler) error

	// ReplayWithType 回放指定类型的事件
	ReplayWithType(ctx context.Context, eventType string, handler base.EventHandler) error
}

// EventSnapshot 事件快照
// 用于定期保存聚合根的状态
type EventSnapshot struct {
	ID            string      `bson:"_id" json:"id"`
	AggregateID   string      `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType string      `bson:"aggregate_type" json:"aggregate_type"`
	Version       int64       `bson:"version" json:"version"`
	State         interface{} `bson:"state" json:"state"`
	CreatedAt     time.Time   `bson:"created_at" json:"created_at"`
}

// EventSnapshotStore 事件快照存储接口
type EventSnapshotStore interface {
	// Save 保存快照
	Save(ctx context.Context, snapshot *EventSnapshot) error

	// Get 获取最新快照
	Get(ctx context.Context, aggregateID, aggregateType string) (*EventSnapshot, error)

	// GetByVersion 获取指定版本的快照
	GetByVersion(ctx context.Context, aggregateID, aggregateType string, version int64) (*EventSnapshot, error)
}
