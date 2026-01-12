package sync

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/pkg/websocket"
	"Qingyu_backend/service/reader"
)

// ProgressSyncService 阅读进度同步服务
type ProgressSyncService struct {
	readerService *reader.ReaderService
	hub           *websocket.ProgressHub
}

// NewProgressSyncService 创建进度同步服务
func NewProgressSyncService(readerService *reader.ReaderService) *ProgressSyncService {
	return &ProgressSyncService{
		readerService: readerService,
		hub:           websocket.NewProgressHub(),
	}
}

// GetHub 获取WebSocket Hub
func (s *ProgressSyncService) GetHub() *websocket.ProgressHub {
	return s.hub
}

// Start 启动同步服务
func (s *ProgressSyncService) Start() {
	go s.hub.Run()
}

// SyncProgress 同步阅读进度
func (s *ProgressSyncService) SyncProgress(ctx context.Context, userID, bookID, chapterID, deviceID string, progress float64) error {
	// 1. 验证用户ID格式
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// 2. 获取当前阅读进度
	currentProgress, err := s.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		// 如果没有现有进度，直接保存新进度
		return s.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)
	}

	// 3. 检查时间戳，使用最新的进度
	now := time.Now()
	if currentProgress.LastReadAt.After(now.Add(-5 * time.Second)) {
		// 如果当前进度是最近更新的（5秒内），可能存在冲突
		// 使用较新的进度
		if currentProgress.LastReadAt.Before(now) {
			return s.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)
		}
		// 否则忽略旧进度
		return nil
	}

	// 4. 保存新进度
	if err := s.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, progress); err != nil {
		return err
	}

	// 5. 通过WebSocket广播给用户的其他设备
	message := &websocket.ProgressMessage{
		Type:      "sync",
		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,
		Progress:  progress,
		Timestamp: now,
		DeviceID:  deviceID,
	}
	s.hub.SyncProgress(message)

	return nil
}

// MergeOfflineProgresses 合并离线进度
func (s *ProgressSyncService) MergeOfflineProgresses(ctx context.Context, userID string, progresses []OfflineProgress) error {
	// 验证用户ID格式
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// 按书籍分组进度
	bookProgresses := make(map[string][]OfflineProgress)
	for _, p := range progresses {
		bookProgresses[p.BookID] = append(bookProgresses[p.BookID], p)
	}

	// 处理每本书的进度
	for bookID, bookProg := range bookProgresses {
		// 找到最新的进度
		var latest *OfflineProgress
		for i := range bookProg {
			if latest == nil || bookProg[i].Timestamp.After(latest.Timestamp) {
				latest = &bookProg[i]
			}
		}

		if latest != nil {
			// 获取服务器当前进度
			currentProgress, err := s.readerService.GetReadingProgress(ctx, userID, bookID)
			if err != nil {
				// 服务器没有进度，直接使用离线进度
				s.readerService.SaveReadingProgress(ctx, userID, bookID, latest.ChapterID, latest.Progress)
				continue
			}

			// 使用时间戳较新的进度
			if latest.Timestamp.After(currentProgress.LastReadAt) {
				s.readerService.SaveReadingProgress(ctx, userID, bookID, latest.ChapterID, latest.Progress)
			}
		}
	}

	return nil
}

// GetSyncStatus 获取同步状态
func (s *ProgressSyncService) GetSyncStatus(userID string) *SyncStatus {
	connectedDevices := s.hub.GetConnectedDevices(userID)
	return &SyncStatus{
		UserID:           userID,
		ConnectedDevices: connectedDevices,
		DeviceCount:      len(connectedDevices),
		IsSyncing:        len(connectedDevices) > 1,
	}
}

// OfflineProgress 离线进度
type OfflineProgress struct {
	UserID    string    `json:"userId"`
	BookID    string    `json:"bookId"`
	ChapterID string    `json:"chapterId"`
	Progress  float64   `json:"progress"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceId"`
}

// SyncStatus 同步状态
type SyncStatus struct {
	UserID           string   `json:"userId"`
	ConnectedDevices []string `json:"connectedDevices"`
	DeviceCount      int      `json:"deviceCount"`
	IsSyncing        bool     `json:"isSyncing"`
}
