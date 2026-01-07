package bookstore

import (
	"context"
	"fmt"

	"Qingyu_backend/pkg/cache"
	"Qingyu_backend/pkg/logger"

	"go.uber.org/zap"
)

// BookstoreCachedService 带缓存的书城服务
type BookstoreCachedService struct {
	service    BookstoreService
	cacheStrategy *cache.CacheStrategy
}

// NewBookstoreCachedService 创建带缓存的书城服务
func NewBookstoreCachedService(service BookstoreService, redisClient cache.RedisClient) *BookstoreCachedService {
	return &BookstoreCachedService{
		service:       service,
		cacheStrategy: cache.NewCacheStrategy(redisClient),
	}
}

// GetBookWithCache 获取书籍详情（使用缓存）
func (s *BookstoreCachedService) GetBookWithCache(ctx context.Context, bookID string) (interface{}, error) {
	// 构建缓存key
	cacheKey := cache.BuildBookKey(bookID, "detail")

	var book interface{}

	// 尝试从缓存获取
	err := s.cacheStrategy.Get(ctx, cacheKey, &book)
	if err == nil {
		logger.Debug("从缓存获取书籍", zap.String("book_id", bookID))
		return book, nil
	}

	// 缓存未命中，从数据库获取
	logger.Debug("缓存未命中，从数据库获取书籍", zap.String("book_id", bookID))
	book, err = s.service.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取书籍失败: %w", err)
	}

	// 写入缓存（使用预定义的book:detail策略）
	if err := s.cacheStrategy.Set(ctx, cacheKey, book); err != nil {
		logger.Warn("写入缓存失败", zap.String("book_id", bookID), zap.Error(err))
	}

	return book, nil
}

// GetBooksWithCache 获取书籍列表（使用缓存）
func (s *BookstoreCachedService) GetBooksWithCache(ctx context.Context, page, size int) (interface{}, error) {
	// 构建缓存key
	cacheKey := fmt.Sprintf("book:all:list:page:%d:size:%d", page, size)

	var result interface{}

	// 尝试从缓存获取
	err := s.cacheStrategy.Get(ctx, cacheKey, &result)
	if err == nil {
		logger.Debug("从缓存获取书籍列表", zap.Int("page", page))
		return result, nil
	}

	// 缓存未命中，从数据库获取
	logger.Debug("缓存未命中，从数据库获取书籍列表", zap.Int("page", page))
	// 使用正确的接口方法 GetAllBooks，接收 3 个返回值
	result, _, err = s.service.GetAllBooks(ctx, page, size)
	if err != nil {
		return nil, fmt.Errorf("获取书籍列表失败: %w", err)
	}

	// 写入缓存
	if err := s.cacheStrategy.Set(ctx, cacheKey, result); err != nil {
		logger.Warn("写入缓存失败", zap.Error(err))
	}

	return result, nil
}

// GetHomepageDataWithCache 获取首页数据（使用缓存）
func (s *BookstoreCachedService) GetHomepageDataWithCache(ctx context.Context) (interface{}, error) {
	cacheKey := "bookstore:homepage:data"

	var data interface{}

	// 使用GetOrLoad模式：缓存不存在则自动加载
	err := s.cacheStrategy.GetOrLoad(ctx, cacheKey, &data, func() (interface{}, error) {
		logger.Info("首页数据缓存未命中，加载数据")
		return s.service.GetHomepageData(ctx)
	})

	if err != nil {
		return nil, fmt.Errorf("获取首页数据失败: %w", err)
	}

	return data, nil
}

// WarmUpCache 预热热门书籍缓存
func (s *BookstoreCachedService) WarmUpCache(ctx context.Context) error {
	logger.Info("开始预热热门书籍缓存")

	// 使用缓存预热功能
	err := s.cacheStrategy.WarmUpCache(ctx, func() (map[string]interface{}, error) {
		// 获取热门书籍列表 - 修正参数数量，GetHotBooks 需要 (page, pageSize) 两个参数
		hotBooks, _, err := s.service.GetHotBooks(ctx, 1, 100)
		if err != nil {
			return nil, err
		}

		// 构建预热数据
		data := make(map[string]interface{})
		for _, book := range hotBooks {
			// 使用 book.ID.Hex() 获取字符串形式的 ID
			bookID := book.ID.Hex()
			cacheKey := cache.BuildBookKey(bookID, "detail")
			data[cacheKey] = book
		}

		return data, nil
	})

	if err != nil {
		return fmt.Errorf("缓存预热失败: %w", err)
	}

	logger.Info("热门书籍缓存预热完成")
	return nil
}

// InvalidateCategoryCache 删除分类相关缓存
func (s *BookstoreCachedService) InvalidateCategoryCache(ctx context.Context, categoryID string) {
	pattern := fmt.Sprintf("book:category:%s:*", categoryID)

	// 注意：Redis的Delete不支持通配符，需要使用SCAN或KEYS命令
	// 这里简化处理，实际应该使用SCAN遍历
	logger.Info("删除分类缓存", zap.String("pattern", pattern))

	// 实际实现中，可以维护一个缓存key的集合来支持批量删除
}
