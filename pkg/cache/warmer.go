package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
)

// HotBookRepository 热门书籍Repository接口（只包含需要的方法）
type HotBookRepository interface {
	GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
}

// ActiveUserRepository 活跃用户Repository接口
type ActiveUserRepository interface {
	GetActiveUsers(ctx context.Context, limit int64) ([]*users.User, error)
}

// CacheWarmer 缓存预热器
type CacheWarmer struct {
	bookRepo HotBookRepository
	userRepo ActiveUserRepository
	redis    *redis.Client
}

// NewCacheWarmer 创建缓存预热器
func NewCacheWarmer(
	bookRepo HotBookRepository,
	userRepo ActiveUserRepository,
	redis *redis.Client,
) *CacheWarmer {
	return &CacheWarmer{
		bookRepo: bookRepo,
		userRepo: userRepo,
		redis:    redis,
	}
}

// WarmUpCache 预热热点数据
func (w *CacheWarmer) WarmUpCache(ctx context.Context) error {
	log.Println("=== 开始缓存预热 ===")
	startTime := time.Now()

	// 预热热门书籍
	if err := w.warmUpHotBooks(ctx); err != nil {
		log.Printf("⚠️  预热热门书籍失败: %v", err)
		// 预热失败不影响应用启动
	} else {
		log.Printf("✅ 预热热门书籍成功")
	}

	// 预热活跃用户（可选）
	if err := w.warmUpActiveUsers(ctx); err != nil {
		log.Printf("⚠️  预热活跃用户失败: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("=== 缓存预热完成，耗时: %v ===", duration)

	return nil
}

// warmUpHotBooks 预热热门书籍
func (w *CacheWarmer) warmUpHotBooks(ctx context.Context) error {
	log.Println("正在预热热门书籍...")

	// 获取热门书籍（按view_count排序）
	hotBooks, err := w.bookRepo.GetHotBooks(ctx, 100, 0)
	if err != nil {
		return fmt.Errorf("获取热门书籍失败: %w", err)
	}

	// 写入缓存
	successCount := 0
	for _, book := range hotBooks {
		if book == nil {
			continue
		}

		key := fmt.Sprintf("book:%s", book.ID.Hex())
		data, err := json.Marshal(book)
		if err != nil {
			log.Printf("序列化书籍失败 [%s]: %v", book.ID.Hex(), err)
			continue
		}

		// 缓存1小时
		if err := w.redis.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
			log.Printf("缓存书籍失败 [%s]: %v", book.ID.Hex(), err)
			continue
		}

		successCount++
	}

	log.Printf("预热热门书籍: %d/%d", successCount, len(hotBooks))
	return nil
}

// warmUpActiveUsers 预热活跃用户
func (w *CacheWarmer) warmUpActiveUsers(ctx context.Context) error {
	log.Println("正在预热活跃用户...")

	// 获取活跃用户（按last_login_at排序）
	activeUsers, err := w.userRepo.GetActiveUsers(ctx, 50)
	if err != nil {
		return fmt.Errorf("获取活跃用户失败: %w", err)
	}

	// 写入缓存
	successCount := 0
	for _, user := range activeUsers {
		if user == nil {
			continue
		}

		key := fmt.Sprintf("user:%s", user.ID.Hex())
		data, err := json.Marshal(user)
		if err != nil {
			log.Printf("序列化用户失败 [%s]: %v", user.ID.Hex(), err)
			continue
		}

		// 缓存30分钟
		if err := w.redis.Set(ctx, key, data, 30*time.Minute).Err(); err != nil {
			log.Printf("缓存用户失败 [%s]: %v", user.ID.Hex(), err)
			continue
		}

		successCount++
	}

	log.Printf("预热活跃用户: %d/%d", successCount, len(activeUsers))
	return nil
}
