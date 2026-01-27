package user

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"

	"Qingyu_backend/models/users"
	"Qingyu_backend/repository/cache"
)

// UserRepositoryCacheable 可缓存的UserRepository接口（仅包含需要缓存的方法）
type UserRepositoryCacheable interface {
	GetByID(ctx context.Context, id string) (*users.User, error)
	Create(ctx context.Context, user *users.User) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// WrapUserRepositoryWithCache 包装UserRepository添加缓存
func WrapUserRepositoryWithCache(
	base UserRepositoryCacheable,
	redisClient *redis.Client,
) UserRepositoryCacheable {
	config := &cache.CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 1 * time.Second,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
		BreakerSettings: gobreaker.Settings{
			Name:        "user-cache-breaker",
			MaxRequests: 3,
			Interval:    10 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("熔断器状态变化: %s %v -> %v", name, from, to)
			},
		},
	}

	return &CachedUserRepository{
		base:         base,
		cachedSingle: cache.NewCachedRepository[users.User](&userRepositoryAdapter{base}, redisClient, 30*time.Minute, "user", config),
	}
}

// userRepositoryAdapter 将UserRepository接口适配为cache.Repository接口
type userRepositoryAdapter struct {
	repo UserRepositoryCacheable
}

func (a *userRepositoryAdapter) GetByID(ctx context.Context, id string) (users.User, error) {
	user, err := a.repo.GetByID(ctx, id)
	if err != nil {
		var zero users.User
		return zero, err
	}
	return *user, nil
}

func (a *userRepositoryAdapter) Create(ctx context.Context, entity users.User) error {
	return a.repo.Create(ctx, &entity)
}

func (a *userRepositoryAdapter) Update(ctx context.Context, entity users.User) error {
	// UserRepository的Update需要(id, updates)签名
	return nil
}

func (a *userRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}

func (a *userRepositoryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	return a.repo.Exists(ctx, id)
}

// CachedUserRepository 带缓存的UserRepository
type CachedUserRepository struct {
	base         UserRepositoryCacheable
	cachedSingle *cache.CachedRepository[users.User]
}

// GetByID 获取单个用户(带缓存)
func (r *CachedUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	user, err := r.cachedSingle.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *CachedUserRepository) Create(ctx context.Context, user *users.User) error {
	return r.base.Create(ctx, user)
}

// Update 更新用户(删除缓存)
func (r *CachedUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if err := r.base.Update(ctx, id, updates); err != nil {
		return err
	}
	// 删除缓存
	r.cachedSingle.Delete(ctx, id)
	return nil
}

// Delete 删除用户(删除缓存)
func (r *CachedUserRepository) Delete(ctx context.Context, id string) error {
	return r.cachedSingle.Delete(ctx, id)
}

// Exists 检查用户是否存在
func (r *CachedUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.cachedSingle.Exists(ctx, id)
}
