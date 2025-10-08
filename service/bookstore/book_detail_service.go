package bookstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// BookDetailFilter 书籍详情筛选条件 - 适用于网络小说平台
type BookDetailFilter struct {
	Title          string                `json:"title,omitempty"`
	Author         string                `json:"author,omitempty"`
	AuthorID       *primitive.ObjectID   `json:"author_id,omitempty"`
	CategoryIDs    []primitive.ObjectID  `json:"category_ids,omitempty"`
	Tags           []string              `json:"tags,omitempty"`
	Status         *bookstore.BookStatus `json:"status,omitempty"`
	IsFree         *bool                 `json:"is_free,omitempty"`
	MinPrice       *float64              `json:"min_price,omitempty"`
	MaxPrice       *float64              `json:"max_price,omitempty"`
	MinRating      *float64              `json:"min_rating,omitempty"`
	MaxRating      *float64              `json:"max_rating,omitempty"`
	MinWordCount   *int64                `json:"min_word_count,omitempty"`
	MaxWordCount   *int64                `json:"max_word_count,omitempty"`
	SerializedFrom *time.Time            `json:"serialized_from,omitempty"` // 开始连载时间范围
	SerializedTo   *time.Time            `json:"serialized_to,omitempty"`   // 开始连载时间范围
	CompletedFrom  *time.Time            `json:"completed_from,omitempty"`  // 完结时间范围
	CompletedTo    *time.Time            `json:"completed_to,omitempty"`    // 完结时间范围
	CreatedAtFrom  *time.Time            `json:"created_at_from,omitempty"`
	CreatedAtTo    *time.Time            `json:"created_at_to,omitempty"`
	UpdatedAtFrom  *time.Time            `json:"updated_at_from,omitempty"`
	UpdatedAtTo    *time.Time            `json:"updated_at_to,omitempty"`
	SortBy         string                `json:"sort_by,omitempty"`    // created_at, updated_at, serialized_at, rating, word_count, view_count, like_count, collect_count
	SortOrder      string                `json:"sort_order,omitempty"` // asc, desc
}

// BookDetailService 书籍详情服务接口 - 专注于书籍详情页面的完整信息管理
// 用于书籍详情页面、章节管理、统计数据等详细场景
type BookDetailService interface {
	// 书籍详情基础操作
	CreateBookDetail(ctx context.Context, bookDetail *bookstore.BookDetail) error
	GetBookDetailByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookDetail, error)
	UpdateBookDetail(ctx context.Context, bookDetail *bookstore.BookDetail) error
	DeleteBookDetail(ctx context.Context, id primitive.ObjectID) error

	// 书籍详情查询
	GetBookDetailByTitle(ctx context.Context, title string) (*bookstore.BookDetail, error)
	GetBookDetailsByAuthor(ctx context.Context, author string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBookDetailsByAuthorID(ctx context.Context, authorID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBookDetailsByCategory(ctx context.Context, category string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBookDetailsByStatus(ctx context.Context, status bookstore.BookStatus, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBookDetailsByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	SearchBookDetails(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	SearchBookDetailsWithFilter(ctx context.Context, filter *BookDetailFilter, page, pageSize int) ([]*bookstore.BookDetail, int64, error)

	// API 兼容方法别名（为了与 API 层命名保持一致）
	GetBooksByTitle(ctx context.Context, title string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBooksByAuthor(ctx context.Context, author string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBooksByCategory(ctx context.Context, category string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBooksByStatus(ctx context.Context, status string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetBooksByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookDetail, int64, error)
	GetRecommendedBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error)
	GetSimilarBooks(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error)
	GetPopularBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error)
	GetLatestBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error)
	CountBooksByCategory(ctx context.Context, category string) (int64, error)

	// 书籍详情统计和交互
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error
	UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64, ratingCount int64) error
	UpdateLastChapter(ctx context.Context, bookID primitive.ObjectID, chapterTitle string) error

	// 书籍详情统计查询
	GetBookDetailStats(ctx context.Context) (map[string]interface{}, error)
	GetBookDetailCountByCategory(ctx context.Context, category string) (int64, error)
	GetBookDetailCountByAuthor(ctx context.Context, author string) (int64, error)
	GetBookDetailCountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error)

	// 书籍详情推荐
	GetRecommendedBookDetails(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error)
	GetSimilarBookDetails(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error)
	GetPopularBookDetails(ctx context.Context, limit int) ([]*bookstore.BookDetail, error)
	GetLatestBookDetails(ctx context.Context, limit int) ([]*bookstore.BookDetail, error)

	// 书籍详情批量操作
	BatchUpdateBookDetailStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error
	BatchUpdateBookDetailTags(ctx context.Context, bookIDs []primitive.ObjectID, tags []string) error
}

// BookDetailServiceImpl 书籍详情服务实现
type BookDetailServiceImpl struct {
	bookDetailRepo BookstoreRepo.BookDetailRepository
	cacheService   CacheService
}

// NewBookDetailService 创建书籍详情服务实例
func NewBookDetailService(bookDetailRepo BookstoreRepo.BookDetailRepository, cacheService CacheService) BookDetailService {
	return &BookDetailServiceImpl{
		bookDetailRepo: bookDetailRepo,
		cacheService:   cacheService,
	}
}

// CreateBookDetail 创建书籍详情
func (s *BookDetailServiceImpl) CreateBookDetail(ctx context.Context, bookDetail *bookstore.BookDetail) error {
	if bookDetail == nil {
		return errors.New("book detail cannot be nil")
	}

	// 验证必填字段
	if bookDetail.Title == "" {
		return errors.New("book title is required")
	}
	if bookDetail.Author == "" {
		return errors.New("book author is required")
	}

	// 检查标题是否已存在
	existingBook, err := s.bookDetailRepo.GetByTitle(ctx, bookDetail.Title)
	if err != nil {
		return fmt.Errorf("failed to check existing book: %w", err)
	}
	if existingBook != nil {
		return errors.New("book with this title already exists")
	}

	// 创建书籍详情
	if err := s.bookDetailRepo.Create(ctx, bookDetail); err != nil {
		return fmt.Errorf("failed to create book detail: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, bookDetail)

	return nil
}

// GetBookDetailByID 根据ID获取书籍详情
func (s *BookDetailServiceImpl) GetBookDetailByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookDetail, error) {
	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedBook, err := s.cacheService.GetBookDetail(ctx, id.Hex()); err == nil && cachedBook != nil {
			return cachedBook, nil
		}
	}

	// 从数据库获取
	bookDetail, err := s.bookDetailRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book detail: %w", err)
	}
	if bookDetail == nil {
		return nil, errors.New("book detail not found")
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetBookDetail(ctx, id.Hex(), bookDetail, 30*time.Minute)
	}

	return bookDetail, nil
}

// UpdateBookDetail 更新书籍详情
func (s *BookDetailServiceImpl) UpdateBookDetail(ctx context.Context, bookDetail *bookstore.BookDetail) error {
	if bookDetail == nil {
		return errors.New("book detail cannot be nil")
	}

	// 验证必填字段
	if bookDetail.ID.IsZero() {
		return errors.New("book detail ID is required")
	}
	if bookDetail.Title == "" {
		return errors.New("book title is required")
	}
	if bookDetail.Author == "" {
		return errors.New("book author is required")
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"title":         bookDetail.Title,
		"subtitle":      bookDetail.Subtitle,
		"author":        bookDetail.Author,
		"author_id":     bookDetail.AuthorID,
		"description":   bookDetail.Description,
		"cover_url":     bookDetail.CoverURL,
		"categories":    bookDetail.Categories,
		"tags":          bookDetail.Tags,
		"status":        bookDetail.Status,
		"word_count":    bookDetail.WordCount,
		"chapter_count": bookDetail.ChapterCount,
		"price":         bookDetail.Price,
		"is_free":       bookDetail.IsFree,
		"serialized_at": bookDetail.SerializedAt,
		"completed_at":  bookDetail.CompletedAt,
		"updated_at":    time.Now(),
	}

	// 更新数据库
	if err := s.bookDetailRepo.Update(ctx, bookDetail.ID, updates); err != nil {
		return fmt.Errorf("failed to update book detail: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, bookDetail)

	return nil
}

// DeleteBookDetail 删除书籍详情
func (s *BookDetailServiceImpl) DeleteBookDetail(ctx context.Context, id primitive.ObjectID) error {
	// 先获取书籍详情用于清除缓存
	bookDetail, err := s.bookDetailRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book detail for deletion: %w", err)
	}
	if bookDetail == nil {
		return errors.New("book detail not found")
	}

	// 删除书籍详情
	if err := s.bookDetailRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete book detail: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, bookDetail)

	return nil
}

// GetBookDetailByTitle 根据标题获取书籍详情
func (s *BookDetailServiceImpl) GetBookDetailByTitle(ctx context.Context, title string) (*bookstore.BookDetail, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	bookDetail, err := s.bookDetailRepo.GetByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get book detail by title: %w", err)
	}

	return bookDetail, nil
}

// GetBookDetailsByAuthor 根据作者获取书籍详情列表
func (s *BookDetailServiceImpl) GetBookDetailsByAuthor(ctx context.Context, author string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if author == "" {
		return nil, 0, errors.New("author cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取书籍列表
	bookDetails, err := s.bookDetailRepo.GetByAuthor(ctx, author, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book details by author: %w", err)
	}

	// 获取总数
	total, err := s.bookDetailRepo.CountByAuthor(ctx, author)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book details by author: %w", err)
	}

	return bookDetails, total, nil
}

// GetBookDetailsByAuthorID 根据作者ID获取书籍详情列表
func (s *BookDetailServiceImpl) GetBookDetailsByAuthorID(ctx context.Context, authorID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if authorID.IsZero() {
		return nil, 0, errors.New("author ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取书籍列表
	bookDetails, err := s.bookDetailRepo.GetByAuthorID(ctx, authorID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book details by author ID: %w", err)
	}

	// TODO: 需要在 Repository 添加 CountByAuthorID 方法
	// 暂时返回书籍列表长度作为总数
	total := int64(len(bookDetails))

	return bookDetails, total, nil
}

// GetBookDetailsByCategory 根据分类获取书籍详情列表
func (s *BookDetailServiceImpl) GetBookDetailsByCategory(ctx context.Context, category string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if category == "" {
		return nil, 0, errors.New("category cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取书籍列表
	bookDetails, err := s.bookDetailRepo.GetByCategory(ctx, category, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book details by category: %w", err)
	}

	// 获取总数
	total, err := s.bookDetailRepo.CountByCategory(ctx, category)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book details by category: %w", err)
	}

	return bookDetails, total, nil
}

// GetBookDetailsByStatus 根据状态获取书籍详情列表
func (s *BookDetailServiceImpl) GetBookDetailsByStatus(ctx context.Context, status bookstore.BookStatus, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取书籍列表
	bookDetails, err := s.bookDetailRepo.GetByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book details by status: %w", err)
	}

	// 获取总数
	total, err := s.bookDetailRepo.CountByStatus(ctx, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book details by status: %w", err)
	}

	return bookDetails, total, nil
}

// GetBookDetailsByTags 根据标签获取书籍详情列表
func (s *BookDetailServiceImpl) GetBookDetailsByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if len(tags) == 0 {
		return nil, 0, errors.New("tags cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取书籍列表
	bookDetails, err := s.bookDetailRepo.GetByTags(ctx, tags, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book details by tags: %w", err)
	}

	// 获取总数
	total, err := s.bookDetailRepo.CountByTags(ctx, tags)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book details by tags: %w", err)
	}

	return bookDetails, total, nil
}

// SearchBookDetails 搜索书籍详情
func (s *BookDetailServiceImpl) SearchBookDetails(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("keyword cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 搜索书籍
	bookDetails, err := s.bookDetailRepo.Search(ctx, keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search book details: %w", err)
	}

	// 这里简化处理，实际应该有专门的搜索计数方法
	total := int64(len(bookDetails))
	if len(bookDetails) == pageSize {
		// 如果返回的结果等于页面大小，可能还有更多结果
		total = int64((page + 1) * pageSize)
	}

	return bookDetails, total, nil
}

// SearchBookDetailsWithFilter 使用过滤器搜索书籍详情
func (s *BookDetailServiceImpl) SearchBookDetailsWithFilter(ctx context.Context, filter *BookDetailFilter, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	if filter == nil {
		return nil, 0, errors.New("filter cannot be nil")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// TODO: Repository 层需要实现 SearchWithFilter 方法
	// 暂时使用简单搜索方式
	var bookDetails []*bookstore.BookDetail
	var err error

	// 根据过滤器条件调用不同的查询方法
	if filter.Title != "" {
		bookDetails, err = s.bookDetailRepo.Search(ctx, filter.Title, pageSize, offset)
	} else if filter.Status != nil {
		bookDetails, err = s.bookDetailRepo.GetByStatus(ctx, *filter.Status, pageSize, offset)
	} else if filter.Author != "" {
		bookDetails, err = s.bookDetailRepo.GetByAuthor(ctx, filter.Author, pageSize, offset)
	} else {
		// 默认搜索
		bookDetails, err = s.bookDetailRepo.Search(ctx, "", pageSize, offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search book details with filter: %w", err)
	}

	// 暂时返回结果长度作为总数
	total := int64(len(bookDetails))

	return bookDetails, total, nil
}

// GetBookDetailStats 获取书籍详情统计
func (s *BookDetailServiceImpl) GetBookDetailStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 各状态书籍数量
	completedCount, err := s.bookDetailRepo.CountByStatus(ctx, bookstore.BookStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed book count: %w", err)
	}
	stats["completed_books"] = completedCount

	ongoingCount, err := s.bookDetailRepo.CountByStatus(ctx, bookstore.BookStatusOngoing)
	if err != nil {
		return nil, fmt.Errorf("failed to get ongoing book count: %w", err)
	}
	stats["ongoing_books"] = ongoingCount

	pausedCount, err := s.bookDetailRepo.CountByStatus(ctx, bookstore.BookStatusPaused)
	if err != nil {
		return nil, fmt.Errorf("failed to get paused book count: %w", err)
	}
	stats["paused_books"] = pausedCount

	// 计算总书籍数
	totalCount := completedCount + ongoingCount + pausedCount
	stats["total_books"] = totalCount

	return stats, nil
}

// GetBookDetailCountByCategory 根据分类统计书籍数量
func (s *BookDetailServiceImpl) GetBookDetailCountByCategory(ctx context.Context, category string) (int64, error) {
	if category == "" {
		return 0, errors.New("category cannot be empty")
	}

	count, err := s.bookDetailRepo.CountByCategory(ctx, category)
	if err != nil {
		return 0, fmt.Errorf("failed to count book details by category: %w", err)
	}

	return count, nil
}

// GetBookDetailCountByAuthor 根据作者统计书籍数量
func (s *BookDetailServiceImpl) GetBookDetailCountByAuthor(ctx context.Context, author string) (int64, error) {
	if author == "" {
		return 0, errors.New("author cannot be empty")
	}

	count, err := s.bookDetailRepo.CountByAuthor(ctx, author)
	if err != nil {
		return 0, fmt.Errorf("failed to count book details by author: %w", err)
	}

	return count, nil
}

// GetBookDetailCountByStatus 根据状态统计书籍数量
func (s *BookDetailServiceImpl) GetBookDetailCountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	count, err := s.bookDetailRepo.CountByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count book details by status: %w", err)
	}

	return count, nil
}

// GetRecommendedBookDetails 获取推荐书籍详情
func (s *BookDetailServiceImpl) GetRecommendedBookDetails(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 先获取当前书籍信息
	currentBook, err := s.bookDetailRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current book: %w", err)
	}
	if currentBook == nil {
		return nil, errors.New("current book not found")
	}

	// 基于分类推荐
	if len(currentBook.Categories) > 0 {
		recommendedBooks, err := s.bookDetailRepo.GetByCategory(ctx, currentBook.Categories[0], limit, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get recommended books: %w", err)
		}

		// 过滤掉当前书籍
		var filteredBooks []*bookstore.BookDetail
		for _, book := range recommendedBooks {
			if book.ID != bookID {
				filteredBooks = append(filteredBooks, book)
			}
		}

		return filteredBooks, nil
	}

	// 如果没有分类，返回最新书籍
	return s.GetLatestBookDetails(ctx, limit)
}

// GetSimilarBookDetails 获取相似书籍详情
func (s *BookDetailServiceImpl) GetSimilarBookDetails(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 先获取当前书籍信息
	currentBook, err := s.bookDetailRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current book: %w", err)
	}
	if currentBook == nil {
		return nil, errors.New("current book not found")
	}

	// 基于标签推荐
	if len(currentBook.Tags) > 0 {
		similarBooks, err := s.bookDetailRepo.GetByTags(ctx, currentBook.Tags, limit, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get similar books: %w", err)
		}

		// 过滤掉当前书籍
		var filteredBooks []*bookstore.BookDetail
		for _, book := range similarBooks {
			if book.ID != bookID {
				filteredBooks = append(filteredBooks, book)
			}
		}

		return filteredBooks, nil
	}

	// 如果没有标签，基于作者推荐
	authorBooks, _, err := s.GetBookDetailsByAuthor(ctx, currentBook.Author, 1, limit)
	return authorBooks, err
}

// GetPopularBookDetails 获取热门书籍详情
func (s *BookDetailServiceImpl) GetPopularBookDetails(ctx context.Context, limit int) ([]*bookstore.BookDetail, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 使用已完成状态的书籍作为热门书籍
	bookDetails, err := s.bookDetailRepo.GetByStatus(ctx, bookstore.BookStatusCompleted, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular books: %w", err)
	}

	return bookDetails, nil
}

// GetLatestBookDetails 获取最新书籍详情
func (s *BookDetailServiceImpl) GetLatestBookDetails(ctx context.Context, limit int) ([]*bookstore.BookDetail, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 使用正在连载状态的书籍作为最新书籍
	bookDetails, err := s.bookDetailRepo.GetByStatus(ctx, bookstore.BookStatusOngoing, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest books: %w", err)
	}

	return bookDetails, nil
}

// IncrementViewCount 增加浏览计数
func (s *BookDetailServiceImpl) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.IncrementViewCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// IncrementLikeCount 增加点赞计数
func (s *BookDetailServiceImpl) IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.IncrementLikeCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment like count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// DecrementLikeCount 减少点赞计数
func (s *BookDetailServiceImpl) DecrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.DecrementLikeCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to decrement like count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// IncrementCommentCount 增加评论计数
func (s *BookDetailServiceImpl) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.IncrementCommentCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment comment count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// DecrementCommentCount 减少评论计数
func (s *BookDetailServiceImpl) DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.DecrementCommentCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to decrement comment count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// IncrementShareCount 增加分享计数
func (s *BookDetailServiceImpl) IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error {
	if err := s.bookDetailRepo.IncrementShareCount(ctx, bookID); err != nil {
		return fmt.Errorf("failed to increment share count: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// UpdateRating 更新评分
func (s *BookDetailServiceImpl) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64, ratingCount int64) error {
	if err := s.bookDetailRepo.UpdateRating(ctx, bookID, rating, ratingCount); err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// UpdateLastChapter 更新最新章节信息
func (s *BookDetailServiceImpl) UpdateLastChapter(ctx context.Context, bookID primitive.ObjectID, chapterTitle string) error {
	if err := s.bookDetailRepo.UpdateLastChapter(ctx, bookID, chapterTitle); err != nil {
		return fmt.Errorf("failed to update last chapter: %w", err)
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
	}

	return nil
}

// BatchUpdateBookDetailStatus 批量更新书籍状态
func (s *BookDetailServiceImpl) BatchUpdateBookDetailStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}

	if err := s.bookDetailRepo.BatchUpdateStatus(ctx, bookIDs, status); err != nil {
		return fmt.Errorf("failed to batch update book status: %w", err)
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
		}
	}

	return nil
}

// BatchUpdateBookDetailCategories 批量更新书籍分类
func (s *BookDetailServiceImpl) BatchUpdateBookDetailCategories(ctx context.Context, bookIDs []primitive.ObjectID, categories []string) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}
	if len(categories) == 0 {
		return errors.New("categories cannot be empty")
	}

	if err := s.bookDetailRepo.BatchUpdateCategories(ctx, bookIDs, categories); err != nil {
		return fmt.Errorf("failed to batch update book categories: %w", err)
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
		}
	}

	return nil
}

// BatchUpdateBookDetailTags 批量更新书籍标签
func (s *BookDetailServiceImpl) BatchUpdateBookDetailTags(ctx context.Context, bookIDs []primitive.ObjectID, tags []string) error {
	if len(bookIDs) == 0 {
		return errors.New("book IDs cannot be empty")
	}
	if len(tags) == 0 {
		return errors.New("tags cannot be empty")
	}

	if err := s.bookDetailRepo.BatchUpdateTags(ctx, bookIDs, tags); err != nil {
		return fmt.Errorf("failed to batch update book tags: %w", err)
	}

	// 清除相关缓存
	for _, bookID := range bookIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookDetailCache(ctx, bookID.Hex())
		}
	}

	return nil
}

// invalidateRelatedCache 清除相关缓存
func (s *BookDetailServiceImpl) invalidateRelatedCache(ctx context.Context, bookDetail *bookstore.BookDetail) {
	if s.cacheService == nil {
		return
	}

	// 清除书籍详情缓存
	s.cacheService.InvalidateBookDetailCache(ctx, bookDetail.ID.Hex())

	// 清除分类相关缓存
	for _, category := range bookDetail.Categories {
		s.cacheService.InvalidateCategoryCache(ctx, category)
	}

	// 清除作者相关缓存
	s.cacheService.InvalidateAuthorCache(ctx, bookDetail.Author)

	// 清除首页缓存
	s.cacheService.InvalidateHomepageCache(ctx)
}

// API 兼容方法别名实现

// GetBooksByTitle 根据标题搜索书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetBooksByTitle(ctx context.Context, title string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	return s.SearchBookDetails(ctx, title, page, pageSize)
}

// GetBooksByAuthor 根据作者获取书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetBooksByAuthor(ctx context.Context, author string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	return s.GetBookDetailsByAuthor(ctx, author, page, pageSize)
}

// GetBooksByCategory 根据分类获取书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetBooksByCategory(ctx context.Context, category string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	return s.GetBookDetailsByCategory(ctx, category, page, pageSize)
}

// GetBooksByStatus 根据状态获取书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetBooksByStatus(ctx context.Context, status string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	// 将字符串状态转换为 BookStatus 类型
	var bookStatus bookstore.BookStatus
	switch status {
	case "serializing", "ongoing":
		bookStatus = bookstore.BookStatusOngoing
	case "completed":
		bookStatus = bookstore.BookStatusCompleted
	case "paused":
		bookStatus = bookstore.BookStatusPaused
	default:
		return nil, 0, fmt.Errorf("无效的书籍状态: %s", status)
	}
	return s.GetBookDetailsByStatus(ctx, bookStatus, page, pageSize)
}

// GetBooksByTags 根据标签获取书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetBooksByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	return s.GetBookDetailsByTags(ctx, tags, page, pageSize)
}

// SearchBooks 搜索书籍（API 兼容方法）
func (s *BookDetailServiceImpl) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookDetail, int64, error) {
	return s.SearchBookDetails(ctx, keyword, page, pageSize)
}

// GetRecommendedBooks 获取推荐书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetRecommendedBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error) {
	// 返回热门书籍作为推荐
	return s.GetPopularBookDetails(ctx, limit)
}

// GetSimilarBooks 获取相似书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetSimilarBooks(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore.BookDetail, error) {
	return s.GetSimilarBookDetails(ctx, bookID, limit)
}

// GetPopularBooks 获取热门书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetPopularBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error) {
	return s.GetPopularBookDetails(ctx, limit)
}

// GetLatestBooks 获取最新书籍（API 兼容方法）
func (s *BookDetailServiceImpl) GetLatestBooks(ctx context.Context, limit int) ([]*bookstore.BookDetail, error) {
	return s.GetLatestBookDetails(ctx, limit)
}

// CountBooksByCategory 根据分类统计书籍数量（API 兼容方法）
func (s *BookDetailServiceImpl) CountBooksByCategory(ctx context.Context, category string) (int64, error) {
	return s.GetBookDetailCountByCategory(ctx, category)
}
