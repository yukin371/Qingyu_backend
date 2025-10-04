package bookstore

import (
	"context"
	"time"

	"Qingyu_backend/models/reading/bookstore"
)

// CachedBookstoreService 带缓存的书城服务
type CachedBookstoreService struct {
	service BookstoreService
	cache   CacheService
}

// NewCachedBookstoreService 创建带缓存的书城服务
func NewCachedBookstoreService(service BookstoreService, cache CacheService) BookstoreService {
	return &CachedBookstoreService{
		service: service,
		cache:   cache,
	}
}

// 缓存过期时间配置
const (
	HomepageCacheExpiration = 5 * time.Minute   // 首页缓存5分钟
	RankingCacheExpiration  = 10 * time.Minute  // 榜单缓存10分钟
	BannerCacheExpiration   = 30 * time.Minute  // Banner缓存30分钟
	BookCacheExpiration     = 1 * time.Hour     // 书籍缓存1小时
	CategoryCacheExpiration = 2 * time.Hour     // 分类缓存2小时
)

// GetHomepageData 获取首页数据（带缓存）
func (c *CachedBookstoreService) GetHomepageData(ctx context.Context) (*HomepageData, error) {
	// 尝试从缓存获取
	if data, err := c.cache.GetHomepageData(ctx); err == nil && data != nil {
		return data, nil
	}

	// 缓存未命中，从服务获取
	data, err := c.service.GetHomepageData(ctx)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetHomepageData(context.Background(), data, HomepageCacheExpiration)
	}()

	return data, nil
}

// GetBookByID 获取书籍详情（带缓存）
func (c *CachedBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
	// 尝试从缓存获取
	if book, err := c.cache.GetBook(ctx, id); err == nil && book != nil {
		return book, nil
	}

	// 缓存未命中，从服务获取
	book, err := c.service.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetBook(context.Background(), id, book, BookCacheExpiration)
	}()

	return book, nil
}

// GetActiveBanners 获取激活的Banner列表（带缓存）
func (c *CachedBookstoreService) GetActiveBanners(ctx context.Context, limit int) ([]*bookstore.Banner, error) {
	// 尝试从缓存获取
	if banners, err := c.cache.GetActiveBanners(ctx); err == nil && banners != nil {
		// 如果缓存的数量足够，直接返回
		if len(banners) >= limit {
			if limit > len(banners) {
				return banners, nil
			}
			return banners[:limit], nil
		}
	}

	// 缓存未命中或数量不足，从服务获取
	banners, err := c.service.GetActiveBanners(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetActiveBanners(context.Background(), banners, BannerCacheExpiration)
	}()

	return banners, nil
}

// GetCategoryTree 获取分类树（带缓存）
func (c *CachedBookstoreService) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	// 尝试从缓存获取
	if tree, err := c.cache.GetCategoryTree(ctx); err == nil && tree != nil {
		return tree, nil
	}

	// 缓存未命中，从服务获取
	tree, err := c.service.GetCategoryTree(ctx)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetCategoryTree(context.Background(), tree, CategoryCacheExpiration)
	}()

	return tree, nil
}

// GetRealtimeRanking 获取实时榜（带缓存）
func (c *CachedBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore.RankingItem, error) {
	period := bookstore.GetPeriodString(bookstore.RankingTypeRealtime, time.Now())
	
	// 尝试从缓存获取
	if items, err := c.cache.GetRanking(ctx, bookstore.RankingTypeRealtime, period); err == nil && items != nil {
		if len(items) >= limit {
			if limit > len(items) {
				return items, nil
			}
			return items[:limit], nil
		}
	}

	// 缓存未命中，从服务获取
	items, err := c.service.GetRealtimeRanking(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetRanking(context.Background(), bookstore.RankingTypeRealtime, period, items, RankingCacheExpiration)
	}()

	return items, nil
}

// GetWeeklyRanking 获取周榜（带缓存）
func (c *CachedBookstoreService) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	if period == "" {
		period = bookstore.GetPeriodString(bookstore.RankingTypeWeekly, time.Now())
	}

	// 尝试从缓存获取
	if items, err := c.cache.GetRanking(ctx, bookstore.RankingTypeWeekly, period); err == nil && items != nil {
		if len(items) >= limit {
			if limit > len(items) {
				return items, nil
			}
			return items[:limit], nil
		}
	}

	// 缓存未命中，从服务获取
	items, err := c.service.GetWeeklyRanking(ctx, period, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetRanking(context.Background(), bookstore.RankingTypeWeekly, period, items, RankingCacheExpiration)
	}()

	return items, nil
}

// GetMonthlyRanking 获取月榜（带缓存）
func (c *CachedBookstoreService) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	if period == "" {
		period = bookstore.GetPeriodString(bookstore.RankingTypeMonthly, time.Now())
	}

	// 尝试从缓存获取
	if items, err := c.cache.GetRanking(ctx, bookstore.RankingTypeMonthly, period); err == nil && items != nil {
		if len(items) >= limit {
			if limit > len(items) {
				return items, nil
			}
			return items[:limit], nil
		}
	}

	// 缓存未命中，从服务获取
	items, err := c.service.GetMonthlyRanking(ctx, period, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetRanking(context.Background(), bookstore.RankingTypeMonthly, period, items, RankingCacheExpiration)
	}()

	return items, nil
}

// GetNewbieRanking 获取新人榜（带缓存）
func (c *CachedBookstoreService) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	if period == "" {
		period = bookstore.GetPeriodString(bookstore.RankingTypeNewbie, time.Now())
	}

	// 尝试从缓存获取
	if items, err := c.cache.GetRanking(ctx, bookstore.RankingTypeNewbie, period); err == nil && items != nil {
		if len(items) >= limit {
			if limit > len(items) {
				return items, nil
			}
			return items[:limit], nil
		}
	}

	// 缓存未命中，从服务获取
	items, err := c.service.GetNewbieRanking(ctx, period, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetRanking(context.Background(), bookstore.RankingTypeNewbie, period, items, RankingCacheExpiration)
	}()

	return items, nil
}

// GetRankingByType 根据类型获取榜单（带缓存）
func (c *CachedBookstoreService) GetRankingByType(ctx context.Context, rankingType bookstore.RankingType, period string, limit int) ([]*bookstore.RankingItem, error) {
	if period == "" {
		period = bookstore.GetPeriodString(rankingType, time.Now())
	}

	// 尝试从缓存获取
	if items, err := c.cache.GetRanking(ctx, rankingType, period); err == nil && items != nil {
		if len(items) >= limit {
			if limit > len(items) {
				return items, nil
			}
			return items[:limit], nil
		}
	}

	// 缓存未命中，从服务获取
	items, err := c.service.GetRankingByType(ctx, rankingType, period, limit)
	if err != nil {
		return nil, err
	}

	// 设置缓存（异步）
	go func() {
		c.cache.SetRanking(context.Background(), rankingType, period, items, RankingCacheExpiration)
	}()

	return items, nil
}

// 以下方法直接委托给原服务，不使用缓存

func (c *CachedBookstoreService) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error) {
	return c.service.GetBooksByCategory(ctx, categoryID, page, pageSize)
}

func (c *CachedBookstoreService) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	return c.service.GetRecommendedBooks(ctx, page, pageSize)
}

func (c *CachedBookstoreService) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	return c.service.GetFeaturedBooks(ctx, page, pageSize)
}

func (c *CachedBookstoreService) SearchBooks(ctx context.Context, keyword string, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	return c.service.SearchBooks(ctx, keyword, filter)
}

func (c *CachedBookstoreService) GetBookStats(ctx context.Context) (*bookstore.BookStats, error) {
	return c.service.GetBookStats(ctx)
}

func (c *CachedBookstoreService) IncrementBookView(ctx context.Context, bookID string) error {
	// 增加浏览量后，清除相关缓存
	err := c.service.IncrementBookView(ctx, bookID)
	if err != nil {
		return err
	}

	// 异步清除缓存
	go func() {
		c.cache.InvalidateBookCache(context.Background(), bookID)
		c.cache.InvalidateHomepageCache(context.Background())
	}()

	return nil
}

func (c *CachedBookstoreService) GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error) {
	return c.service.GetCategoryByID(ctx, id)
}

func (c *CachedBookstoreService) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	return c.service.GetRootCategories(ctx)
}

func (c *CachedBookstoreService) IncrementBannerClick(ctx context.Context, bannerID string) error {
	// 增加点击次数后，清除Banner缓存
	err := c.service.IncrementBannerClick(ctx, bannerID)
	if err != nil {
		return err
	}

	// 异步清除缓存
	go func() {
		c.cache.InvalidateBannerCache(context.Background())
		c.cache.InvalidateHomepageCache(context.Background())
	}()

	return nil
}

func (c *CachedBookstoreService) UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string) error {
	// 更新榜单后，清除相关缓存
	err := c.service.UpdateRankings(ctx, rankingType, period)
	if err != nil {
		return err
	}

	// 异步清除缓存
	go func() {
		c.cache.InvalidateRankingCache(context.Background(), rankingType, period)
		c.cache.InvalidateHomepageCache(context.Background())
	}()

	return nil
}