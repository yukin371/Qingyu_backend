package reader

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoremodels "Qingyu_backend/models/bookstore"
	"Qingyu_backend/service/bookstore"
)

var (
	// ErrChapterNotFound 章节不存在
	ErrChapterNotFound = errors.New("chapter not found")
	// ErrChapterNotPublished 章节未发布
	ErrChapterNotPublished = errors.New("chapter is not published")
	// ErrAccessDenied 无权访问
	ErrAccessDenied = errors.New("access denied to this chapter")
)

// ChapterService 阅读器章节服务接口
type ChapterService interface {
	// GetChapterContent 获取章节内容（带权限检查、进度保存）
	GetChapterContent(ctx context.Context, userID, bookID, chapterID string) (*ChapterContentResponse, error)
	GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*ChapterContentResponse, error)

	// GetNextChapter 获取下一章
	GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error)
	GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error)

	// GetChapterList 获取章节目录
	GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*ChapterListResponse, error)

	// GetChapterInfo 获取章节信息（不含内容）
	GetChapterInfo(ctx context.Context, userID, chapterID string) (*ChapterInfo, error)
}

// ChapterServiceImpl 章节服务实现
type ChapterServiceImpl struct {
	chapterService    bookstore.ChapterService
	readerService     *ReaderService
	vipService        VIPPermissionService
}

// ChapterContentResponse 章节内容响应
type ChapterContentResponse struct {
	ChapterID     string    `json:"chapterId"`
	BookID        string    `json:"bookId"`
	Title         string    `json:"title"`
	ChapterNum    int       `json:"chapterNum"`
	Content       string    `json:"content"`
	WordCount     int       `json:"wordCount"`
	HasNext       bool      `json:"hasNext"`
	HasPrevious   bool      `json:"hasPrevious"`
	Progress      float64   `json:"progress"`
	ReadingTime   int64     `json:"readingTime"`
	LastReadAt    time.Time `json:"lastReadAt"`
	CanAccess     bool      `json:"canAccess"`
	AccessReason  string    `json:"accessReason,omitempty"`
}

// ChapterInfo 章节信息（不含内容）
type ChapterInfo struct {
	ChapterID    string    `json:"chapterId"`
	BookID       string    `json:"bookId"`
	Title        string    `json:"title"`
	ChapterNum   int       `json:"chapterNum"`
	WordCount    int       `json:"wordCount"`
	IsFree       bool      `json:"isFree"`
	Price        float64   `json:"price"`
	PublishTime  time.Time `json:"publishTime"`
	Progress     float64   `json:"progress"`
	IsRead       bool      `json:"isRead"`
	CanAccess    bool      `json:"canAccess"`
	AccessReason string    `json:"accessReason,omitempty"`
}

// ChapterListResponse 章节列表响应
type ChapterListResponse struct {
	Chapters   []*ChapterInfo `json:"chapters"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	BookID     string         `json:"bookId"`
	BookTitle  string         `json:"bookTitle"`
	Author     string         `json:"author"`
	TotalWords int64          `json:"totalWords"`
}

// NewChapterService 创建章节服务
func NewChapterService(
	chapterService bookstore.ChapterService,
	readerService *ReaderService,
	vipService VIPPermissionService,
) ChapterService {
	return &ChapterServiceImpl{
		chapterService: chapterService,
		readerService:  readerService,
		vipService:     vipService,
	}
}

// GetChapterContent 获取章节内容（阅读器专用）
func (s *ChapterServiceImpl) GetChapterContent(ctx context.Context, userID, bookID, chapterID string) (*ChapterContentResponse, error) {
	// 解析ID
	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}

	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	// 获取章节元数据
	chapter, err := s.chapterService.GetChapterByID(ctx, chapterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, ErrChapterNotFound
	}

	// 检查发布状态
	if !chapter.IsPublished() {
		return nil, ErrChapterNotPublished
	}

	// 权限检查
	userOID := primitive.NilObjectID
	if userID != "" {
		userOID, _ = primitive.ObjectIDFromHex(userID)
	}

	canAccess, accessReason := s.checkChapterAccess(ctx, userOID, chapter, bookOID)
	if !canAccess {
		// 返回章节信息但不含内容
		return &ChapterContentResponse{
			ChapterID:    chapterID,
			BookID:       bookID,
			Title:        chapter.Title,
			ChapterNum:   chapter.ChapterNum,
			Content:      "",
			WordCount:    chapter.WordCount,
			CanAccess:    false,
			AccessReason: accessReason,
		}, ErrAccessDenied
	}

	// 获取章节内容
	content, err := s.chapterService.GetChapterContent(ctx, chapterOID, userOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter content: %w", err)
	}

	// 获取导航信息
	hasNext, hasPrevious := s.getNavigationInfo(ctx, chapter)

	// 获取或创建阅读进度
	var progress float64
	var readingTime int64
	var lastReadAt time.Time

	if userID != "" {
		readingProgress, err := s.readerService.GetReadingProgress(ctx, userID, bookID)
		if err == nil && readingProgress != nil {
			progress = readingProgress.Progress
			readingTime = readingProgress.ReadingTime
			lastReadAt = readingProgress.LastReadAt
		}

		// 更新阅读进度到当前章节
		_ = s.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, 0)
	}

	return &ChapterContentResponse{
		ChapterID:    chapterID,
		BookID:       bookID,
		Title:        chapter.Title,
		ChapterNum:   chapter.ChapterNum,
		Content:      content,
		WordCount:    chapter.WordCount,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		Progress:     progress,
		ReadingTime:  readingTime,
		LastReadAt:   lastReadAt,
		CanAccess:    true,
		AccessReason: "",
	}, nil
}

// GetChapterByNumber 根据章节号获取内容
func (s *ChapterServiceImpl) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*ChapterContentResponse, error) {
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	// 获取章节
	chapter, err := s.chapterService.GetChapterByBookIDAndNum(ctx, bookOID, chapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, ErrChapterNotFound
	}

	return s.GetChapterContent(ctx, userID, bookID, chapter.ID.Hex())
}

// GetNextChapter 获取下一章信息
func (s *ChapterServiceImpl) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error) {
	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}

	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	// 获取当前章节
	currentChapter, err := s.chapterService.GetChapterByID(ctx, chapterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current chapter: %w", err)
	}
	if currentChapter == nil {
		return nil, ErrChapterNotFound
	}

	// 获取下一章
	nextChapter, err := s.chapterService.GetNextChapter(ctx, bookOID, currentChapter.ChapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get next chapter: %w", err)
	}
	if nextChapter == nil {
		return nil, nil // 没有下一章
	}

	return s.buildChapterInfo(ctx, userID, nextChapter)
}

// GetPreviousChapter 获取上一章信息
func (s *ChapterServiceImpl) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error) {
	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}

	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	// 获取当前章节
	currentChapter, err := s.chapterService.GetChapterByID(ctx, chapterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current chapter: %w", err)
	}
	if currentChapter == nil {
		return nil, ErrChapterNotFound
	}

	// 获取上一章
	prevChapter, err := s.chapterService.GetPreviousChapter(ctx, bookOID, currentChapter.ChapterNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous chapter: %w", err)
	}
	if prevChapter == nil {
		return nil, nil // 没有上一章
	}

	return s.buildChapterInfo(ctx, userID, prevChapter)
}

// GetChapterList 获取章节目录
func (s *ChapterServiceImpl) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*ChapterListResponse, error) {
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	// 获取章节列表
	chapters, total, err := s.chapterService.GetChaptersByBookID(ctx, bookOID, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter list: %w", err)
	}

	// 构建响应
	chapterInfos := make([]*ChapterInfo, 0, len(chapters))
	for _, chapter := range chapters {
		info, _ := s.buildChapterInfo(ctx, userID, chapter)
		chapterInfos = append(chapterInfos, info)
	}

	// 获取书籍统计信息
	totalWords, _ := s.chapterService.GetTotalWordCountByBookID(ctx, bookOID)

	return &ChapterListResponse{
		Chapters:   chapterInfos,
		Total:      total,
		Page:       page,
		Size:       size,
		BookID:     bookID,
		BookTitle:  "", // TODO: 从书籍信息获取
		Author:     "",
		TotalWords: totalWords,
	}, nil
}

// GetChapterInfo 获取章节信息（不含内容）
func (s *ChapterServiceImpl) GetChapterInfo(ctx context.Context, userID, chapterID string) (*ChapterInfo, error) {
	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}

	chapter, err := s.chapterService.GetChapterByID(ctx, chapterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, ErrChapterNotFound
	}

	return s.buildChapterInfo(ctx, userID, chapter)
}

// checkChapterAccess 检查章节访问权限
func (s *ChapterServiceImpl) checkChapterAccess(ctx context.Context, userID primitive.ObjectID, chapter *bookstoremodels.Chapter, bookID primitive.ObjectID) (bool, string) {
	// 免费章节，所有人可访问
	if chapter.IsFree {
		return true, ""
	}

	// 未登录用户无法访问付费章节
	if userID.IsZero() {
		return false, "需要登录后阅读付费章节"
	}

	// 检查VIP权限
	if s.vipService != nil {
		// 使用 CheckVIPAccess 检查VIP权限
		hasAccess, err := s.vipService.CheckVIPAccess(ctx, userID.Hex(), chapter.ID.Hex(), !chapter.IsFree)
		if err == nil && hasAccess {
			return true, ""
		}
		// 如果VIP检查失败，继续检查是否已购买该章节
		purchased, err := s.vipService.CheckChapterPurchased(ctx, userID.Hex(), chapter.ID.Hex())
		if err == nil && purchased {
			return true, ""
		}
	}

	return false, "该章节需要购买或开通VIP后阅读"
}

// getNavigationInfo 获取导航信息
func (s *ChapterServiceImpl) getNavigationInfo(ctx context.Context, chapter *bookstoremodels.Chapter) (hasNext, hasPrevious bool) {
	// 检查是否有下一章
	nextChapter, err := s.chapterService.GetNextChapter(ctx, chapter.BookID, chapter.ChapterNum)
	hasNext = (err == nil && nextChapter != nil)

	// 检查是否有上一章
	prevChapter, err := s.chapterService.GetPreviousChapter(ctx, chapter.BookID, chapter.ChapterNum)
	hasPrevious = (err == nil && prevChapter != nil)

	return hasNext, hasPrevious
}

// buildChapterInfo 构建章节信息
func (s *ChapterServiceImpl) buildChapterInfo(ctx context.Context, userID string, chapter *bookstoremodels.Chapter) (*ChapterInfo, error) {
	userOID := primitive.NilObjectID
	if userID != "" {
		userOID, _ = primitive.ObjectIDFromHex(userID)
	}

	canAccess, accessReason := s.checkChapterAccess(ctx, userOID, chapter, chapter.BookID)

	// 获取阅读进度
	var progress float64
	isRead := false

	if userID != "" {
		readingProgress, err := s.readerService.GetReadingProgress(ctx, userID, chapter.BookID.Hex())
		if err == nil && readingProgress != nil {
			progress = readingProgress.Progress
			isRead = (readingProgress.ChapterID == chapter.ID.Hex())
		}
	}

	return &ChapterInfo{
		ChapterID:    chapter.ID.Hex(),
		BookID:       chapter.BookID.Hex(),
		Title:        chapter.Title,
		ChapterNum:   chapter.ChapterNum,
		WordCount:    chapter.WordCount,
		IsFree:       chapter.IsFree,
		Price:        chapter.Price,
		PublishTime:  chapter.PublishTime,
		Progress:     progress,
		IsRead:       isRead,
		CanAccess:    canAccess,
		AccessReason: accessReason,
	}, nil
}
