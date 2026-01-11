package interfaces

import (
	"context"

	"Qingyu_backend/pkg/sync"
	"Qingyu_backend/pkg/websocket"
)

// ProgressSyncService 阅读进度同步服务接口
type ProgressSyncService interface {
	GetHub() *websocket.ProgressHub
	SyncProgress(ctx context.Context, userID, bookID, chapterID, deviceID string, progress float64) error
	MergeOfflineProgresses(ctx context.Context, userID string, progresses []sync.OfflineProgress) error
	GetSyncStatus(userID string) *sync.SyncStatus
}
