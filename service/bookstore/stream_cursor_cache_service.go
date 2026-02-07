package bookstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Qingyu_backend/models/bookstore"

	"github.com/redis/go-redis/v9"
)

// StreamCursorCacheService 流式游标缓存服务
// 用于缓存游标状态，实现断点续传和性能优化
type StreamCursorCacheService struct {
	redisClient *redis.Client
	prefix      string
	defaultTTL  time.Duration
}

// CursorCacheEntry 游标缓存条目
type CursorCacheEntry struct {
	Cursor     string    `json:"cursor"`
	SearchKey  string    `json:"search_key"`
	Timestamp  time.Time `json:"timestamp"`
	BookCount  int       `json:"book_count"`
	LastBookID string    `json:"last_book_id,omitempty"`
}

// NewStreamCursorCacheService 创建流式游标缓存服务
func NewStreamCursorCacheService(redisClient *redis.Client) *StreamCursorCacheService {
	return &StreamCursorCacheService{
		redisClient: redisClient,
		prefix:      "stream_cursor:",
		defaultTTL:  10 * time.Minute, // 默认10分钟过期
	}
}

// SaveCursorState 保存游标状态
func (s *StreamCursorCacheService) SaveCursorState(
	ctx context.Context,
	searchKey string,
	cursor string,
	books []*bookstore.Book,
	ttl time.Duration,
) error {
	if ttl <= 0 {
		ttl = s.defaultTTL
	}

	var lastBookID string
	if len(books) > 0 {
		lastBookID = books[len(books)-1].ID.Hex()
	}

	entry := CursorCacheEntry{
		Cursor:     cursor,
		SearchKey:  searchKey,
		Timestamp:  time.Now(),
		BookCount:  len(books),
		LastBookID: lastBookID,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal cursor cache entry failed: %w", err)
	}

	key := s.prefix + searchKey
	if err := s.redisClient.Set(ctx, key, string(data), ttl).Err(); err != nil {
		return fmt.Errorf("save cursor state failed: %w", err)
	}

	return nil
}

// GetCursorState 获取游标状态
func (s *StreamCursorCacheService) GetCursorState(
	ctx context.Context,
	searchKey string,
) (*CursorCacheEntry, error) {
	key := s.prefix + searchKey
	cmd := s.redisClient.Get(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	data, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	if data == "" {
		return nil, nil // 游标不存在或已过期
	}

	var entry CursorCacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, fmt.Errorf("unmarshal cursor cache entry failed: %w", err)
	}

	return &entry, nil
}

// InvalidateCursor 使游标失效
func (s *StreamCursorCacheService) InvalidateCursor(
	ctx context.Context,
	searchKey string,
) error {
	key := s.prefix + searchKey
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("invalidate cursor failed: %w", err)
	}
	return nil
}

// GenerateSearchKey 生成搜索的唯一键
func (s *StreamCursorCacheService) GenerateSearchKey(
	keyword string,
	categoryID string,
	author string,
	status string,
	tags []string,
	sortBy string,
	sortOrder string,
) string {
	// 创建一个规范化的搜索键
	// 可以使用哈希函数来减少键的长度
	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s",
		keyword,
		categoryID,
		author,
		status,
		tags,
		sortBy,
		sortOrder,
	)
	return key
}

// ExtendCursorTTL 延长游标有效期
func (s *StreamCursorCacheService) ExtendCursorTTL(
	ctx context.Context,
	searchKey string,
	additionalTTL time.Duration,
) error {
	entry, err := s.GetCursorState(ctx, searchKey)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("cursor not found")
	}

	// 重新保存以延长TTL
	return s.SaveCursorState(ctx, searchKey, entry.Cursor, nil, additionalTTL)
}

// BatchInvalidateCursors 批量使游标失效
func (s *StreamCursorCacheService) BatchInvalidateCursors(
	ctx context.Context,
	searchKeys []string,
) error {
	if len(searchKeys) == 0 {
		return nil
	}

	keys := make([]string, len(searchKeys))
	for i, key := range searchKeys {
		keys[i] = s.prefix + key
	}

	if err := s.redisClient.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("batch invalidate cursors failed: %w", err)
	}

	return nil
}

// CleanupExpiredCursors 清理过期的游标（Redis会自动处理，这个方法用于显式清理）
func (s *StreamCursorCacheService) CleanupExpiredCursors(ctx context.Context) error {
	// Redis会自动处理过期键，这个方法可以用于特殊情况下的手动清理
	// 实现取决于具体的Redis客户端
	return nil
}

// GetActiveCursorCount 获取活跃游标数量（统计）
func (s *StreamCursorCacheService) GetActiveCursorCount(ctx context.Context) (int64, error) {
	// 这个方法需要Redis的SCAN命令来统计匹配的键
	// 实现取决于具体的Redis客户端
	return 0, nil
}

// ValidateCursor 验证游标是否有效
func (s *StreamCursorCacheService) ValidateCursor(
	ctx context.Context,
	searchKey string,
	cursor string,
) bool {
	entry, err := s.GetCursorState(ctx, searchKey)
	if err != nil || entry == nil {
		return false
	}

	return entry.Cursor == cursor
}

// UpdateCursorProgress 更新游标进度（用于长连接场景）
func (s *StreamCursorCacheService) UpdateCursorProgress(
	ctx context.Context,
	searchKey string,
	additionalBooks int,
) error {
	entry, err := s.GetCursorState(ctx, searchKey)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("cursor not found")
	}

	// 更新书籍数量
	entry.BookCount += additionalBooks
	entry.Timestamp = time.Now()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal cursor cache entry failed: %w", err)
	}

	key := s.prefix + searchKey
	// 保持原有的TTL
	if err := s.redisClient.Set(ctx, key, string(data), s.defaultTTL).Err(); err != nil {
		return fmt.Errorf("update cursor progress failed: %w", err)
	}

	return nil
}
