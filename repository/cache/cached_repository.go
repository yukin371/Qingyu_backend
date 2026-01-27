package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
)

// ErrNotFound 记录未找到错误
var ErrNotFound = errors.New("record not found")

// Cacheable 可缓存的数据接口
type Cacheable interface {
	GetID() string
}

// Repository 通用Repository接口
type Repository[T Cacheable] interface {
	GetByID(ctx context.Context, id string) (T, error)
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// CacheConfig 缓存配置结构体
type CacheConfig struct {
	Enabled           bool                // 总开关
	DoubleDeleteDelay time.Duration       // 双删策略延迟
	NullCacheTTL      time.Duration       // 空值缓存TTL
	NullCachePrefix   string              // 空值缓存前缀 (@@NULL@@)
	BreakerSettings   gobreaker.Settings  // 熔断器设置
}

// CachedRepository 缓存装饰器
type CachedRepository[T Cacheable] struct {
	base    Repository[T]
	client  *redis.Client
	ttl     time.Duration
	prefix  string
	enabled bool
	breaker *gobreaker.CircuitBreaker
	config  *CacheConfig
}

// NewCachedRepository 创建缓存装饰器
func NewCachedRepository[T Cacheable](
	base Repository[T],
	client *redis.Client,
	ttl time.Duration,
	prefix string,
	config *CacheConfig,
) *CachedRepository[T] {
	if config == nil {
		config = &CacheConfig{
			Enabled:           true,
			DoubleDeleteDelay: 1 * time.Second,
			NullCacheTTL:      30 * time.Second,
			NullCachePrefix:   "@@NULL@@",
		}
	}

	// 设置默认的熔断器配置
	if config.BreakerSettings.ReadyToTrip == nil {
		config.BreakerSettings = gobreaker.Settings{
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
		}
	}

	breaker := gobreaker.NewCircuitBreaker(config.BreakerSettings)

	return &CachedRepository[T]{
		base:    base,
		client:  client,
		ttl:     ttl,
		prefix:  prefix,
		enabled: config.Enabled,
		breaker: breaker,
		config:  config,
	}
}

// GetByID 带缓存的查询
func (r *CachedRepository[T]) GetByID(ctx context.Context, id string) (T, error) {
	if !r.enabled {
		return r.base.GetByID(ctx, id)
	}

	var result T
	_, err := r.breaker.Execute(func() (interface{}, error) {
		var e error
		result, e = r.getFromCacheOrDB(ctx, id)
		return result, e
	})

	if err != nil {
		// 熔断器触发，降级到直连DB
		return r.base.GetByID(ctx, id)
	}

	return result, nil
}

// getFromCacheOrDB 从缓存或数据库获取数据（核心逻辑）
func (r *CachedRepository[T]) getFromCacheOrDB(ctx context.Context, id string) (T, error) {
	key := r.cacheKey(id)
	cached, err := r.client.Get(ctx, key).Result()

	if err == nil {
		// 检查空值缓存
		if cached == r.config.NullCachePrefix {
			var zero T
			return zero, ErrNotFound
		}

		var entity T
		if err := json.Unmarshal([]byte(cached), &entity); err == nil {
			return entity, nil
		}
	}

	// 缓存错误不应影响业务
	if err != nil && err != redis.Nil {
		log.Printf("缓存读取失败(降级): %v", err)
	}

	// 查询数据库
	entity, err := r.base.GetByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			// 缓存空值防止穿透
			go func() {
				nullKey := r.cacheKey(id)
				r.client.Set(context.Background(), nullKey, r.config.NullCachePrefix, r.config.NullCacheTTL)
			}()
		}
		var zero T
		return zero, err
	}

	// 异步写入缓存
	go func() {
		data, _ := json.Marshal(entity)
		r.client.Set(context.Background(), key, data, r.ttl)
	}()

	return entity, nil
}

// cacheKey 生成缓存键
func (r *CachedRepository[T]) cacheKey(id string) string {
	return fmt.Sprintf("%s:%s", r.prefix, id)
}

// Create 创建实体（不缓存）
func (r *CachedRepository[T]) Create(ctx context.Context, entity T) error {
	if err := r.base.Create(ctx, entity); err != nil {
		return err
	}

	// 缓存由异步预热填充
	return nil
}

// Update 更新实体（删除缓存）
func (r *CachedRepository[T]) Update(ctx context.Context, entity T) error {
	if err := r.base.Update(ctx, entity); err != nil {
		return err
	}

	key := r.cacheKey(entity.GetID())
	r.client.Del(ctx, key)

	// 双删策略：延迟后再次删除
	go func() {
		time.Sleep(r.config.DoubleDeleteDelay)
		r.client.Del(context.Background(), key)
	}()

	return nil
}

// Delete 删除实体（删除缓存）
func (r *CachedRepository[T]) Delete(ctx context.Context, id string) error {
	if err := r.base.Delete(ctx, id); err != nil {
		return err
	}

	key := r.cacheKey(id)
	r.client.Del(ctx, key)

	// 双删策略
	go func() {
		time.Sleep(r.config.DoubleDeleteDelay)
		r.client.Del(context.Background(), key)
	}()

	return nil
}

// Exists 检查实体是否存在
func (r *CachedRepository[T]) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
