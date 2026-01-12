package reader

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	readermodels "Qingyu_backend/models/reader"
	readerrepo "Qingyu_backend/repository/interfaces/reader"
)

var (
	// ErrBookmarkNotFound 书签不存在
	ErrBookmarkNotFound = errors.New("bookmark not found")
	// ErrBookmarkAlreadyExists 书签已存在
	ErrBookmarkAlreadyExists = errors.New("bookmark already exists")
)

// BookmarkService 书签服务接口
type BookmarkService interface {
	// CreateBookmark 创建书签
	CreateBookmark(ctx context.Context, bookmark *readermodels.Bookmark) error

	// GetBookmark 获取书签详情
	GetBookmark(ctx context.Context, bookmarkID string) (*readermodels.Bookmark, error)

	// GetUserBookmarks 获取用户书签列表
	GetUserBookmarks(ctx context.Context, userID string, filter *readermodels.BookmarkFilter, page, size int) (*BookmarkListResponse, error)

	// GetBookBookmarks 获取书籍书签
	GetBookBookmarks(ctx context.Context, userID, bookID string, page, size int) (*BookmarkListResponse, error)

	// UpdateBookmark 更新书签
	UpdateBookmark(ctx context.Context, bookmarkID string, bookmark *readermodels.Bookmark) error

	// DeleteBookmark 删除书签
	DeleteBookmark(ctx context.Context, bookmarkID string) error

	// ExportBookmarks 导出书签
	ExportBookmarks(ctx context.Context, userID, format string) ([]byte, string, error)

	// GetBookmarkStats 获取书签统计
	GetBookmarkStats(ctx context.Context, userID string) (*readermodels.BookmarkStats, error)

	// SearchBookmarks 搜索书签
	SearchBookmarks(ctx context.Context, userID, keyword string, page, size int) (*BookmarkListResponse, error)
}

// BookmarkServiceImpl 书签服务实现
type BookmarkServiceImpl struct {
	bookmarkRepo   readerrepo.BookmarkRepository
	chapterService ChapterService
}

// BookmarkListResponse 书签列表响应
type BookmarkListResponse struct {
	Bookmarks []*readermodels.Bookmark `json:"bookmarks"`
	Total     int64                    `json:"total"`
	Page      int                      `json:"page"`
	Size      int                      `json:"size"`
}

// NewBookmarkService 创建书签服务
func NewBookmarkService(
	bookmarkRepo readerrepo.BookmarkRepository,
	chapterService ChapterService,
) BookmarkService {
	return &BookmarkServiceImpl{
		bookmarkRepo:   bookmarkRepo,
		chapterService: chapterService,
	}
}

// CreateBookmark 创建书签
func (s *BookmarkServiceImpl) CreateBookmark(ctx context.Context, bookmark *readermodels.Bookmark) error {
	// 验证必填字段
	if bookmark.UserID.IsZero() {
		return errors.New("user ID is required")
	}
	if bookmark.BookID.IsZero() {
		return errors.New("book ID is required")
	}
	if bookmark.ChapterID.IsZero() {
		return errors.New("chapter ID is required")
	}

	// 设置默认颜色
	if bookmark.Color == "" {
		bookmark.Color = "yellow"
	}

	// 检查是否已存在同一位置的书签
	exists, err := s.bookmarkRepo.Exists(ctx, bookmark.UserID, bookmark.ChapterID, bookmark.Position)
	if err != nil {
		return fmt.Errorf("failed to check bookmark existence: %w", err)
	}
	if exists {
		return ErrBookmarkAlreadyExists
	}

	return s.bookmarkRepo.Create(ctx, bookmark)
}

// GetBookmark 获取书签详情
func (s *BookmarkServiceImpl) GetBookmark(ctx context.Context, bookmarkID string) (*readermodels.Bookmark, error) {
	id, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("invalid bookmark ID: %w", err)
	}

	bookmark, err := s.bookmarkRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookmarkNotFound
	}

	return bookmark, nil
}

// GetUserBookmarks 获取用户书签列表
func (s *BookmarkServiceImpl) GetUserBookmarks(ctx context.Context, userID string, filter *readermodels.BookmarkFilter, page, size int) (*BookmarkListResponse, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	bookmarks, total, err := s.bookmarkRepo.GetByUser(ctx, userOID, filter, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks: %w", err)
	}

	return &BookmarkListResponse{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		Size:      size,
	}, nil
}

// GetBookBookmarks 获取书籍书签
func (s *BookmarkServiceImpl) GetBookBookmarks(ctx context.Context, userID, bookID string, page, size int) (*BookmarkListResponse, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	bookmarks, total, err := s.bookmarkRepo.GetByBook(ctx, userOID, bookOID, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks: %w", err)
	}

	return &BookmarkListResponse{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		Size:      size,
	}, nil
}

// UpdateBookmark 更新书签
func (s *BookmarkServiceImpl) UpdateBookmark(ctx context.Context, bookmarkID string, bookmark *readermodels.Bookmark) error {
	id, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %w", err)
	}

	// 获取现有书签
	existing, err := s.bookmarkRepo.GetByID(ctx, id)
	if err != nil {
		return ErrBookmarkNotFound
	}

	// 更新字段
	bookmark.ID = id
	bookmark.UserID = existing.UserID
	bookmark.BookID = existing.BookID
	bookmark.ChapterID = existing.ChapterID

	if bookmark.Note == "" {
		bookmark.Note = existing.Note
	}
	if bookmark.Color == "" {
		bookmark.Color = existing.Color
	}
	if bookmark.Position == 0 {
		bookmark.Position = existing.Position
	}

	return s.bookmarkRepo.Update(ctx, bookmark)
}

// DeleteBookmark 删除书签
func (s *BookmarkServiceImpl) DeleteBookmark(ctx context.Context, bookmarkID string) error {
	id, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %w", err)
	}

	return s.bookmarkRepo.Delete(ctx, id)
}

// ExportBookmarks 导出书签
func (s *BookmarkServiceImpl) ExportBookmarks(ctx context.Context, userID, format string) ([]byte, string, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, "", fmt.Errorf("invalid user ID: %w", err)
	}

	// 获取所有书签
	bookmarks, _, err := s.bookmarkRepo.GetByUser(ctx, userOID, nil, 1, 10000)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get bookmarks: %w", err)
	}

	switch strings.ToLower(format) {
	case "json":
		return s.exportJSON(bookmarks)
	case "csv":
		return s.exportCSV(bookmarks)
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportJSON 导出为JSON格式
func (s *BookmarkServiceImpl) exportJSON(bookmarks []*readermodels.Bookmark) ([]byte, string, error) {
	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return nil, "", err
	}
	return data, "application/json", nil
}

// exportCSV 导出为CSV格式
func (s *BookmarkServiceImpl) exportCSV(bookmarks []*readermodels.Bookmark) ([]byte, string, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"书名", "章节", "位置", "笔记", "颜色", "引用", "创建时间"}
	writer.Write(headers)

	// 写入数据
	for _, b := range bookmarks {
		record := []string{
			b.BookID.Hex(),
			b.ChapterID.Hex(),
			fmt.Sprintf("%d", b.Position),
			b.Note,
			b.Color,
			b.Quote,
			b.CreatedAt.Format(time.RFC3339),
		}
		writer.Write(record)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	return []byte(buf.String()), "text/csv", nil
}

// GetBookmarkStats 获取书签统计
func (s *BookmarkServiceImpl) GetBookmarkStats(ctx context.Context, userID string) (*readermodels.BookmarkStats, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	stats, err := s.bookmarkRepo.GetStats(ctx, userOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}

// SearchBookmarks 搜索书签
func (s *BookmarkServiceImpl) SearchBookmarks(ctx context.Context, userID, keyword string, page, size int) (*BookmarkListResponse, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	bookmarks, total, err := s.bookmarkRepo.Search(ctx, userOID, keyword, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to search bookmarks: %w", err)
	}

	return &BookmarkListResponse{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		Size:      size,
	}, nil
}

// ImportBookmarks 导入书签
func (s *BookmarkServiceImpl) ImportBookmarks(ctx context.Context, userID string, reader io.Reader) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// 解析CSV
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	// 跳过表头
	for i, record := range records {
		if i == 0 {
			continue // 跳过表头
		}

		// 解析记录
		position := 0
		fmt.Sscanf(record[2], "%d", &position)

		bookmark := &readermodels.Bookmark{
			UserID:    userOID,
			Color:     record[4],
			Note:      record[3],
			Position:  position,
			IsPublic:  false,
			CreatedAt: time.Now(),
		}

		// 尝试解析ID
		if bookID, err := primitive.ObjectIDFromHex(record[0]); err == nil {
			bookmark.BookID = bookID
		}
		if chapterID, err := primitive.ObjectIDFromHex(record[1]); err == nil {
			bookmark.ChapterID = chapterID
		}

		// 创建书签
		s.bookmarkRepo.Create(ctx, bookmark)
	}

	return nil
}
