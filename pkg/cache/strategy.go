package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// CachePolicy 缓存策略
type CachePolicy struct {
	TTL           time.Duration // 基础过期时间
	RandomTTL     time.Duration // 随机过期时间范围（防止雪崩）
	Strategy      StrategyType  // 缓存策略类型
	Compress      bool          // 是否压缩
	SerializeFunc func(interface{}) ([]byte, error)
}

// StrategyType 缓存策略类型
type StrategyType string

const (
	StrategyCacheAside   StrategyType = "cache_aside"   // 旁路缓存
	StrategyWriteThrough StrategyType = "write_through" // 直写缓存
	StrategyWriteBehind  StrategyType = "write_behind"  // 异步写回
)

// CacheStrategy 缓存策略管理器
type CacheStrategy struct {
	client     RedisClient
	policies   map[string]*CachePolicy
	localCache LocalCache // 本地缓存（多级缓存）
}

// LocalCache 本地缓存接口
type LocalCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

// NewCacheStrategy 创建缓存策略管理器
func NewCacheStrategy(client RedisClient) *CacheStrategy {
	strategy := &CacheStrategy{
		client:   client,
		policies: make(map[string]*CachePolicy),
	}

	// 初始化默认策略
	strategy.initDefaultPolicies()

	return strategy
}

// initDefaultPolicies 初始化默认缓存策略
func (s *CacheStrategy) initDefaultPolicies() {
	// 用户会话 - 长期缓存
	s.RegisterPolicy("user:session", &CachePolicy{
		TTL:       7 * 24 * time.Hour,
		RandomTTL: 1 * time.Hour,
		Strategy:  StrategyCacheAside,
	})

	// 用户信息 - 中期缓存
	s.RegisterPolicy("user:info", &CachePolicy{
		TTL:       30 * time.Minute,
		RandomTTL: 5 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 书籍详情 - 热点数据
	s.RegisterPolicy("book:detail", &CachePolicy{
		TTL:       1 * time.Hour,
		RandomTTL: 10 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 章节内容 - 较长缓存
	s.RegisterPolicy("book:chapter", &CachePolicy{
		TTL:       6 * time.Hour,
		RandomTTL: 30 * time.Minute,
		Strategy:  StrategyCacheAside,
		Compress:  true, // 章节内容较大，启用压缩
	})

	// 书籍目录 - 长期缓存
	s.RegisterPolicy("book:catalog", &CachePolicy{
		TTL:       6 * time.Hour,
		RandomTTL: 30 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 热门列表 - 短期缓存
	s.RegisterPolicy("book:hot", &CachePolicy{
		TTL:       15 * time.Minute,
		RandomTTL: 3 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 推荐列表 - 短期缓存
	s.RegisterPolicy("book:recommend", &CachePolicy{
		TTL:       15 * time.Minute,
		RandomTTL: 3 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 搜索结果 - 短期缓存
	s.RegisterPolicy("search:result", &CachePolicy{
		TTL:       10 * time.Minute,
		RandomTTL: 2 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 统计数据 - 短期缓存
	s.RegisterPolicy("stats", &CachePolicy{
		TTL:       5 * time.Minute,
		RandomTTL: 1 * time.Minute,
		Strategy:  StrategyWriteBehind,
	})

	// 阅读进度 - 实时更新，但保留短期缓存
	s.RegisterPolicy("reading:progress", &CachePolicy{
		TTL:       5 * time.Minute,
		RandomTTL: 1 * time.Minute,
		Strategy:  StrategyWriteThrough,
	})

	// 评论列表 - 中期缓存
	s.RegisterPolicy("comment:list", &CachePolicy{
		TTL:       30 * time.Minute,
		RandomTTL: 5 * time.Minute,
		Strategy:  StrategyCacheAside,
	})

	// 配置信息 - 长期缓存
	s.RegisterPolicy("config", &CachePolicy{
		TTL:       1 * time.Hour,
		RandomTTL: 10 * time.Minute,
		Strategy:  StrategyCacheAside,
	})
}

// RegisterPolicy 注册缓存策略
func (s *CacheStrategy) RegisterPolicy(keyPrefix string, policy *CachePolicy) {
	s.policies[keyPrefix] = policy
}

// GetPolicy 获取缓存策略
func (s *CacheStrategy) GetPolicy(key string) (*CachePolicy, bool) {
	// 提取key前缀
	for prefix, policy := range s.policies {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return policy, true
		}
	}

	// 默认策略
	return &CachePolicy{
		TTL:       5 * time.Minute,
		RandomTTL: 1 * time.Minute,
		Strategy:  StrategyCacheAside,
	}, false
}

// Get 获取缓存（使用策略）
func (s *CacheStrategy) Get(ctx context.Context, key string, target interface{}) error {
	// 1. 检查本地缓存（多级缓存）
	if s.localCache != nil {
		if val, ok := s.localCache.Get(key); ok {
			// 从本地缓存返回
			return json.Unmarshal(val.([]byte), target)
		}
	}

	// 2. 从Redis获取
	val, err := s.client.Get(ctx, key)
	if err != nil {
		return err
	}

	// 3. 反序列化
	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		return fmt.Errorf("反序列化失败: %w", err)
	}

	// 4. 写入本地缓存
	if s.localCache != nil {
		s.localCache.Set(key, []byte(val), 5*time.Minute)
	}

	return nil
}

// Set 设置缓存（使用策略）
func (s *CacheStrategy) Set(ctx context.Context, key string, value interface{}) error {
	// 1. 获取策略
	policy, _ := s.GetPolicy(key)

	// 2. 序列化
	var data []byte
	var err error

	if policy.SerializeFunc != nil {
		data, err = policy.SerializeFunc(value)
	} else {
		data, err = json.Marshal(value)
	}

	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	// 3. 压缩（如果需要）
	if policy.Compress {
		data, err = compressData(data)
		if err != nil {
			return fmt.Errorf("压缩失败: %w", err)
		}
	}

	// 4. 计算过期时间（加入随机值防止雪崩）
	ttl := s.calculateTTL(policy.TTL, policy.RandomTTL)

	// 5. 写入Redis
	err = s.client.Set(ctx, key, data, ttl)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %w", err)
	}

	// 6. 写入本地缓存
	if s.localCache != nil {
		s.localCache.Set(key, data, 5*time.Minute)
	}

	return nil
}

// Delete 删除缓存
func (s *CacheStrategy) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// 删除Redis缓存
	err := s.client.Delete(ctx, keys...)
	if err != nil {
		return fmt.Errorf("删除Redis缓存失败: %w", err)
	}

	// 删除本地缓存
	if s.localCache != nil {
		for _, key := range keys {
			s.localCache.Delete(key)
		}
	}

	return nil
}

// GetOrLoad 获取缓存，如果不存在则加载
func (s *CacheStrategy) GetOrLoad(ctx context.Context, key string, target interface{}, loader func() (interface{}, error)) error {
	// 1. 尝试从缓存获取
	err := s.Get(ctx, key, target)
	if err == nil {
		return nil
	}

	// 2. 缓存不存在，加载数据
	data, err := loader()
	if err != nil {
		return fmt.Errorf("加载数据失败: %w", err)
	}

	// 3. 写入缓存
	err = s.Set(ctx, key, data)
	if err != nil {
		// 写入缓存失败不影响返回数据
		return nil
	}

	// 4. 转换为目标类型
	dataBytes, _ := json.Marshal(data)
	err = json.Unmarshal(dataBytes, target)
	if err != nil {
		return fmt.Errorf("数据转换失败: %w", err)
	}

	return nil
}

// MGet 批量获取（使用Pipeline优化）
func (s *CacheStrategy) MGet(ctx context.Context, keys []string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return make(map[string][]byte), nil
	}

	// 使用MGET批量获取
	vals, err := s.client.MGet(ctx, keys...)
	if err != nil {
		return nil, fmt.Errorf("批量获取失败: %w", err)
	}

	// 组装结果
	result := make(map[string][]byte)
	for i, key := range keys {
		if vals[i] != nil {
			result[key] = []byte(vals[i].(string))
		}
	}

	return result, nil
}

// MSet 批量设置（使用Pipeline优化）
func (s *CacheStrategy) MSet(ctx context.Context, items map[string]interface{}) error {
	if len(items) == 0 {
		return nil
	}

	// 准备批量数据
	pairs := make([]interface{}, 0, len(items)*2)
	for key, value := range items {
		// 序列化
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化失败: %w", err)
		}

		// 获取策略
		policy, _ := s.GetPolicy(key)

		// 压缩
		if policy.Compress {
			data, err = compressData(data)
			if err != nil {
				return fmt.Errorf("压缩失败: %w", err)
			}
		}

		pairs = append(pairs, key, data)
	}

	// 批量设置
	err := s.client.MSet(ctx, pairs...)
	if err != nil {
		return fmt.Errorf("批量设置失败: %w", err)
	}

	return nil
}

// calculateTTL 计算带随机值的过期时间
func (s *CacheStrategy) calculateTTL(baseTTL, randomTTL time.Duration) time.Duration {
	if randomTTL == 0 {
		return baseTTL
	}

	// 添加随机值（确保为正数）
	randomSeconds := rand.Int63() % int64(randomTTL)
	if randomSeconds < 0 {
		randomSeconds = -randomSeconds
	}
	return baseTTL + time.Duration(randomSeconds)*time.Second
}

// compressData 压缩数据
func compressData(data []byte) ([]byte, error) {
	// TODO: 实现压缩逻辑
	// 可以使用 gzip 或 zstd
	// 这里先简单返回原始数据
	return data, nil
}

// ============ 缓存预热功能 ============

// WarmUpCache 预热缓存
func (s *CacheStrategy) WarmUpCache(ctx context.Context, loader func() (map[string]interface{}, error)) error {
	// 1. 加载预热数据
	data, err := loader()
	if err != nil {
		return fmt.Errorf("加载预热数据失败: %w", err)
	}

	// 2. 批量写入缓存
	err = s.MSet(ctx, data)
	if err != nil {
		return fmt.Errorf("写入预热缓存失败: %w", err)
	}

	return nil
}

// ============ 缓存统计功能 ============

// CacheStats 缓存统计
type CacheStats struct {
	Hits   int64 `json:"hits"`
	Misses int64 `json:"misses"`
	Sets   int64 `json:"sets"`
	Deletes int64 `json:"deletes"`
}

func (s *CacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}

// ============ 缓存Key生成工具 ============

// BuildKey 构建缓存Key
func BuildKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}

// BuildUserKey 构建用户相关Key
func BuildUserKey(userID, keyType string) string {
	return BuildKey("user", keyType, userID)
}

// BuildBookKey 构建书籍相关Key
func BuildBookKey(bookID, keyType string) string {
	return BuildKey("book", keyType, bookID)
}

// BuildChapterKey 构建章节相关Key
func BuildChapterKey(bookID string, chapterNumber int) string {
	return BuildKey("book", "chapter", bookID, fmt.Sprintf("%d", chapterNumber))
}

// BuildSearchKey 构建搜索结果Key
func BuildSearchKey(query string, page int) string {
	return BuildKey("search", "result", query, fmt.Sprintf("page:%d", page))
}
