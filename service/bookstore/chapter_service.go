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

// ChapterService 章节服务接口
type ChapterService interface {
	// 章节基础操作
	CreateChapter(ctx context.Context, chapter *bookstore.Chapter) error
	GetChapterByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Chapter, error)
	UpdateChapter(ctx context.Context, chapter *bookstore.Chapter) error
	DeleteChapter(ctx context.Context, id primitive.ObjectID) error

	// 章节查询
	GetChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error)
	GetChapterByBookIDAndNum(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetChaptersByTitle(ctx context.Context, title string, page, pageSize int) ([]*bookstore.Chapter, int64, error)
	GetFreeChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error)
	GetPaidChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error)
	GetPublishedChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error)

	// 章节导航
	GetPreviousChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetNextChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetFirstChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error)
	GetLastChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error)

	// 章节统计
	GetChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetFreeChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetPaidChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetTotalWordCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetChapterStats(ctx context.Context, bookID primitive.ObjectID) (map[string]interface{}, error)

	// 章节内容管理
	GetChapterContent(ctx context.Context, chapterID primitive.ObjectID, userID primitive.ObjectID) (string, error)
	UpdateChapterContent(ctx context.Context, chapterID primitive.ObjectID, content string) error
	PublishChapter(ctx context.Context, chapterID primitive.ObjectID) error
	UnpublishChapter(ctx context.Context, chapterID primitive.ObjectID) error

	// 章节批量操作
	BatchUpdateChapterPrice(ctx context.Context, chapterIDs []primitive.ObjectID, price float64) error
	BatchPublishChapters(ctx context.Context, chapterIDs []primitive.ObjectID) error
	BatchDeleteChapters(ctx context.Context, chapterIDs []primitive.ObjectID) error
	BatchDeleteChaptersByBookID(ctx context.Context, bookID primitive.ObjectID) error

	// 章节搜索
	SearchChapters(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Chapter, int64, error)
}

// ChapterServiceImpl 章节服务实现
type ChapterServiceImpl struct {
	chapterRepo  BookstoreRepo.ChapterRepository
	cacheService CacheService
}

// NewChapterService 创建章节服务实例
func NewChapterService(chapterRepo BookstoreRepo.ChapterRepository, cacheService CacheService) ChapterService {
	return &ChapterServiceImpl{
		chapterRepo:  chapterRepo,
		cacheService: cacheService,
	}
}

// CreateChapter 创建章节
func (s *ChapterServiceImpl) CreateChapter(ctx context.Context, chapter *bookstore.Chapter) error {
	if chapter == nil {
		return errors.New("chapter cannot be nil")
	}

	// 验证必填字段
	if chapter.BookID.IsZero() {
		return errors.New("book ID is required")
	}
	if chapter.Title == "" {
		return errors.New("chapter title is required")
	}
	if chapter.ChapterNum <= 0 {
		return errors.New("chapter number must be positive")
	}

	// 检查章节号是否已存在
	existingChapter, err := s.chapterRepo.GetByBookIDAndChapterNum(ctx, chapter.BookID, chapter.ChapterNum)
	if err != nil {
		return fmt.Errorf("failed to check existing chapter: %w", err)
	}
	if existingChapter != nil {
		return errors.New("chapter with this number already exists for this book")
	}

	// 创建章节
	if err := s.chapterRepo.Create(ctx, chapter); err != nil {
		return fmt.Errorf("failed to create chapter: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// GetChapterByID 根据ID获取章节
func (s *ChapterServiceImpl) GetChapterByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Chapter, error) {
	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedChapter, err := s.cacheService.GetChapter(ctx, id.Hex()); err == nil && cachedChapter != nil {
			return cachedChapter, nil
		}
	}

	// 从数据库获取
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, errors.New("chapter not found")
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetChapter(ctx, id.Hex(), chapter, 30*time.Minute)
	}

	return chapter, nil
}

// UpdateChapter 更新章节
func (s *ChapterServiceImpl) UpdateChapter(ctx context.Context, chapter *bookstore.Chapter) error {
	if chapter == nil {
		return errors.New("chapter cannot be nil")
	}

	// 验证必填字段
	if chapter.BookID.IsZero() {
		return errors.New("book ID is required")
	}
	if chapter.Title == "" {
		return errors.New("chapter title is required")
	}
	if chapter.ChapterNum <= 0 {
		return errors.New("chapter number must be positive")
	}

	// 准备更新数据
	updates := map[string]interface{}{
		"book_id":      chapter.BookID,
		"title":        chapter.Title,
		"content":      chapter.Content,
		"chapter_num":  chapter.ChapterNum,
		"word_count":   chapter.WordCount,
		"is_free":      chapter.IsFree,
		"price":        chapter.Price,
		"publish_time": chapter.PublishTime,
	}

	// 更新章节
	if err := s.chapterRepo.Update(ctx, chapter.ID, updates); err != nil {
		return fmt.Errorf("failed to update chapter: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// DeleteChapter 删除章节
func (s *ChapterServiceImpl) DeleteChapter(ctx context.Context, id primitive.ObjectID) error {
	// 先获取章节信息用于清除缓存
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get chapter for deletion: %w", err)
	}
	if chapter == nil {
		return errors.New("chapter not found")
	}

	// 删除章节
	if err := s.chapterRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete chapter: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// GetChaptersByBookID 根据书籍ID获取章节列表
func (s *ChapterServiceImpl) GetChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取章节列表
	chapters, err := s.chapterRepo.GetByBookID(ctx, bookID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chapters by book ID: %w", err)
	}

	// 获取总数
	total, err := s.chapterRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chapters by book ID: %w", err)
	}

	return chapters, total, nil
}

// GetChapterByBookIDAndNum 根据书籍ID和章节号获取章节
func (s *ChapterServiceImpl) GetChapterByBookIDAndNum(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}
	if chapterNum <= 0 {
		return nil, errors.New("chapter number must be positive")
	}

	chapter, err := s.chapterRepo.GetByBookIDAndChapterNum(ctx, bookID, chapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter by book ID and number: %w", err)
	}

	return chapter, nil
}

// GetChaptersByTitle 根据标题搜索章节
func (s *ChapterServiceImpl) GetChaptersByTitle(ctx context.Context, title string, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
	if title == "" {
		return nil, 0, errors.New("title cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取章节列表
	chapters, err := s.chapterRepo.GetByTitle(ctx, title, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chapters by title: %w", err)
	}

	// 这里简化处理，实际应该有专门的计数方法
	total := int64(len(chapters))
	if len(chapters) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return chapters, total, nil
}

// GetFreeChaptersByBookID 获取免费章节列表
func (s *ChapterServiceImpl) GetFreeChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取免费章节列表
	chapters, err := s.chapterRepo.GetFreeChapters(ctx, bookID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get free chapters: %w", err)
	}

	// 获取总数
	total, err := s.chapterRepo.CountFreeChapters(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count free chapters: %w", err)
	}

	return chapters, total, nil
}

// GetPaidChaptersByBookID 获取付费章节列表
func (s *ChapterServiceImpl) GetPaidChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取付费章节列表
	chapters, err := s.chapterRepo.GetPaidChapters(ctx, bookID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paid chapters: %w", err)
	}

	// 获取总数
	total, err := s.chapterRepo.CountPaidChapters(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count paid chapters: %w", err)
	}

	return chapters, total, nil
}

// GetPublishedChaptersByBookID 获取已发布章节列表
func (s *ChapterServiceImpl) GetPublishedChaptersByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取已发布章节列表
	chapters, err := s.chapterRepo.GetPublishedChapters(ctx, bookID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get published chapters: %w", err)
	}

	// 这里简化处理，实际应该有专门的计数方法
	total := int64(len(chapters))
	if len(chapters) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return chapters, total, nil
}

// GetPreviousChapter 获取上一章节
func (s *ChapterServiceImpl) GetPreviousChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}
	if chapterNum <= 1 {
		return nil, errors.New("no previous chapter available")
	}

	chapter, err := s.chapterRepo.GetPreviousChapter(ctx, bookID, chapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous chapter: %w", err)
	}

	return chapter, nil
}

// GetNextChapter 获取下一章节
func (s *ChapterServiceImpl) GetNextChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}
	if chapterNum <= 0 {
		return nil, errors.New("invalid chapter number")
	}

	chapter, err := s.chapterRepo.GetNextChapter(ctx, bookID, chapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get next chapter: %w", err)
	}

	return chapter, nil
}

// GetFirstChapter 获取第一章节
func (s *ChapterServiceImpl) GetFirstChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	chapter, err := s.chapterRepo.GetFirstChapter(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get first chapter: %w", err)
	}

	return chapter, nil
}

// GetLastChapter 获取最后章节
func (s *ChapterServiceImpl) GetLastChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	chapter, err := s.chapterRepo.GetLastChapter(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last chapter: %w", err)
	}

	return chapter, nil
}

// GetChapterCountByBookID 获取书籍章节总数
func (s *ChapterServiceImpl) GetChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	count, err := s.chapterRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to count chapters by book ID: %w", err)
	}

	return count, nil
}

// GetFreeChapterCountByBookID 获取免费章节数量
func (s *ChapterServiceImpl) GetFreeChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	count, err := s.chapterRepo.CountFreeChapters(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to count free chapters: %w", err)
	}

	return count, nil
}

// GetPaidChapterCountByBookID 获取付费章节数量
func (s *ChapterServiceImpl) GetPaidChapterCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	count, err := s.chapterRepo.CountPaidChapters(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to count paid chapters: %w", err)
	}

	return count, nil
}

// GetTotalWordCountByBookID 获取书籍总字数
func (s *ChapterServiceImpl) GetTotalWordCountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	wordCount, err := s.chapterRepo.GetTotalWordCount(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total word count: %w", err)
	}

	return wordCount, nil
}

// GetChapterStats 获取章节统计信息
func (s *ChapterServiceImpl) GetChapterStats(ctx context.Context, bookID primitive.ObjectID) (map[string]interface{}, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	stats := make(map[string]interface{})

	// 总章节数
	totalCount, err := s.chapterRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total chapter count: %w", err)
	}
	stats["total_chapters"] = totalCount

	// 免费章节数
	freeCount, err := s.chapterRepo.CountFreeChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get free chapter count: %w", err)
	}
	stats["free_chapters"] = freeCount

	// 付费章节数
	paidCount, err := s.chapterRepo.CountPaidChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paid chapter count: %w", err)
	}
	stats["paid_chapters"] = paidCount

	// 总字数
	totalWordCount, err := s.chapterRepo.GetTotalWordCount(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total word count: %w", err)
	}
	stats["total_word_count"] = totalWordCount

	return stats, nil
}

// GetChapterContent 获取章节内容（考虑权限）
func (s *ChapterServiceImpl) GetChapterContent(ctx context.Context, chapterID primitive.ObjectID, userID primitive.ObjectID) (string, error) {
	if chapterID.IsZero() {
		return "", errors.New("chapter ID cannot be empty")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return "", fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return "", errors.New("chapter not found")
	}

	// 检查是否已发布
	if !chapter.IsPublished() {
		return "", errors.New("chapter is not published")
	}

	// 如果是免费章节，直接返回内容
	if chapter.IsFree {
		return chapter.Content, nil
	}

	// 付费章节需要检查用户权限（这里简化处理）
	if userID.IsZero() {
		return "", errors.New("user authentication required for paid content")
	}

	// TODO: 实际应该检查用户是否已购买该章节
	// 这里简化处理，假设已购买
	return chapter.Content, nil
}

// UpdateChapterContent 更新章节内容
func (s *ChapterServiceImpl) UpdateChapterContent(ctx context.Context, chapterID primitive.ObjectID, content string) error {
	if chapterID.IsZero() {
		return errors.New("chapter ID cannot be empty")
	}
	if content == "" {
		return errors.New("content cannot be empty")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return errors.New("chapter not found")
	}

	// 更新内容和字数
	updates := map[string]interface{}{
		"content":    content,
		"word_count": int64(len([]rune(content))), // 使用rune计算字符数
	}

	// 保存更新
	if err := s.chapterRepo.Update(ctx, chapterID, updates); err != nil {
		return fmt.Errorf("failed to update chapter content: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// PublishChapter 发布章节
func (s *ChapterServiceImpl) PublishChapter(ctx context.Context, chapterID primitive.ObjectID) error {
	if chapterID.IsZero() {
		return errors.New("chapter ID cannot be empty")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return errors.New("chapter not found")
	}

	// 设置发布时间
	now := time.Now()
	updates := map[string]interface{}{
		"is_published": true,
		"publish_time": &now,
	}

	// 保存更新
	if err := s.chapterRepo.Update(ctx, chapterID, updates); err != nil {
		return fmt.Errorf("failed to publish chapter: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// UnpublishChapter 取消发布章节
func (s *ChapterServiceImpl) UnpublishChapter(ctx context.Context, chapterID primitive.ObjectID) error {
	if chapterID.IsZero() {
		return errors.New("chapter ID cannot be empty")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return errors.New("chapter not found")
	}

	// 清除发布时间
	updates := map[string]interface{}{
		"is_published": false,
		"publish_time": nil,
	}

	// 保存更新
	if err := s.chapterRepo.Update(ctx, chapterID, updates); err != nil {
		return fmt.Errorf("failed to unpublish chapter: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, chapter)

	return nil
}

// BatchUpdateChapterPrice 批量更新章节价格
func (s *ChapterServiceImpl) BatchUpdateChapterPrice(ctx context.Context, chapterIDs []primitive.ObjectID, price float64) error {
	if len(chapterIDs) == 0 {
		return errors.New("chapter IDs cannot be empty")
	}
	if price < 0 {
		return errors.New("price cannot be negative")
	}

	if err := s.chapterRepo.BatchUpdatePrice(ctx, chapterIDs, price); err != nil {
		return fmt.Errorf("failed to batch update chapter price: %w", err)
	}

	// 清除相关缓存
	for _, chapterID := range chapterIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateChapterCache(ctx, chapterID.Hex())
		}
	}

	return nil
}

// BatchPublishChapters 批量发布章节
func (s *ChapterServiceImpl) BatchPublishChapters(ctx context.Context, chapterIDs []primitive.ObjectID) error {
	if len(chapterIDs) == 0 {
		return errors.New("chapter IDs cannot be empty")
	}

	now := time.Now()
	if err := s.chapterRepo.BatchUpdatePublishTime(ctx, chapterIDs, now); err != nil {
		return fmt.Errorf("failed to batch publish chapters: %w", err)
	}

	// 清除相关缓存
	for _, chapterID := range chapterIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateChapterCache(ctx, chapterID.Hex())
		}
	}

	return nil
}

// BatchDeleteChapters 批量删除章节
func (s *ChapterServiceImpl) BatchDeleteChapters(ctx context.Context, chapterIDs []primitive.ObjectID) error {
	if len(chapterIDs) == 0 {
		return errors.New("chapter IDs cannot be empty")
	}

	if err := s.chapterRepo.BatchDelete(ctx, chapterIDs); err != nil {
		return fmt.Errorf("failed to batch delete chapters: %w", err)
	}

	// 清除相关缓存
	for _, chapterID := range chapterIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateChapterCache(ctx, chapterID.Hex())
		}
	}

	return nil
}

// BatchDeleteChaptersByBookID 批量删除书籍的所有章节
func (s *ChapterServiceImpl) BatchDeleteChaptersByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	if err := s.chapterRepo.DeleteByBookID(ctx, bookID); err != nil {
		return fmt.Errorf("failed to batch delete chapters by book ID: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookChaptersCache(ctx, bookID.Hex())
	}

	return nil
}

// SearchChapters 搜索章节
func (s *ChapterServiceImpl) SearchChapters(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Chapter, int64, error) {
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

	// 搜索章节
	chapters, err := s.chapterRepo.Search(ctx, keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search chapters: %w", err)
	}

	// TODO: 这里简化处理，实际应该有专门的搜索计数方法
	total := int64(len(chapters))
	if len(chapters) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return chapters, total, nil
}

// invalidateRelatedCache 清除相关缓存
func (s *ChapterServiceImpl) invalidateRelatedCache(ctx context.Context, chapter *bookstore.Chapter) {
	if s.cacheService == nil {
		return
	}

	// 清除章节缓存
	s.cacheService.InvalidateChapterCache(ctx, chapter.ID.Hex())

	// 清除书籍章节列表缓存
	s.cacheService.InvalidateBookChaptersCache(ctx, chapter.BookID.Hex())

	// 清除书籍详情缓存（因为章节数量可能变化）
	s.cacheService.InvalidateBookDetailCache(ctx, chapter.BookID.Hex())
}
