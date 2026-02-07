package bookstore

import (
	"context"
	"errors"
	"fmt"

	"Qingyu_backend/models/bookstore"
)

// BookStreamRepositoryInterface 书籍流式仓储接口（为了解耦）
type BookStreamRepositoryInterface interface {
	StreamSearch(ctx context.Context, filter *bookstore.BookFilter) (BookCursor, error)
	StreamByCursor(ctx context.Context, filter *bookstore.BookFilter) (BookCursor, error)
	CountWithFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error)
	GetCursorManager() CursorManagerInterface
}

// BookCursor 数据库游标接口
type BookCursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
	Err() error
}

// CursorManagerInterface 游标管理器接口
type CursorManagerInterface interface {
	EncodeCursor(cursorType bookstore.CursorType, value interface{}) (string, error)
	DecodeCursor(encoded string) (*bookstore.StreamCursor, error)
	GenerateNextCursor(book *bookstore.Book, cursorType bookstore.CursorType, sortField string) (string, error)
}

// BookstoreStreamService 书城流式服务
type BookstoreStreamService struct {
	streamRepo BookStreamRepositoryInterface
}

// NewBookstoreStreamService 创建书城流式服务
func NewBookstoreStreamService(streamRepo BookStreamRepositoryInterface) *BookstoreStreamService {
	return &BookstoreStreamService{
		streamRepo: streamRepo,
	}
}

// StreamSearchResult 流式搜索结果
type StreamSearchResult struct {
	Books     []*bookstore.Book `json:"books"`
	NextCursor string           `json:"nextCursor,omitempty"`
	HasMore   bool             `json:"hasMore"`
	Total     int64            `json:"total,omitempty"`
}

// StreamSearch 流式搜索书籍
func (s *BookstoreStreamService) StreamSearch(ctx context.Context, filter *bookstore.BookFilter) (*StreamSearchResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	// 设置默认值
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // 最大每批100条
	}

	// 执行流式搜索
	cursor, err := s.streamRepo.StreamSearch(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("stream search failed: %w", err)
	}
	defer cursor.Close(context.Background())

	// 读取数据
	var books []*bookstore.Book
	for cursor.Next(ctx) {
		var book bookstore.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("decode book failed: %w", err)
		}
		books = append(books, &book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// 生成下一个游标
	var nextCursor string
	hasMore := len(books) == filter.Limit
	if hasMore && len(books) > 0 {
		lastBook := books[len(books)-1]
		cm := s.streamRepo.GetCursorManager()
		cursorType := s.getCursorType(filter.SortBy)
		nextCursor, err = cm.GenerateNextCursor(lastBook, cursorType, filter.SortBy)
		if err != nil {
			return nil, fmt.Errorf("generate cursor failed: %w", err)
		}
	}

	return &StreamSearchResult{
		Books:      books,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// StreamSearchWithCursor 使用游标继续流式搜索
func (s *BookstoreStreamService) StreamSearchWithCursor(ctx context.Context, filter *bookstore.BookFilter) (*StreamSearchResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	if filter.Cursor == nil || *filter.Cursor == "" {
		return s.StreamSearch(ctx, filter)
	}

	// 设置默认值
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// 执行基于游标的流式搜索
	cursor, err := s.streamRepo.StreamByCursor(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("stream by cursor failed: %w", err)
	}
	defer cursor.Close(context.Background())

	// 读取数据
	var books []*bookstore.Book
	for cursor.Next(ctx) {
		var book bookstore.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("decode book failed: %w", err)
		}
		books = append(books, &book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// 生成下一个游标
	var nextCursor string
	hasMore := len(books) == filter.Limit
	if hasMore && len(books) > 0 {
		lastBook := books[len(books)-1]
		cm := s.streamRepo.GetCursorManager()
		cursorType := s.getCursorType(filter.SortBy)
		nextCursor, err = cm.GenerateNextCursor(lastBook, cursorType, filter.SortBy)
		if err != nil {
			return nil, fmt.Errorf("generate cursor failed: %w", err)
		}
	}

	return &StreamSearchResult{
		Books:      books,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// CountWithFilter 使用过滤条件统计书籍数量
func (s *BookstoreStreamService) CountWithFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	if filter == nil {
		filter = &bookstore.BookFilter{}
	}

	return s.streamRepo.CountWithFilter(ctx, filter)
}

// getCursorType 根据排序字段确定游标类型
func (s *BookstoreStreamService) getCursorType(sortBy string) bookstore.CursorType {
	switch sortBy {
	case "created_at", "updated_at", "published_at":
		return bookstore.CursorTypeTimestamp
	case "_id":
		return bookstore.CursorTypeID
	default:
		return bookstore.CursorTypeTimestamp
	}
}
