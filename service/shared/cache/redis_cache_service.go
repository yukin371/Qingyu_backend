package cache

import (
	"context"
	"fmt"
	"time"

	serviceInterfaces "Qingyu_backend/service/interfaces"

	"github.com/redis/go-redis/v9"
)

// RedisCacheService Redis 缓存服务实现
type RedisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService 创建 Redis 缓存服务
func NewRedisCacheService(client *redis.Client) serviceInterfaces.CacheService {
	return &RedisCacheService{
		client: client,
	}
}

// ============ 基础操作 ============

// Get 获取缓存
func (s *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	result, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("获取缓存失败: %w", err)
	}
	return result, nil
}

// Set 设置缓存
func (s *RedisCacheService) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := s.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}
	return nil
}

// Delete 删除缓存
func (s *RedisCacheService) Delete(ctx context.Context, key string) error {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// Exists 检查键是否存在
func (s *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查键存在性失败: %w", err)
	}
	return count > 0, nil
}

// ============ 批量操作 ============

// MGet 批量获取
func (s *RedisCacheService) MGet(ctx context.Context, keys ...string) ([]string, error) {
	results, err := s.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("批量获取失败: %w", err)
	}

	// 转换为字符串切片
	values := make([]string, len(results))
	for i, result := range results {
		if result != nil {
			values[i] = result.(string)
		}
	}

	return values, nil
}

// MSet 批量设置
func (s *RedisCacheService) MSet(ctx context.Context, kvPairs map[string]string, expiration time.Duration) error {
	// Redis MSET 不支持过期时间，需要使用 Pipeline
	pipe := s.client.Pipeline()

	for key, value := range kvPairs {
		pipe.Set(ctx, key, value, expiration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("批量设置失败: %w", err)
	}

	return nil
}

// MDelete 批量删除
func (s *RedisCacheService) MDelete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := s.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}

	return nil
}

// ============ 高级操作 ============

// Expire 设置过期时间
func (s *RedisCacheService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := s.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置过期时间失败: %w", err)
	}
	return nil
}

// TTL 获取剩余过期时间
func (s *RedisCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取TTL失败: %w", err)
	}
	return ttl, nil
}

// Increment 递增
func (s *RedisCacheService) Increment(ctx context.Context, key string) (int64, error) {
	result, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("递增失败: %w", err)
	}
	return result, nil
}

// Decrement 递减
func (s *RedisCacheService) Decrement(ctx context.Context, key string) (int64, error) {
	result, err := s.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("递减失败: %w", err)
	}
	return result, nil
}

// ============ 哈希操作 ============

// HGet 获取哈希字段
func (s *RedisCacheService) HGet(ctx context.Context, key, field string) (string, error) {
	result, err := s.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("field not found: %s", field)
	}
	if err != nil {
		return "", fmt.Errorf("获取哈希字段失败: %w", err)
	}
	return result, nil
}

// HSet 设置哈希字段
func (s *RedisCacheService) HSet(ctx context.Context, key, field, value string) error {
	err := s.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		return fmt.Errorf("设置哈希字段失败: %w", err)
	}
	return nil
}

// HGetAll 获取所有哈希字段
func (s *RedisCacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取所有哈希字段失败: %w", err)
	}
	return result, nil
}

// HDelete 删除哈希字段
func (s *RedisCacheService) HDelete(ctx context.Context, key string, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}

	err := s.client.HDel(ctx, key, fields...).Err()
	if err != nil {
		return fmt.Errorf("删除哈希字段失败: %w", err)
	}
	return nil
}

// ============ 集合操作 ============

// SAdd 添加集合成员
func (s *RedisCacheService) SAdd(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	// 转换为 interface{} 切片
	interfaceMembers := make([]interface{}, len(members))
	for i, member := range members {
		interfaceMembers[i] = member
	}

	err := s.client.SAdd(ctx, key, interfaceMembers...).Err()
	if err != nil {
		return fmt.Errorf("添加集合成员失败: %w", err)
	}
	return nil
}

// SMembers 获取集合所有成员
func (s *RedisCacheService) SMembers(ctx context.Context, key string) ([]string, error) {
	result, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取集合成员失败: %w", err)
	}
	return result, nil
}

// SIsMember 检查是否为集合成员
func (s *RedisCacheService) SIsMember(ctx context.Context, key, member string) (bool, error) {
	result, err := s.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, fmt.Errorf("检查集合成员失败: %w", err)
	}
	return result, nil
}

// SRemove 移除集合成员
func (s *RedisCacheService) SRemove(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	// 转换为 interface{} 切片
	interfaceMembers := make([]interface{}, len(members))
	for i, member := range members {
		interfaceMembers[i] = member
	}

	err := s.client.SRem(ctx, key, interfaceMembers...).Err()
	if err != nil {
		return fmt.Errorf("移除集合成员失败: %w", err)
	}
	return nil
}

// ============ 有序集合操作 ============

// ZAdd 添加有序集合成员
func (s *RedisCacheService) ZAdd(ctx context.Context, key string, score float64, member string) error {
	err := s.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
	if err != nil {
		return fmt.Errorf("添加有序集合成员失败: %w", err)
	}
	return nil
}

// ZRange 获取有序集合范围
func (s *RedisCacheService) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	result, err := s.client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("获取有序集合范围失败: %w", err)
	}
	return result, nil
}

// ZRangeWithScores 获取有序集合范围及分数
func (s *RedisCacheService) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (map[string]float64, error) {
	result, err := s.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("获取有序集合范围失败: %w", err)
	}

	// 转换为 map
	scores := make(map[string]float64)
	for _, z := range result {
		scores[z.Member.(string)] = z.Score
	}

	return scores, nil
}

// ZRemove 移除有序集合成员
func (s *RedisCacheService) ZRemove(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	// 转换为 interface{} 切片
	interfaceMembers := make([]interface{}, len(members))
	for i, member := range members {
		interfaceMembers[i] = member
	}

	err := s.client.ZRem(ctx, key, interfaceMembers...).Err()
	if err != nil {
		return fmt.Errorf("移除有序集合成员失败: %w", err)
	}
	return nil
}

// ============ 服务管理 ============

// Ping 健康检查
func (s *RedisCacheService) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// FlushDB 清空数据库
func (s *RedisCacheService) FlushDB(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}

// Close 关闭连接
func (s *RedisCacheService) Close() error {
	return s.client.Close()
}
