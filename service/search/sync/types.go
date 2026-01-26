package sync

import (
	"time"

	"Qingyu_backend/models/search"
)

// DeadLetterEvent 死信队列事件
type DeadLetterEvent struct {
	Event         *search.SyncEvent `json:"event"`          // 原始事件
	ErrorMessage  string            `json:"error_message"`  // 错误信息
	RetryCount    int               `json:"retry_count"`    // 重试次数
	LastRetryAt   int64             `json:"last_retry_at"`  // 最后重试时间
	OriginalEvent *search.SyncEvent `json:"original_event"` // 原始事件（用于追踪）
}

// ConsistencyReport 一致性校验报告
type ConsistencyReport struct {
	ID          string                 `json:"id"`
	Collection  string                 `json:"collection"`
	MongoCount  int64                  `json:"mongo_count"`
	ESCount     int64                  `json:"es_count"`
	MissingDocs []string               `json:"missing_docs"` // MongoDB 有但 ES 没有的文档 ID
	ExtraDocs   []string               `json:"extra_docs"`   // ES 有但 MongoDB 没有的文档 ID
	Status      string                 `json:"status"`       // "consistent", "inconsistent", "error"
	CheckedAt   time.Time              `json:"checked_at"`
	Details     map[string]interface{} `json:"details"`
}
