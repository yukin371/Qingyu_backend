package reading

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Qingyu_backend/models/reading/reader"

	"github.com/redis/go-redis/v9"
)

// AnnotationCacheService 注记缓存服务
type AnnotationCacheService struct {
	redisClient *redis.Client
	keyPrefix   string
	defaultTTL  time.Duration
}

// NewAnnotationCacheService 创建注记缓存服务
func NewAnnotationCacheService(redisClient *redis.Client) *AnnotationCacheService {
	return &AnnotationCacheService{
		redisClient: redisClient,
		keyPrefix:   "annotation:",
		defaultTTL:  30 * time.Minute, // 默认30分钟过期
	}
}

// GetAnnotation 获取单个注记缓存
func (s *AnnotationCacheService) GetAnnotation(ctx context.Context, annotationID string) (*reader.Annotation, error) {
	key := s.getAnnotationKey(annotationID)

	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("获取注记缓存失败: %w", err)
	}

	var annotation reader.Annotation
	if err := json.Unmarshal([]byte(data), &annotation); err != nil {
		return nil, fmt.Errorf("解析注记数据失败: %w", err)
	}

	return &annotation, nil
}

// SetAnnotation 设置单个注记缓存
func (s *AnnotationCacheService) SetAnnotation(ctx context.Context, annotation *reader.Annotation) error {
	key := s.getAnnotationKey(annotation.ID)

	data, err := json.Marshal(annotation)
	if err != nil {
		return fmt.Errorf("序列化注记数据失败: %w", err)
	}

	return s.redisClient.Set(ctx, key, data, s.defaultTTL).Err()
}

// DeleteAnnotation 删除注记缓存
func (s *AnnotationCacheService) DeleteAnnotation(ctx context.Context, annotationID string) error {
	key := s.getAnnotationKey(annotationID)
	return s.redisClient.Del(ctx, key).Err()
}

// GetChapterAnnotations 获取章节注记列表缓存
func (s *AnnotationCacheService) GetChapterAnnotations(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	key := s.getChapterAnnotationsKey(userID, bookID, chapterID)

	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("获取章节注记缓存失败: %w", err)
	}

	var annotations []*reader.Annotation
	if err := json.Unmarshal([]byte(data), &annotations); err != nil {
		return nil, fmt.Errorf("解析章节注记数据失败: %w", err)
	}

	return annotations, nil
}

// SetChapterAnnotations 设置章节注记列表缓存
func (s *AnnotationCacheService) SetChapterAnnotations(ctx context.Context, userID, bookID, chapterID string, annotations []*reader.Annotation) error {
	key := s.getChapterAnnotationsKey(userID, bookID, chapterID)

	data, err := json.Marshal(annotations)
	if err != nil {
		return fmt.Errorf("序列化章节注记数据失败: %w", err)
	}

	return s.redisClient.Set(ctx, key, data, s.defaultTTL).Err()
}

// InvalidateChapterAnnotations 使章节注记缓存失效
func (s *AnnotationCacheService) InvalidateChapterAnnotations(ctx context.Context, userID, bookID, chapterID string) error {
	key := s.getChapterAnnotationsKey(userID, bookID, chapterID)
	return s.redisClient.Del(ctx, key).Err()
}

// GetBookAnnotations 获取书籍注记列表缓存
func (s *AnnotationCacheService) GetBookAnnotations(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	key := s.getBookAnnotationsKey(userID, bookID)

	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("获取书籍注记缓存失败: %w", err)
	}

	var annotations []*reader.Annotation
	if err := json.Unmarshal([]byte(data), &annotations); err != nil {
		return nil, fmt.Errorf("解析书籍注记数据失败: %w", err)
	}

	return annotations, nil
}

// SetBookAnnotations 设置书籍注记列表缓存
func (s *AnnotationCacheService) SetBookAnnotations(ctx context.Context, userID, bookID string, annotations []*reader.Annotation) error {
	key := s.getBookAnnotationsKey(userID, bookID)

	data, err := json.Marshal(annotations)
	if err != nil {
		return fmt.Errorf("序列化书籍注记数据失败: %w", err)
	}

	return s.redisClient.Set(ctx, key, data, s.defaultTTL).Err()
}

// InvalidateBookAnnotations 使书籍注记缓存失效
func (s *AnnotationCacheService) InvalidateBookAnnotations(ctx context.Context, userID, bookID string) error {
	key := s.getBookAnnotationsKey(userID, bookID)
	return s.redisClient.Del(ctx, key).Err()
}

// InvalidateUserAnnotations 使用户所有注记缓存失效
func (s *AnnotationCacheService) InvalidateUserAnnotations(ctx context.Context, userID string) error {
	pattern := s.getUserAnnotationsPattern(userID)

	// 使用SCAN命令查找所有匹配的键
	iter := s.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("删除缓存失败: %w", err)
		}
	}

	return iter.Err()
}

// BatchGetAnnotations 批量获取注记
func (s *AnnotationCacheService) BatchGetAnnotations(ctx context.Context, annotationIDs []string) ([]*reader.Annotation, error) {
	if len(annotationIDs) == 0 {
		return []*reader.Annotation{}, nil
	}

	keys := make([]string, len(annotationIDs))
	for i, id := range annotationIDs {
		keys[i] = s.getAnnotationKey(id)
	}

	// 使用MGET批量获取
	results, err := s.redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("批量获取注记缓存失败: %w", err)
	}

	annotations := make([]*reader.Annotation, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue // 跳过不存在的键
		}

		data, ok := result.(string)
		if !ok {
			continue
		}

		var annotation reader.Annotation
		if err := json.Unmarshal([]byte(data), &annotation); err != nil {
			continue // 跳过解析失败的数据
		}

		annotations = append(annotations, &annotation)
	}

	return annotations, nil
}

// BatchSetAnnotations 批量设置注记缓存
func (s *AnnotationCacheService) BatchSetAnnotations(ctx context.Context, annotations []*reader.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	// 使用Pipeline批量设置
	pipe := s.redisClient.Pipeline()

	for _, annotation := range annotations {
		key := s.getAnnotationKey(annotation.ID)
		data, err := json.Marshal(annotation)
		if err != nil {
			return fmt.Errorf("序列化注记数据失败: %w", err)
		}
		pipe.Set(ctx, key, data, s.defaultTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetAnnotationStats 获取注记统计缓存
func (s *AnnotationCacheService) GetAnnotationStats(ctx context.Context, userID, bookID string) (*AnnotationStats, error) {
	key := s.getAnnotationStatsKey(userID, bookID)

	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("获取注记统计缓存失败: %w", err)
	}

	var stats AnnotationStats
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, fmt.Errorf("解析注记统计数据失败: %w", err)
	}

	return &stats, nil
}

// SetAnnotationStats 设置注记统计缓存
func (s *AnnotationCacheService) SetAnnotationStats(ctx context.Context, userID, bookID string, stats *AnnotationStats) error {
	key := s.getAnnotationStatsKey(userID, bookID)

	data, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("序列化注记统计数据失败: %w", err)
	}

	return s.redisClient.Set(ctx, key, data, s.defaultTTL).Err()
}

// AnnotationStats 注记统计
type AnnotationStats struct {
	TotalCount     int       `json:"totalCount"`
	BookmarkCount  int       `json:"bookmarkCount"`
	HighlightCount int       `json:"highlightCount"`
	NoteCount      int       `json:"noteCount"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

// 键生成辅助方法
func (s *AnnotationCacheService) getAnnotationKey(annotationID string) string {
	return fmt.Sprintf("%s%s", s.keyPrefix, annotationID)
}

func (s *AnnotationCacheService) getChapterAnnotationsKey(userID, bookID, chapterID string) string {
	return fmt.Sprintf("%suser:%s:book:%s:chapter:%s:list", s.keyPrefix, userID, bookID, chapterID)
}

func (s *AnnotationCacheService) getBookAnnotationsKey(userID, bookID string) string {
	return fmt.Sprintf("%suser:%s:book:%s:list", s.keyPrefix, userID, bookID)
}

func (s *AnnotationCacheService) getUserAnnotationsPattern(userID string) string {
	return fmt.Sprintf("%suser:%s:*", s.keyPrefix, userID)
}

func (s *AnnotationCacheService) getAnnotationStatsKey(userID, bookID string) string {
	return fmt.Sprintf("%suser:%s:book:%s:stats", s.keyPrefix, userID, bookID)
}
