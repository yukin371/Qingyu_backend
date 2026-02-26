package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	searchModels "Qingyu_backend/models/search"
	"Qingyu_backend/models/shared/types"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"


	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// 元数据查询 - 用于筛选项
	GetYears(ctx context.Context) ([]int, error)
	GetTags(ctx context.Context, categoryID *string) ([]string, error)

	// 搜索相关 - 集成SearchService
	SearchByTitle(ctx context.Context, title string, page, size int) ([]*bookstore2.Book, int64, error)
	SearchByAuthor(ctx context.Context, author string, page, size int) ([]*bookstore2.Book, int64, error)
	GetSimilarBooks(ctx context.Context, bookID string, limit int) ([]*bookstore2.Book, error)
}

// BookstoreServiceImpl 书城服务实现
type BookstoreServiceImpl struct {
	bookRepo      BookstoreRepo.BookRepository
	categoryRepo  BookstoreRepo.CategoryRepository
	bannerRepo    BookstoreRepo.BannerRepository
	rankingRepo   BookstoreRepo.RankingRepository
	searchService interface{} // SearchService接口（避免循环依赖）
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

// SetSearchService 设置SearchService（避免循环依赖）
func (s *BookstoreServiceImpl) SetSearchService(searchService interface{}) {
	s.searchService = searchService
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
	// Repository 层现在接受 string 类型的 ID
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("[DEBUG] GetBookByID(%s) repository error: %v\n", id, err)
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	if book == nil {
		fmt.Printf("[DEBUG] GetBookByID(%s) book not found (nil)\n", id)
		return nil, errors.New("book not found")
	}

	fmt.Printf("[DEBUG] GetBookByID(%s) found book: %s, status: %s\n", id, book.Title, book.Status)

	// 只有连载中和已完结的书籍可以访问
	if book.Status != bookstore2.BookStatusOngoing && book.Status != bookstore2.BookStatusCompleted {
		fmt.Printf("[DEBUG] GetBookByID(%s) book status check failed: %s not in [ongoing, completed]\n", id, book.Status)
		return nil, errors.New("book not available")
	}

	fmt.Printf("[DEBUG] GetBookByID(%s) returning book successfully\n", id)
	return book, nil
}

// GetBooksByCategory 根据分类获取书籍列表
func (s *BookstoreServiceImpl) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取书籍列表 - Repository 层接受 string 类型的 categoryID
	books, err := s.bookRepo.GetByCategory(ctx, categoryID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by category: %w", err)
	}

	// 获取总数
	total, err := s.bookRepo.CountByCategory(ctx, categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by category: %w", err)
	}

	// 过滤只返回已发布的书籍（连载中或已完结）
	var publishedBooks []*bookstore2.Book
	for _, book := range books {
		if book.Status == bookstore2.BookStatusOngoing || book.Status == bookstore2.BookStatusCompleted {
			publishedBooks = append(publishedBooks, book)
		}
	}

	return publishedBooks, total, nil
}

// GetBooksByAuthorID 根据作者ID获取书籍列表
func (s *BookstoreServiceImpl) GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取书籍列表 - Repository 层接受 string 类型的 authorID
	books, err := s.bookRepo.GetByAuthorID(ctx, authorID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by author: %w", err)
	}

	// 获取总数
	total, err := s.bookRepo.CountByAuthor(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by author: %w", err)
	}

	// 过滤只返回已发布的书籍（连载中或已完结）
	var publishedBooks []*bookstore2.Book
	for _, book := range books {
		if book.Status == bookstore2.BookStatusOngoing || book.Status == bookstore2.BookStatusCompleted {
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
	ongoingStatus := bookstore2.BookStatusOngoing
	keywordPtr := keyword
	filter := &bookstore2.BookFilter{
		Keyword: &keywordPtr,
		Status:  &ongoingStatus,
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

	// 如果有关键词，先在数据库中获取所有符合其他条件的书籍，然后在Go代码中过滤关键词
	// 这样可以避免MongoDB正则表达式的UTF-8编码问题
	var books []*bookstore2.Book
	var err error

	if filter.Keyword != nil && *filter.Keyword != "" {
		// 创建一个没有关键词的过滤器，并设置为不分页
		filterWithoutKeyword := *filter
		filterWithoutKeyword.Keyword = nil
		filterWithoutKeyword.Limit = 0  // 不限制数量
		filterWithoutKeyword.Offset = 0

		// 获取所有符合其他条件的书籍
		allBooks, err := s.bookRepo.SearchWithFilter(ctx, &filterWithoutKeyword)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to search books: %w", err)
		}

		fmt.Printf("[DEBUG] 获取到 %d 本书籍，搜索关键词: %s\n", len(allBooks), *filter.Keyword)

		// 在Go代码中过滤关键词和其他条件
		keyword := *filter.Keyword  // 不使用ToLower，因为中文没有大小写
		filteredBooks := make([]*bookstore2.Book, 0)
		matchCount := 0
		for _, book := range allBooks {
			// 检查关键词匹配（不区分大小写）
			keywordMatch := strings.Contains(strings.ToLower(book.Title), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(book.Author), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(book.Introduction), strings.ToLower(keyword))

			if keywordMatch {
				matchCount++
				fmt.Printf("[DEBUG] 匹配到书籍: %s (作者: %s)\n", book.Title, book.Author)
			}

			if !keywordMatch {
				continue
			}

			// 检查其他条件
			if filter.Status != nil && book.Status != *filter.Status {
				continue
			}
			if filter.CategoryID != nil {
				found := false
				for _, catID := range book.CategoryIDs {
					if catID.Hex() == *filter.CategoryID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if filter.Author != nil && book.Author != *filter.Author {
				continue
			}
			if len(filter.Tags) > 0 {
				tagMatch := false
				for _, tag := range filter.Tags {
					for _, bookTag := range book.Tags {
						if bookTag == tag {
							tagMatch = true
							break
						}
					}
					if tagMatch {
						break
					}
				}
				if !tagMatch {
					continue
				}
			}

			filteredBooks = append(filteredBooks, book)
		}

		// 应用排序
		// TODO: 实现排序逻辑

		// 应用分页
		total := int64(len(filteredBooks))
		start := filter.Offset
		end := start + filter.Limit
		if end > len(filteredBooks) {
			end = len(filteredBooks)
		}
		if start >= len(filteredBooks) {
			return []*bookstore2.Book{}, total, nil
		}

		books = filteredBooks[start:end]
		return books, total, nil
	}

	// 没有关键词，直接查询
	books, err = s.bookRepo.SearchWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	// 计算总数
	total, err := s.bookRepo.CountByFilter(ctx, filter)
	if err != nil {
		return books, int64(len(books)), nil // 降级：使用当前页数量
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
	// 先检查书籍是否存在且已发布
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	if book == nil {
		return errors.New("book not found")
	}

	if book.Status != bookstore2.BookStatusOngoing && book.Status != bookstore2.BookStatusCompleted {
		return errors.New("book not available")
	}

	// 增加浏览量
	err = s.bookRepo.IncrementViewCount(ctx, bookID)
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
	// 直接使用 string ID，无需转换
	category, err := s.categoryRepo.GetByID(ctx, id)
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
	// 先检查Banner是否存在且激活
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
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
	err = s.bannerRepo.IncrementClickCount(ctx, bannerID)
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
	var firstError error
	for i := 0; i < 8; i++ {
		if err := <-errChan; err != nil {
			// 记录错误但不中断，让其他goroutine继续
			fmt.Printf("[WARNING] GetHomepageData部分数据获取失败: %v\n", err)
			if firstError == nil {
				firstError = err
			}
		}
	}

	// 只在所有关键数据都失败时才返回错误
	// 关键数据：推荐书籍、精选书籍
	if data.RecommendedBooks == nil && data.FeaturedBooks == nil {
		return nil, firstError
	}
	// 其他数据缺失是可以接受的，返回部分数据
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

// GetYears 获取所有书籍的发布年份列表（去重，倒序）
func (s *BookstoreServiceImpl) GetYears(ctx context.Context) ([]int, error) {
	years, err := s.bookRepo.GetYears(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get years: %w", err)
	}
	return years, nil
}

// GetTags 获取所有标签列表（去重，排序）
// 如果提供了 categoryID，则只返回该分类下的书籍标签
func (s *BookstoreServiceImpl) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	tags, err := s.bookRepo.GetTags(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

func boolPtr(b bool) *bool {
	return &b
}

// ========== 搜索相关辅助方法 ==========

// BuildSearchFilter 构建搜索过滤条件
// 将前端参数转换为后端过滤条件
func BuildSearchFilter(categoryID, author, status string, tags []string) map[string]interface{} {
	filter := make(map[string]interface{})

	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	if author != "" {
		filter["author"] = author
	}

	if status != "" {
		// 映射前端状态值到后端状态值
		var backendStatus string
		switch status {
		case "serializing":
			backendStatus = "ongoing"
		case "completed", "paused":
			backendStatus = status
		default:
			backendStatus = status
		}
		filter["status"] = backendStatus
	}

	if len(tags) > 0 {
		filter["tags"] = tags
	}

	return filter
}

// BuildSearchSort 构建搜索排序条件
func BuildSearchSort(sortBy, sortOrder string) []searchModels.SortField {
	var ascending bool
	if sortOrder == "asc" {
		ascending = true
	}

	return []searchModels.SortField{
		{
			Field:     sortBy,
			Ascending: ascending,
		},
	}
}

// DeduplicateAndLimitBooks 去重并限制数量
func DeduplicateAndLimitBooks(books []*bookstore2.Book, excludeID string, limit int) []*bookstore2.Book {
	seen := make(map[string]bool)
	result := make([]*bookstore2.Book, 0, len(books))

	for _, b := range books {
		bookID := b.ID.Hex()
		// 排除指定ID和已添加的书籍
		if bookID != excludeID && !seen[bookID] {
			result = append(result, b)
			seen[bookID] = true
		}
		if len(result) >= limit {
			break
		}
	}
	return result
}

// ConvertSearchResponseToBooks 将搜索响应转换为 Book 切片
func ConvertSearchResponseToBooks(items []searchModels.SearchItem) []*bookstore2.Book {
	books := make([]*bookstore2.Book, 0, len(items))

	for _, item := range items {
		book := &bookstore2.Book{}

		// 从 Data 中提取字段
		if id, ok := item.Data["id"].(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				book.ID = objectID
			}
		}
		if title, ok := item.Data["title"].(string); ok {
			book.Title = title
		}
		if author, ok := item.Data["author"].(string); ok {
			book.Author = author
		}
		if intro, ok := item.Data["introduction"].(string); ok {
			book.Introduction = intro
		}
		if coverURL, ok := item.Data["cover_url"].(string); ok {
			book.Cover = coverURL
		}
		if viewCount, ok := item.Data["view_count"].(int64); ok {
			book.ViewCount = viewCount
		}
		if rating, ok := item.Data["rating"].(float64); ok {
			book.Rating = types.Rating(rating)
		}
		if wordCount, ok := item.Data["word_count"].(int64); ok {
			book.WordCount = wordCount
		}
		if status, ok := item.Data["status"].(string); ok {
			book.Status = bookstore2.BookStatus(status)
		}

		books = append(books, book)
	}

	return books
}

// SearchSimilarBooks 搜索相似书籍
// 这是一个辅助方法，被GetSimilarBooks调用
// 参数：
//   - book: 原书籍
//   - limit: 返回数量限制
//   - useCategory: 是否使用分类过滤
//   - useTags: 是否使用标签过滤
func SearchSimilarBooks(
	ctx context.Context,
	book *bookstore2.Book,
	limit int,
	useCategory bool,
	useTags bool,
	bookRepo BookstoreRepo.BookRepository,
) []*bookstore2.Book {

	filter := &bookstore2.BookFilter{
		Limit:     limit + 1, // 多查一条用于排除当前书籍
		Offset:    0,
		SortBy:    "view_count",
		SortOrder: "desc",
	}

	// 注意：Book 模型使用 CategoryIDs 数组而非单个 CategoryID
	// 如果书籍有分类ID，使用第一个分类ID进行查询
	if useCategory && len(book.CategoryIDs) > 0 {
		firstCategoryID := book.CategoryIDs[0].Hex()
		filter.CategoryID = &firstCategoryID
	}

	if useTags && len(book.Tags) > 0 {
		filter.Tags = book.Tags
	}

	books, err := bookRepo.SearchWithFilter(ctx, filter)
	if err != nil {
		return []*bookstore2.Book{}
	}

	return books
}

// ========== SearchService集成方法 ==========

// SearchByTitle 按标题搜索书籍
// 优先使用SearchService (Milvus向量搜索)，失败或空结果时fallback到MongoDB
func (s *BookstoreServiceImpl) SearchByTitle(ctx context.Context, title string, page, size int) ([]*bookstore2.Book, int64, error) {
	// 优先使用新路径 (SearchService - Milvus向量搜索)
	// TODO: 实现SearchService集成
	_ = s.searchService // 暂时避免未使用警告

	// Fallback 到旧路径 (MongoDB查询)
	filter := &bookstore2.BookFilter{
		Keyword:   &title,
		SortBy:    "view_count",
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := s.SearchBooksWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books by title: %w", err)
	}

	// 确保返回空数组而不是nil
	if books == nil {
		books = make([]*bookstore2.Book, 0)
	}

	return books, total, nil
}

// SearchByAuthor 按作者搜索书籍
// 优先使用SearchService (Milvus向量搜索)，失败或空结果时fallback到MongoDB
func (s *BookstoreServiceImpl) SearchByAuthor(ctx context.Context, author string, page, size int) ([]*bookstore2.Book, int64, error) {
	// 优先使用新路径 (SearchService - Milvus向量搜索)
	// TODO: 实现SearchService集成
	_ = s.searchService // 暂时避免未使用警告

	// Fallback 到旧路径 (MongoDB查询)
	filter := &bookstore2.BookFilter{
		Author:    &author,
		SortBy:    "view_count",
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := s.SearchBooksWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books by author: %w", err)
	}

	// 确保返回空数组而不是nil
	if books == nil {
		books = make([]*bookstore2.Book, 0)
	}

	return books, total, nil
}

// GetSimilarBooks 获取相似书籍推荐
// 基于书籍分类、标签等推荐相似书籍，有四层降级策略
func (s *BookstoreServiceImpl) GetSimilarBooks(ctx context.Context, bookID string, limit int) ([]*bookstore2.Book, error) {
	// 获取原书籍信息
	book, err := s.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	var result []*bookstore2.Book

	// v1.2实现：四层降级策略
	// 策略1: 同分类 + 标签
	result = SearchSimilarBooks(ctx, book, limit, true, true, s.bookRepo)

	// 策略2: 如果结果不足，尝试同分类
	if len(result) < limit {
		additional := SearchSimilarBooks(ctx, book, limit-len(result), true, false, s.bookRepo)
		result = append(result, additional...)
	}

	// 策略3: 如果还不足，尝试标签匹配
	if len(result) < limit {
		additional := SearchSimilarBooks(ctx, book, limit-len(result), false, true, s.bookRepo)
		result = append(result, additional...)
	}

	// 策略4: 兜底 - 返回热门书籍（禁止返回空列表）
	if len(result) == 0 {
		popularBooks, _, err := s.GetHotBooks(ctx, 1, limit)
		if err == nil && len(popularBooks) > 0 {
			result = popularBooks
		} else {
			// 如果连热门书籍都获取失败，尝试推荐书籍
			recommendedBooks, _, err := s.GetRecommendedBooks(ctx, 1, limit)
			if err == nil && len(recommendedBooks) > 0 {
				result = recommendedBooks
			}
		}
	}

	// 去重并限制数量
	uniqueResult := DeduplicateAndLimitBooks(result, bookID, limit)

	return uniqueResult, nil
}
