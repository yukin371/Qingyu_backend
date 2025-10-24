package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 阅读相关事件 ============

// ReadingEvent 事件类型常量
const (
	EventTypeChapterRead      = "reading.chapter_read"
	EventTypeBookmarkAdded    = "reading.bookmark_added"
	EventTypeNoteCreated      = "reading.note_created"
	EventTypeReadingProgress  = "reading.progress_updated"
	EventTypeReadingCompleted = "reading.book_completed"
)

// ReadingEventData 阅读事件数据
type ReadingEventData struct {
	UserID    string                 `json:"user_id"`
	BookID    string                 `json:"book_id"`
	ChapterID string                 `json:"chapter_id,omitempty"`
	Action    string                 `json:"action"`
	Progress  int                    `json:"progress,omitempty"`
	Time      time.Time              `json:"time"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewChapterReadEvent 创建章节阅读事件
func NewChapterReadEvent(userID, bookID, chapterID string) base.Event {
	return &base.BaseEvent{
		EventType: EventTypeChapterRead,
		EventData: ReadingEventData{
			UserID:    userID,
			BookID:    bookID,
			ChapterID: chapterID,
			Action:    "chapter_read",
			Time:      time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "ReaderService",
	}
}

// NewReadingProgressEvent 创建阅读进度事件
func NewReadingProgressEvent(userID, bookID string, progress int) base.Event {
	return &base.BaseEvent{
		EventType: EventTypeReadingProgress,
		EventData: ReadingEventData{
			UserID:   userID,
			BookID:   bookID,
			Action:   "progress_updated",
			Progress: progress,
			Time:     time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "ReaderService",
	}
}

// ============ 事件处理器 ============

// ReadingStatisticsHandler 阅读统计处理器
// 更新阅读统计数据
type ReadingStatisticsHandler struct {
	name string
}

// NewReadingStatisticsHandler 创建阅读统计处理器
func NewReadingStatisticsHandler() *ReadingStatisticsHandler {
	return &ReadingStatisticsHandler{
		name: "ReadingStatisticsHandler",
	}
}

// Handle 处理事件
func (h *ReadingStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	// 解析事件数据
	data, ok := event.GetEventData().(ReadingEventData)
	if !ok {
		return fmt.Errorf("事件数据类型错误")
	}

	// 更新统计信息
	switch event.GetEventType() {
	case EventTypeChapterRead:
		log.Printf("[ReadingStatistics] 用户 %s 阅读了章节 %s", data.UserID, data.ChapterID)
		// 实际项目中这里应该更新：
		// - 书籍阅读次数
		// - 用户阅读时长
		// - 章节热度

	case EventTypeReadingProgress:
		log.Printf("[ReadingStatistics] 用户 %s 的阅读进度更新为 %d%%", data.UserID, data.Progress)
		// 更新用户阅读进度统计

	case EventTypeReadingCompleted:
		log.Printf("[ReadingStatistics] 用户 %s 完成了书籍 %s", data.UserID, data.BookID)
		// 更新书籍完成统计
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *ReadingStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *ReadingStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventTypeChapterRead,
		EventTypeReadingProgress,
		EventTypeReadingCompleted,
	}
}

// RecommendationUpdateHandler 推荐更新处理器
// 根据阅读行为更新推荐
type RecommendationUpdateHandler struct {
	name string
}

// NewRecommendationUpdateHandler 创建推荐更新处理器
func NewRecommendationUpdateHandler() *RecommendationUpdateHandler {
	return &RecommendationUpdateHandler{
		name: "RecommendationUpdateHandler",
	}
}

// Handle 处理事件
func (h *RecommendationUpdateHandler) Handle(ctx context.Context, event base.Event) error {
	// 解析事件数据
	data, ok := event.GetEventData().(ReadingEventData)
	if !ok {
		return fmt.Errorf("事件数据类型错误")
	}

	// 根据阅读行为更新推荐
	log.Printf("[RecommendationUpdate] 基于用户 %s 的阅读行为更新推荐列表", data.UserID)

	// 实际项目中这里应该:
	// 1. 分析用户阅读偏好
	// 2. 更新用户画像
	// 3. 重新计算推荐书籍
	// recommendationService.UpdateUserRecommendations(data.UserID)

	return nil
}

// GetHandlerName 获取处理器名称
func (h *RecommendationUpdateHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *RecommendationUpdateHandler) GetSupportedEventTypes() []string {
	return []string{
		EventTypeChapterRead,
		EventTypeReadingCompleted,
	}
}
