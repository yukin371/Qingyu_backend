package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/cache"
)

// BookRepositoryCacheable 可缓存的BookRepository接口（仅包含需要缓存的方法）
type BookRepositoryCacheable interface {
	GetByID(ctx context.Context, id string) (*bookstore.Book, error)
	Create(ctx context.Context, book *bookstore.Book) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// WrapBookRepositoryWithCache 包装BookRepository添加缓存
func WrapBookRepositoryWithCache(
	base BookRepositoryCacheable,
	redisClient *redis.Client,
) BookRepositoryCacheable {
	config := &cache.CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 1 * time.Second,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",
		BreakerSettings: gobreaker.Settings{
			Name:        "book-cache-breaker",
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

	return &CachedBookRepository{
		base:         base,
		cachedSingle: cache.NewCachedRepository[bookstore.Book](&bookRepositoryAdapter{base}, redisClient, 1*time.Hour, "book", config),
	}
}

// bookRepositoryAdapter 将BookRepository接口适配为cache.Repository接口
type bookRepositoryAdapter struct {
	repo BookRepositoryCacheable
}

func (a *bookRepositoryAdapter) GetByID(ctx context.Context, id string) (bookstore.Book, error) {
	book, err := a.repo.GetByID(ctx, id)
	if err != nil {
		var zero bookstore.Book
		return zero, err
	}
	return *book, nil
}

func (a *bookRepositoryAdapter) Create(ctx context.Context, entity bookstore.Book) error {
	return a.repo.Create(ctx, &entity)
}

func (a *bookRepositoryAdapter) Update(ctx context.Context, entity bookstore.Book) error {
	// BookRepository的Update需要(id, updates)签名
	// 这里我们只标记需要删除缓存
	return nil
}

func (a *bookRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}

func (a *bookRepositoryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	return a.repo.Exists(ctx, id)
}

// CachedBookRepository 带缓存的BookRepository
type CachedBookRepository struct {
	base         BookRepositoryCacheable
	cachedSingle *cache.CachedRepository[bookstore.Book]
}

// GetByID 获取单本书籍(带缓存)
func (r *CachedBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
	book, err := r.cachedSingle.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Create 创建书籍
func (r *CachedBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	return r.base.Create(ctx, book)
}

// Update 更新书籍(删除缓存)
func (r *CachedBookRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if err := r.base.Update(ctx, id, updates); err != nil {
		return err
	}
	// 删除缓存
	r.cachedSingle.Delete(ctx, id)
	return nil
}

// Delete 删除书籍(删除缓存)
func (r *CachedBookRepository) Delete(ctx context.Context, id string) error {
	return r.cachedSingle.Delete(ctx, id)
}

// Exists 检查书籍是否存在
func (r *CachedBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.cachedSingle.Exists(ctx, id)
}
