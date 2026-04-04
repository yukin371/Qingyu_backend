package events

import (
	"context"
	"fmt"
	"log"
	"time"

	baseInterfaces "Qingyu_backend/service/interfaces/base"
	"Qingyu_backend/service/search/engine"
)

// SearchIndexHandler 书籍发布事件索引处理器
// 订阅 project.published 事件，将发布后的书籍自动加入搜索索引
type SearchIndexHandler struct {
	name   string
	engine engine.Engine
}

// NewSearchIndexHandler 创建搜索索引处理器
func NewSearchIndexHandler(eng engine.Engine) *SearchIndexHandler {
	return &SearchIndexHandler{
		name:   "SearchIndexHandler",
		engine: eng,
	}
}

// Handle 处理 project.published 事件
func (h *SearchIndexHandler) Handle(ctx context.Context, event baseInterfaces.Event) error {
	if event == nil {
		return nil
	}

	eventType := event.GetEventType()
	if eventType != "project.published" {
		return nil
	}

	// 从事件数据中提取 bookstoreId（书籍 ID）
	data, ok := event.GetEventData().(map[string]interface{})
	if !ok {
		return fmt.Errorf("SearchIndexHandler: event data is not a map")
	}

	bookstoreID, ok := data["bookstoreId"].(string)
	if !ok || bookstoreID == "" {
		return fmt.Errorf("SearchIndexHandler: bookstoreId not found in event data")
	}

	// 直接用 bookstoreId 作为文档 ID 索引
	// 书籍的其他信息从事件中获取，构建搜索文档
	doc := h.buildSearchDocument(data, bookstoreID)

	// 尝试插入文档；如果已存在（索引重建场景），则静默忽略
	err := h.engine.Index(ctx, booksIndexName, []engine.Document{doc})
	if err != nil {
		// 如果是重复键错误，说明文档已存在，不需要处理
		if isDuplicateKeyError(err) {
			log.Printf("[SearchIndexHandler] Book %s already indexed, skipping", bookstoreID)
			return nil
		}
		return fmt.Errorf("SearchIndexHandler: failed to index book %s: %w", bookstoreID, err)
	}

	log.Printf("[SearchIndexHandler] Successfully indexed book %s", bookstoreID)
	return nil
}

// buildSearchDocument 从事件数据构建搜索文档
func (h *SearchIndexHandler) buildSearchDocument(data map[string]interface{}, bookID string) engine.Document {
	source := make(map[string]interface{})

	// 从事件数据中提取可用的书籍信息
	if projectID, ok := data["projectId"].(string); ok {
		source["project_id"] = projectID
	}
	if bookstoreID, ok := data["bookstoreId"].(string); ok {
		source["bookstore_id"] = bookstoreID
	}
	if publishedAt, ok := data["publishedAt"].(time.Time); ok {
		source["published_at"] = publishedAt
	} else if publishedAtStr, ok := data["publishedAt"].(string); ok {
		if t, err := time.Parse(time.RFC3339, publishedAtStr); err == nil {
			source["published_at"] = t
		}
	}

	// 状态设为公开
	source["status"] = "ongoing"

	return engine.Document{
		ID:     bookID,
		Source: source,
	}
}

// GetHandlerName 返回处理器名称
func (h *SearchIndexHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 返回支持的事件类型
func (h *SearchIndexHandler) GetSupportedEventTypes() []string {
	return []string{"project.published"}
}

// booksIndexName 搜索索引名称，与 BookProvider 保持一致
const booksIndexName = "books"

// isDuplicateKeyError 判断是否为 MongoDB 重复键错误
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "duplicate key") || contains(errStr, "E11000")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
