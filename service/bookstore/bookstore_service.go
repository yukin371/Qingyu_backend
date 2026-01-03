package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// BookstoreService 书城服务接口 - 专注于书城列表展示和首页聚合
// 用于书城首页、分类页面、搜索结果等列表场景
type BookstoreService interface {
	// 书籍列表相关服务 - 使用Book模型
	GetAllBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetBookByID(ctx context.Context, id string) (*bookstore2.Book, error)
	GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error)
	SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore2.Book, int64, error)
	SearchBooksWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, int64, error)

	// 分类相关服务
	GetCategoryTree(ctx context.Context) ([]*bookstore2.CategoryTree, error)
	GetCategoryByID(ctx context.Context, id string) (*bookstore2.Category, error)
	GetRootCategories(ctx context.Context) ([]*bookstore2.Category, error)

	// Banner相关方法
	GetActiveBanners(ctx context.Context, limit int) ([]*bookstore2.Banner, error)
	IncrementBannerClick(ctx context.Context, bannerID string) error

	// 榜单相关方法
	GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error)
	GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error)
	GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error)
	UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error

	// 首页数据聚合
	GetHomepageData(ctx context.Context) (*HomepageData, error)

	// 统计和计数
	GetBookStats(ctx context.Context) (*bookstore2.BookStats, error)
	IncrementBookView(ctx context.Context, bookID string) error
}

// BookstoreServiceImpl 书城服务实现
type BookstoreServiceImpl struct {
	bookRepo     BookstoreRepo.BookRepository
	categoryRepo BookstoreRepo.CategoryRepository
	bannerRepo   BookstoreRepo.BannerRepository
	rankingRepo  BookstoreRepo.RankingRepository
}

// HomepageData 首页数据结构
type HomepageData struct {
	Banners          []*bookstore2.Banner                 `json:"banners"`
	RecommendedBooks []*bookstore2.Book                   `json:"recommendedBooks"`
	FeaturedBooks    []*bookstore2.Book                   `json:"featuredBooks"`
	Categories       []*bookstore2.Category               `json:"categories"`
	Stats            *bookstore2.BookStats                `json:"stats"`
	Rankings         map[string][]*bookstore2.RankingItem `json:"rankings"` // 各类榜单
}

// NewBookstoreService 创建书城服务实例
func NewBookstoreService(
	bookRepo BookstoreRepo.BookRepository,
	categoryRepo BookstoreRepo.CategoryRepository,
	bannerRepo BookstoreRepo.BannerRepository,
	rankingRepo BookstoreRepo.RankingRepository,
) BookstoreService {
	return &BookstoreServiceImpl{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
		bannerRepo:   bannerRepo,
		rankingRepo:  rankingRepo,
	}
}

// GetAllBooks 获取所有书籍列表（分页）
func (s *BookstoreServiceImpl) GetAllBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	offset := (page - 1) * pageSize

	// 获取所有已发布的书籍
	books, err := s.bookRepo.GetHotBooks(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	// 获取总数
	total, err := s.bookRepo.CountByFilter(ctx, &bookstore2.BookFilter{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	return books, total, nil
}

// GetBookByID 根据ID获取书籍详情
func (s *BookstoreServiceImpl) GetBookByID(ctx context.Context, id string) (*bookstore2.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	book, err := s.bookRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	// 只返回已发布的书籍
	if book.Status != "published" {
		return nil, errors.New("book not available")
	}

	return book, nil
}

// GetBooksByCategory 根据分类获取书籍列表
func (s *BookstoreServiceImpl) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid category ID: %w", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取书籍列表
	books, err := s.bookRepo.GetByCategory(ctx, objectID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by category: %w", err)
	}

	// 获取总数
	total, err := s.bookRepo.CountByCategory(ctx, objectID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by category: %w", err)
	}

	// 过滤只返回已发布的书籍
	var publishedBooks []*bookstore2.Book
	for _, book := range books {
		if book.Status == "published" {
			publishedBooks = append(publishedBooks, book)
		}
	}

	return publishedBooks, total, nil
}

// GetBooksByAuthorID 根据作者ID获取书籍列表
func (s *BookstoreServiceImpl) GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid author ID: %w", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取书籍列表
	books, err := s.bookRepo.GetByAuthorID(ctx, objectID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by author: %w", err)
	}

	// 获取总数
	total, err := s.bookRepo.CountByAuthor(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by author: %w", err)
	}

	// 过滤只返回已发布的书籍
	var publishedBooks []*bookstore2.Book
	for _, book := range books {
		if book.Status == "published" {
			publishedBooks = append(publishedBooks, book)
		}
	}

	return publishedBooks, total, nil
}

// GetRecommendedBooks 获取所有书籍（书库列表）
// 修改：返回所有书籍，按创建时间倒序排序
func (s *BookstoreServiceImpl) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	offset := (page - 1) * pageSize

	// 直接调用 repository 的 GetRecommended 方法
	// repository 已修改为返回所有书籍
	books, err := s.bookRepo.GetRecommended(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	// 获取所有书籍总数
	total, err := s.bookRepo.Count(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	return books, total, nil
}

// GetFeaturedBooks 获取精选书籍
func (s *BookstoreServiceImpl) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	// 检查repository是否为nil
	if s.bookRepo == nil {
		return []*bookstore2.Book{}, 0, nil
	}

	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetFeatured(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get featured books: %w", err)
	}

	// 获取精选书籍总数
	filter := &bookstore2.BookFilter{
		IsFeatured: boolPtr(true),
	}
	total, err := s.bookRepo.CountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count featured books: %w", err)
	}

	return books, total, nil
}

// GetHotBooks 获取热门书籍
func (s *BookstoreServiceImpl) GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetHotBooks(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get hot books: %w", err)
	}

	// 获取热门书籍总数（所有已发布状态的书籍）
	total, err := s.bookRepo.CountByFilter(ctx, &bookstore2.BookFilter{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count hot books: %w", err)
	}

	return books, total, nil
}

// GetNewReleases 获取新书列表
func (s *BookstoreServiceImpl) GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetNewReleases(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get new releases: %w", err)
	}

	// 获取新书总数（所有已发布状态的书籍）
	total, err := s.bookRepo.CountByFilter(ctx, &bookstore2.BookFilter{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count new releases: %w", err)
	}

	return books, total, nil
}

// GetFreeBooks 获取免费书籍列表
func (s *BookstoreServiceImpl) GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetFreeBooks(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get free books: %w", err)
	}

	// 获取免费书籍总数
	filter := &bookstore2.BookFilter{
		IsFree: boolPtr(true),
	}
	total, err := s.bookRepo.CountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count free books: %w", err)
	}

	return books, total, nil
}

// SearchBooks 搜索书籍 - 简单搜索
func (s *BookstoreServiceImpl) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("keyword is required")
	}

	offset := (page - 1) * pageSize

	// 调用Repository的Search方法
	books, err := s.bookRepo.Search(ctx, keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	// 计算总数
	publishedStatus := bookstore2.BookStatusPublished
	keywordPtr := keyword
	filter := &bookstore2.BookFilter{
		Keyword: &keywordPtr,
		Status:  &publishedStatus,
	}
	total, err := s.bookRepo.CountByFilter(ctx, filter)
	if err != nil {
		return books, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return books, total, nil
}

// SearchBooksWithFilter 搜索书籍 - 高级搜索
func (s *BookstoreServiceImpl) SearchBooksWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, int64, error) {
	if filter == nil {
		return nil, 0, errors.New("filter is required")
	}

	// 确保只搜索已发布的书籍
	if filter.Status == nil {
		publishedStatus := bookstore2.BookStatusPublished
		filter.Status = &publishedStatus
	}

	// 执行搜索
	books, err := s.bookRepo.SearchWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	// 计算总数
	total, err := s.bookRepo.CountByFilter(ctx, filter)
	if err != nil {
		return books, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return books, total, nil
}

// GetBookStats 获取书籍统计信息
func (s *BookstoreServiceImpl) GetBookStats(ctx context.Context) (*bookstore2.BookStats, error) {
	// 检查repository是否为nil
	if s.bookRepo == nil {
		return &bookstore2.BookStats{}, nil
	}

	stats, err := s.bookRepo.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get book stats: %w", err)
	}

	return stats, nil
}

// IncrementBookView 增加书籍浏览量
func (s *BookstoreServiceImpl) IncrementBookView(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}

	// 先检查书籍是否存在且已发布
	book, err := s.bookRepo.GetByID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	if book == nil {
		return errors.New("book not found")
	}

	if book.Status != "published" {
		return errors.New("book not available")
	}

	// 增加浏览量
	err = s.bookRepo.IncrementViewCount(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	return nil
}

// GetCategoryTree 获取分类树
func (s *BookstoreServiceImpl) GetCategoryTree(ctx context.Context) ([]*bookstore2.CategoryTree, error) {
	tree, err := s.categoryRepo.GetCategoryTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category tree: %w", err)
	}

	return tree, nil
}

// GetCategoryByID 根据ID获取分类
func (s *BookstoreServiceImpl) GetCategoryByID(ctx context.Context, id string) (*bookstore2.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	category, err := s.categoryRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	if !category.IsActive {
		return nil, errors.New("category not available")
	}

	return category, nil
}

// GetRootCategories 获取根分类列表
func (s *BookstoreServiceImpl) GetRootCategories(ctx context.Context) ([]*bookstore2.Category, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}

	return categories, nil
}

// GetActiveBanners 获取激活的Banner列表
func (s *BookstoreServiceImpl) GetActiveBanners(ctx context.Context, limit int) ([]*bookstore2.Banner, error) {
	// 检查repository是否为nil
	if s.bannerRepo == nil {
		return []*bookstore2.Banner{}, nil
	}

	banners, err := s.bannerRepo.GetActive(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get active banners: %w", err)
	}

	return banners, nil
}

// IncrementBannerClick 增加Banner点击次数
func (s *BookstoreServiceImpl) IncrementBannerClick(ctx context.Context, bannerID string) error {
	objectID, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		return fmt.Errorf("invalid banner ID: %w", err)
	}

	// 先检查Banner是否存在且激活
	banner, err := s.bannerRepo.GetByID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return errors.New("banner not found")
	}

	if !banner.IsActive {
		return errors.New("banner not available")
	}

	// 检查时间范围
	now := time.Now()
	if banner.StartTime != nil && now.Before(*banner.StartTime) {
		return errors.New("banner not started")
	}
	if banner.EndTime != nil && now.After(*banner.EndTime) {
		return errors.New("banner expired")
	}

	// 增加点击次数
	err = s.bannerRepo.IncrementClickCount(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	return nil
}

// GetHomepageData 获取首页数据
func (s *BookstoreServiceImpl) GetHomepageData(ctx context.Context) (*HomepageData, error) {
	data := &HomepageData{
		Rankings: make(map[string][]*bookstore2.RankingItem),
	}

	// 并发获取各种数据
	errChan := make(chan error, 8) // 增加到8个goroutine

	// 获取Banner
	go func() {
		banners, err := s.GetActiveBanners(ctx, 5)
		if err != nil {
			errChan <- fmt.Errorf("failed to get banners: %w", err)
			return
		}
		data.Banners = banners
		errChan <- nil
	}()

	// 获取推荐书籍
	go func() {
		books, _, err := s.GetRecommendedBooks(ctx, 1, 10)
		if err != nil {
			errChan <- fmt.Errorf("failed to get recommended books: %w", err)
			return
		}
		data.RecommendedBooks = books
		errChan <- nil
	}()

	// 获取精选书籍
	go func() {
		books, _, err := s.GetFeaturedBooks(ctx, 1, 10)
		if err != nil {
			errChan <- fmt.Errorf("failed to get featured books: %w", err)
			return
		}
		data.FeaturedBooks = books
		errChan <- nil
	}()

	// 获取根分类
	go func() {
		categories, err := s.GetRootCategories(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to get categories: %w", err)
			return
		}
		data.Categories = categories
		errChan <- nil
	}()

	// 获取统计信息
	go func() {
		stats, err := s.GetBookStats(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to get stats: %w", err)
			return
		}
		data.Stats = stats
		errChan <- nil
	}()

	// 使用mutex保护并发写入map
	var mu sync.Mutex

	// 获取实时榜
	go func() {
		rankings, err := s.GetRealtimeRanking(ctx, 10)
		if err != nil {
			errChan <- fmt.Errorf("failed to get realtime ranking: %w", err)
			return
		}
		mu.Lock()
		data.Rankings["realtime"] = rankings
		mu.Unlock()
		errChan <- nil
	}()

	// 获取周榜
	go func() {
		rankings, err := s.GetWeeklyRanking(ctx, "", 10)
		if err != nil {
			errChan <- fmt.Errorf("failed to get weekly ranking: %w", err)
			return
		}
		mu.Lock()
		data.Rankings["weekly"] = rankings
		mu.Unlock()
		errChan <- nil
	}()

	// 获取月榜
	go func() {
		rankings, err := s.GetMonthlyRanking(ctx, "", 10)
		if err != nil {
			errChan <- fmt.Errorf("failed to get monthly ranking: %w", err)
			return
		}
		mu.Lock()
		data.Rankings["monthly"] = rankings
		mu.Unlock()
		errChan <- nil
	}()

	// 等待所有goroutine完成
	for i := 0; i < 8; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	return data, nil
}

// GetRealtimeRanking 获取实时榜
func (s *BookstoreServiceImpl) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error) {
	// 检查repository是否为nil
	if s.rankingRepo == nil {
		return []*bookstore2.RankingItem{}, nil
	}

	period := bookstore2.GetPeriodString(bookstore2.RankingTypeRealtime, time.Now())
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeRealtime, period, limit, 0)
}

// GetWeeklyRanking 获取周榜
func (s *BookstoreServiceImpl) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	// 检查repository是否为nil
	if s.rankingRepo == nil {
		return []*bookstore2.RankingItem{}, nil
	}

	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeWeekly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeWeekly, period, limit, 0)
}

// GetMonthlyRanking 获取月榜
func (s *BookstoreServiceImpl) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	// 检查repository是否为nil
	if s.rankingRepo == nil {
		return []*bookstore2.RankingItem{}, nil
	}

	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeMonthly, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeMonthly, period, limit, 0)
}

// GetNewbieRanking 获取新人榜
func (s *BookstoreServiceImpl) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	// 检查repository是否为nil
	if s.rankingRepo == nil {
		return []*bookstore2.RankingItem{}, nil
	}

	if period == "" {
		period = bookstore2.GetPeriodString(bookstore2.RankingTypeNewbie, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, bookstore2.RankingTypeNewbie, period, limit, 0)
}

// GetRankingByType 根据类型获取榜单
func (s *BookstoreServiceImpl) GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	if period == "" {
		period = bookstore2.GetPeriodString(rankingType, time.Now())
	}
	return s.rankingRepo.GetByTypeWithBooks(ctx, rankingType, period, limit, 0)
}

// UpdateRankings 更新榜单数据
func (s *BookstoreServiceImpl) UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error {
	var items []*bookstore2.RankingItem
	var err error

	switch rankingType {
	case bookstore2.RankingTypeRealtime:
		items, err = s.rankingRepo.CalculateRealtimeRanking(ctx, period)
	case bookstore2.RankingTypeWeekly:
		items, err = s.rankingRepo.CalculateWeeklyRanking(ctx, period)
	case bookstore2.RankingTypeMonthly:
		items, err = s.rankingRepo.CalculateMonthlyRanking(ctx, period)
	case bookstore2.RankingTypeNewbie:
		items, err = s.rankingRepo.CalculateNewbieRanking(ctx, period)
	default:
		return fmt.Errorf("unsupported ranking type: %s", rankingType)
	}

	if err != nil {
		return fmt.Errorf("failed to calculate ranking: %w", err)
	}

	return s.rankingRepo.UpdateRankings(ctx, rankingType, period, items)
}

// boolPtr 返回bool值的指针
func boolPtr(b bool) *bool {
	return &b
}
