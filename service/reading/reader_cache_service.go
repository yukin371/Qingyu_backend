package reading

import (
	reader2 "Qingyu_backend/models/reader"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ReaderCacheService 阅读器缓存服务接口
type ReaderCacheService interface {
	// 阅读设置缓存
	GetReadingSettings(ctx context.Context, userID string) (*reader2.ReadingSettings, error)
	SetReadingSettings(ctx context.Context, userID string, settings *reader2.ReadingSettings, expiration time.Duration) error
	InvalidateReadingSettings(ctx context.Context, userID string) error

	// 阅读进度缓存
	GetReadingProgress(ctx context.Context, userID, bookID string) (*reader2.ReadingProgress, error)
	SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader2.ReadingProgress, expiration time.Duration) error
	InvalidateReadingProgress(ctx context.Context, userID, bookID string) error

	// 标注缓存
	GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error)
	SetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string, annotations []*reader2.Annotation, expiration time.Duration) error
	InvalidateAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) error

	// 批量清理
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
// 阅读设置缓存
// =========================

// GetReadingSettings 获取阅读设置缓存
func (c *RedisReaderCacheService) GetReadingSettings(ctx context.Context, userID string) (*reader2.ReadingSettings, error) {
	key := c.getKey(fmt.Sprintf("settings:%s", userID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get reading settings cache: %w", err)
	}

	var settings reader2.ReadingSettings
	if err := json.Unmarshal([]byte(data), &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reading settings data: %w", err)
	}

	return &settings, nil
}

// SetReadingSettings 设置阅读设置缓存
func (c *RedisReaderCacheService) SetReadingSettings(ctx context.Context, userID string, settings *reader2.ReadingSettings, expiration time.Duration) error {
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
func (c *RedisReaderCacheService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader2.ReadingProgress, error) {
	key := c.getKey(fmt.Sprintf("progress:%s:%s", userID, bookID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get reading progress cache: %w", err)
	}

	var progress reader2.ReadingProgress
	if err := json.Unmarshal([]byte(data), &progress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reading progress data: %w", err)
	}

	return &progress, nil
}

// SetReadingProgress 设置阅读进度缓存
func (c *RedisReaderCacheService) SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader2.ReadingProgress, expiration time.Duration) error {
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
// 标注缓存
// =========================

// GetAnnotationsByChapter 获取章节标注缓存
func (c *RedisReaderCacheService) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error) {
	key := c.getKey(fmt.Sprintf("annotations:%s:%s:%s", userID, bookID, chapterID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get annotations cache: %w", err)
	}

	var annotations []*reader2.Annotation
	if err := json.Unmarshal([]byte(data), &annotations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotations data: %w", err)
	}

	return annotations, nil
}

// SetAnnotationsByChapter 设置章节标注缓存
func (c *RedisReaderCacheService) SetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string, annotations []*reader2.Annotation, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("annotations:%s:%s:%s", userID, bookID, chapterID))
	jsonData, err := json.Marshal(annotations)
	if err != nil {
		return fmt.Errorf("failed to marshal annotations data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateAnnotationsByChapter 清除章节标注缓存
func (c *RedisReaderCacheService) InvalidateAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) error {
	key := c.getKey(fmt.Sprintf("annotations:%s:%s:%s", userID, bookID, chapterID))
	return c.client.Del(ctx, key).Err()
}

// =========================
// 批量清理
// =========================

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
