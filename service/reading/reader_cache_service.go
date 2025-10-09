package reading

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Qingyu_backend/models/reading/reader"

	"github.com/redis/go-redis/v9"
)

// ReaderCacheService 阅读器缓存服务接口
type ReaderCacheService interface {
	// 章节内容缓存
	GetChapterContent(ctx context.Context, chapterID string) (string, error)
	SetChapterContent(ctx context.Context, chapterID string, content string, expiration time.Duration) error
	InvalidateChapterContent(ctx context.Context, chapterID string) error

	// 章节信息缓存
	GetChapter(ctx context.Context, chapterID string) (*reader.Chapter, error)
	SetChapter(ctx context.Context, chapterID string, chapter *reader.Chapter, expiration time.Duration) error
	InvalidateChapter(ctx context.Context, chapterID string) error

	// 阅读设置缓存
	GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error)
	SetReadingSettings(ctx context.Context, userID string, settings *reader.ReadingSettings, expiration time.Duration) error
	InvalidateReadingSettings(ctx context.Context, userID string) error

	// 阅读进度缓存
	GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error)
	SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader.ReadingProgress, expiration time.Duration) error
	InvalidateReadingProgress(ctx context.Context, userID, bookID string) error

	// 批量清理
	InvalidateBookChapters(ctx context.Context, bookID string) error
	InvalidateUserData(ctx context.Context, userID string) error
}

// RedisReaderCacheService Redis缓存服务实现
type RedisReaderCacheService struct {
	client *redis.Client
	prefix string
}

// NewRedisReaderCacheService 创建Redis阅读器缓存服务
func NewRedisReaderCacheService(client *redis.Client, prefix string) ReaderCacheService {
	if prefix == "" {
		prefix = "qingyu"
	}
	return &RedisReaderCacheService{
		client: client,
		prefix: prefix,
	}
}

// 缓存键生成
func (c *RedisReaderCacheService) getKey(key string) string {
	return fmt.Sprintf("%s:reader:%s", c.prefix, key)
}

// =========================
// 章节内容缓存
// =========================

// GetChapterContent 获取章节内容缓存
func (c *RedisReaderCacheService) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	key := c.getKey(fmt.Sprintf("chapter_content:%s", chapterID))
	content, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // 缓存不存在
		}
		return "", fmt.Errorf("failed to get chapter content cache: %w", err)
	}
	return content, nil
}

// SetChapterContent 设置章节内容缓存
func (c *RedisReaderCacheService) SetChapterContent(ctx context.Context, chapterID string, content string, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("chapter_content:%s", chapterID))
	return c.client.Set(ctx, key, content, expiration).Err()
}

// InvalidateChapterContent 清除章节内容缓存
func (c *RedisReaderCacheService) InvalidateChapterContent(ctx context.Context, chapterID string) error {
	key := c.getKey(fmt.Sprintf("chapter_content:%s", chapterID))
	return c.client.Del(ctx, key).Err()
}

// =========================
// 章节信息缓存
// =========================

// GetChapter 获取章节信息缓存
func (c *RedisReaderCacheService) GetChapter(ctx context.Context, chapterID string) (*reader.Chapter, error) {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get chapter cache: %w", err)
	}

	var chapter reader.Chapter
	if err := json.Unmarshal([]byte(data), &chapter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chapter data: %w", err)
	}

	return &chapter, nil
}

// SetChapter 设置章节信息缓存
func (c *RedisReaderCacheService) SetChapter(ctx context.Context, chapterID string, chapter *reader.Chapter, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	jsonData, err := json.Marshal(chapter)
	if err != nil {
		return fmt.Errorf("failed to marshal chapter data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateChapter 清除章节信息缓存
func (c *RedisReaderCacheService) InvalidateChapter(ctx context.Context, chapterID string) error {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	return c.client.Del(ctx, key).Err()
}

// =========================
// 阅读设置缓存
// =========================

// GetReadingSettings 获取阅读设置缓存
func (c *RedisReaderCacheService) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	key := c.getKey(fmt.Sprintf("settings:%s", userID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get reading settings cache: %w", err)
	}

	var settings reader.ReadingSettings
	if err := json.Unmarshal([]byte(data), &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reading settings data: %w", err)
	}

	return &settings, nil
}

// SetReadingSettings 设置阅读设置缓存
func (c *RedisReaderCacheService) SetReadingSettings(ctx context.Context, userID string, settings *reader.ReadingSettings, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("settings:%s", userID))
	jsonData, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal reading settings data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateReadingSettings 清除阅读设置缓存
func (c *RedisReaderCacheService) InvalidateReadingSettings(ctx context.Context, userID string) error {
	key := c.getKey(fmt.Sprintf("settings:%s", userID))
	return c.client.Del(ctx, key).Err()
}

// =========================
// 阅读进度缓存
// =========================

// GetReadingProgress 获取阅读进度缓存
func (c *RedisReaderCacheService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	key := c.getKey(fmt.Sprintf("progress:%s:%s", userID, bookID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get reading progress cache: %w", err)
	}

	var progress reader.ReadingProgress
	if err := json.Unmarshal([]byte(data), &progress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reading progress data: %w", err)
	}

	return &progress, nil
}

// SetReadingProgress 设置阅读进度缓存
func (c *RedisReaderCacheService) SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader.ReadingProgress, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("progress:%s:%s", userID, bookID))
	jsonData, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal reading progress data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateReadingProgress 清除阅读进度缓存
func (c *RedisReaderCacheService) InvalidateReadingProgress(ctx context.Context, userID, bookID string) error {
	key := c.getKey(fmt.Sprintf("progress:%s:%s", userID, bookID))
	return c.client.Del(ctx, key).Err()
}

// =========================
// 批量清理
// =========================

// InvalidateBookChapters 清除书籍所有章节缓存
func (c *RedisReaderCacheService) InvalidateBookChapters(ctx context.Context, bookID string) error {
	pattern := c.getKey(fmt.Sprintf("chapter*:%s:*", bookID))
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get chapter cache keys: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// InvalidateUserData 清除用户所有缓存数据
func (c *RedisReaderCacheService) InvalidateUserData(ctx context.Context, userID string) error {
	// 清除设置缓存
	settingsKey := c.getKey(fmt.Sprintf("settings:%s", userID))
	_ = c.client.Del(ctx, settingsKey)

	// 清除进度缓存
	progressPattern := c.getKey(fmt.Sprintf("progress:%s:*", userID))
	progressKeys, err := c.client.Keys(ctx, progressPattern).Result()
	if err == nil && len(progressKeys) > 0 {
		_ = c.client.Del(ctx, progressKeys...)
	}

	return nil
}
